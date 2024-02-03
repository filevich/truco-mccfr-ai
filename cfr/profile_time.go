package cfr

import (
	"fmt"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/utils"
)

const (
	NEVER = time.Duration(math.MaxInt64)
)

type ProfileTime struct {
	//
	TotalRunningTime  time.Duration
	Prunning_treshold time.Duration

	// privadas
	start     time.Time
	last_save time.Time
	last_eval time.Time
	last_GC   time.Time

	// Multi
	Threads int
	Mu      *sync.Mutex

	// report
	last_info_len  int
	first_info_inc int

	// io
	Save_every  time.Duration
	Silent      bool
	FullySilent bool
	Save_dir    string
	Save_prefix string

	// gc
	GC_every time.Duration

	// exploitability
	Exploiting bot.Agent

	// eval
	Eval_every time.Duration
	PostSave   func()
}

// is prunable: iteracion (int ~ 2k), tiempo (tiempo ~ 3hs), porcentage done (float ~ 25%)

func (p *ProfileTime) Init(trainer ITrainer) {
	trainer.SetWorkers(p.Threads)

	now := time.Now()
	p.start = now
	p.last_save = now
	p.last_eval = time.Now().Add(-1 * 24 * time.Hour)
	p.last_GC = now

	if p.IsFullySilent() {
		return
	}

	// no seteo T porque no se
	mode := "[verbose]"
	if p.Silent {
		mode = "[silent]"
	}

	prunning := "[no prunning]"
	if p.Prunning_treshold != NEVER {
		prunning = fmt.Sprintf("[prunning @ %s]", p.Prunning_treshold.Round(time.Second))
	}

	fmt.Printf("\nRunning %s x %s for %s [saving every %s][GC every %s] %s %s [%d threads] starting at %s\n\n",
		trainer.String(),
		trainer.Get_abs().String(),
		p.TotalRunningTime,
		p.Save_every,
		p.GC_every,
		prunning,
		mode,
		p.Threads,
		utils.Current_time_file_format(),
	)
}

func (p ProfileTime) Continue(trainer ITrainer) bool {
	return time.Since(p.start) < p.TotalRunningTime
}

func (p ProfileTime) IsPrunable(trainer ITrainer) bool {
	return p.Prunning_treshold != NEVER &&
		time.Since(p.start) >= p.Prunning_treshold
}

func (p ProfileTime) IsMulti() bool {
	return p.Threads > 1
}

func (p ProfileTime) GetThreads() int {
	return p.Threads
}

func (p ProfileTime) IsSilent() bool {
	return p.Silent
}

func (p ProfileTime) IsFullySilent() bool {
	return p.FullySilent
}

func (p ProfileTime) shouldSave(trainer ITrainer) bool {
	// en caso de que haya caducado el tiempo
	// el que tiene que salvar es el ULTIMO THREAD vivo
	lastIter := time.Since(p.start) >= p.TotalRunningTime
	lastThread := false
	if lastIter {
		// me desprocupo aca, porque para editar trainer.Working
		// necesitar p.Mu, y se que solo yo lo tengo
		lastThread = trainer.AllDones()
	}
	// lo mismo con el last_save
	return time.Since(p.last_save) >= p.Save_every || (lastIter && lastThread)
}

func (p ProfileTime) shouldEval(trainer ITrainer) bool {
	return time.Since(p.last_eval) >= p.Eval_every
}

func (p ProfileTime) shouldGC(trainer ITrainer) bool {
	return time.Since(p.last_GC) >= p.GC_every
}

func (p ProfileTime) PercentageDone(t int) float32 {
	done := time.Since(p.start).Seconds() / p.TotalRunningTime.Seconds()
	return float32(done * 100)
}

func (p *ProfileTime) Check(trainer ITrainer) {
	// estoy checkeado, solo yo
	// solo checkea 1 por vez
	p.Mu.Lock()
	defer p.Mu.Unlock()

	// termine?
	termine := !p.Continue(trainer) // el thread termino sii se acabo el tiempo.
	if termine {
		trainer.WorkerDone() // <- aca actualiza al cantidad de workers remaining
	}

	// incrementos
	trainer.inc_T()

	p.PrintProgress(trainer)
	p.Checkpoint(trainer) // <- aca actualiza el p.last_save
	p.CheckGC(trainer)
}

func (p *ProfileTime) PrintProgress(trainer ITrainer) {
	t := trainer.Get_t() + 1
	verbose := !p.IsSilent()

	if p.shouldSave(trainer) || verbose {

		P := ""
		if p.IsPrunable(trainer) {
			P = " [P]"
		}

		total := trainer.Count_infosets() // unsafe, pero no importa
		inc := total - p.last_info_len
		if p.first_info_inc == 0 {
			p.first_info_inc = inc
		}
		relative_inc := float32(inc) / float32(p.first_info_inc)
		p.last_info_len = total

		AGV := trainer.Max_Avg_Game_Value()

		fmt.Printf("![%3.f%%] - (%s/%s) - #iter:%d - AGV:%.4f - #infos: %d (+%d ~ %.3f) %s @%s",
			p.PercentageDone(t),
			time.Since(p.start).Round(time.Second),
			p.TotalRunningTime,
			t,
			AGV,
			total,
			inc,
			relative_inc,
			P,
			utils.Mini_Current_time()[2:],
		)

		// if !p.IsSilent() {
		// 	fmt.Println()
		// }
	}
}

func (p *ProfileTime) Checkpoint(t ITrainer) {
	// eval ?
	if p.shouldEval(t) {
		if p.PostSave != nil {
			p.PostSave()
		}
		p.last_eval = time.Now()
	}

	if p.Save_every == NEVER {
		return
	}

	if !p.shouldSave(t) {
		if !p.Silent {
			fmt.Println()
		}
		return
	}

	// anres de guardar, corro el GC 5 veces
	total := 5
	for i := 0; i < total; i++ {
		fmt.Printf("PRE save GC...")
		runtime.GC()
		fmt.Printf("; sleep 10s... [%d/%d]\n", i+1, total)
		time.Sleep(time.Second * 10)
	}

	p.last_save = time.Now() // <- seguro porque estoy con el Mu

	// iter 0 *hecha* -> t1
	// ok, debo salvar la estrategia actual
	iter_name := t.Get_t() + 1
	d := time.Since(p.start).Round(time.Minute).String()
	D := p.TotalRunningTime.Round(time.Minute).String()

	prunned := 0
	if p.IsPrunable(t) {
		prunned = 1
	}

	filename := fmt.Sprintf("%s/%s%s_d%s_D%s_t%d_p%d_%s_%s.json",
		p.Save_dir,
		p.Save_prefix,
		t.String(),
		d[:len(d)-2],
		D[:len(D)-2],
		iter_name,
		prunned,
		t.Get_abs().String(),
		utils.Mini_Current_time(),
	)

	dotmodel := filename[:len(filename)-len(".json")] + ".model"

	// iter 0 *hecha* -> t1
	// pero la strucy dice Current_Iter=0
	// entonces hago inc&desc

	// deprecated
	// t.Save(filename)
	t.Save_model(dotmodel, 1_000_000, t.String(), nil)

	// report el filesave?
	// minimo
	fmt.Println(" [*]")

	// completo
	// p := float32(iter_name) / float32(t.T) * 100
	// fmt.Printf("[%.f%%] - (%d/%d) %s\n", p, iter_name, t.T, filename)

	if p.PostSave != nil {
		p.PostSave()
	}
}

func (p *ProfileTime) CheckGC(t ITrainer) {
	if p.GC_every == NEVER {
		return
	}

	if !p.shouldGC(t) {
		if !p.Silent {
			fmt.Println()
		}
		return
	}

	p.last_GC = time.Now() // <- seguro porque estoy con el Mu
	// fmt.Printf("GC...")
	runtime.GC()
	// fmt.Println(" [done]")
}

func (p *ProfileTime) Exploit() bot.Agent {
	return p.Exploiting
}
