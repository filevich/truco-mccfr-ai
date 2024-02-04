package cfr

import (
	"fmt"
	"log"
	"math"
	"runtime"
	"sync"
	"time"

	"github.com/filevich/truco-ai/bot"
	"github.com/filevich/truco-ai/utils"
)

const (
	NEVER = time.Duration(math.MaxInt64)
)

type ProfileTime struct {
	//
	TotalRunningTime time.Duration
	PrunningTreshold time.Duration

	// privadas
	start    time.Time
	lastSave time.Time
	lastEval time.Time
	lastGC   time.Time

	// Multi
	Threads int
	Mu      *sync.Mutex

	// report
	lastInfoLen  int
	firstInfoInc int

	// io
	SaveEvery   time.Duration
	Silent      bool
	FullySilent bool
	SaveDir     string
	SavePrefix  string
	PostSave    func()

	// gc
	GCEvery time.Duration

	// exploitability
	Exploiting bot.Agent

	// eval
	EvalEvery time.Duration
	Evaluator func()
}

// is prunable: iteracion (int ~ 2k), tiempo (tiempo ~ 3hs), porcentage done (float ~ 25%)

func (p *ProfileTime) Init(trainer ITrainer) {
	trainer.SetWorkers(p.Threads)

	now := time.Now()
	p.start = now
	p.lastSave = now
	p.lastEval = time.Now().Add(-1 * 24 * time.Hour)
	p.lastGC = now

	if p.IsFullySilent() {
		return
	}

	// no seteo T porque no se
	mode := "[verbose]"
	if p.Silent {
		mode = "[silent]"
	}

	prunning := "[no prunning]"
	if p.PrunningTreshold != NEVER {
		prunning = fmt.Sprintf("[prunning @ %s]", p.PrunningTreshold.Round(time.Second))
	}

	log.Printf("Running %s x %s for %s [saving every %s][GC every %s] %s %s [%d threads]\n",
		trainer.String(),
		trainer.GetAbs().String(),
		p.TotalRunningTime,
		p.SaveEvery,
		p.GCEvery,
		prunning,
		mode,
		p.Threads,
	)
}

func (p ProfileTime) Continue(trainer ITrainer) bool {
	return time.Since(p.start) < p.TotalRunningTime
}

func (p ProfileTime) IsPrunable(trainer ITrainer) bool {
	return p.PrunningTreshold != NEVER &&
		time.Since(p.start) >= p.PrunningTreshold
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
	return time.Since(p.lastSave) >= p.SaveEvery || (lastIter && lastThread)
}

func (p ProfileTime) shouldEval(trainer ITrainer) bool {
	return time.Since(p.lastEval) >= p.EvalEvery
}

func (p ProfileTime) shouldGC(trainer ITrainer) bool {
	return time.Since(p.lastGC) >= p.GCEvery
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

		total := trainer.CountInfosets() // unsafe, pero no importa
		inc := total - p.lastInfoLen
		if p.firstInfoInc == 0 {
			p.firstInfoInc = inc
		}
		relative_inc := float32(inc) / float32(p.firstInfoInc)
		p.lastInfoLen = total

		AGV := trainer.MaxAvgGameValue()

		log.Printf("[%3.f%%] - (%s/%s) - #iter:%d - AGV:%.4f - #infos: %d (+%d ~ %.3f) %s @%s",
			p.PercentageDone(t),
			time.Since(p.start).Round(time.Second),
			p.TotalRunningTime,
			t,
			AGV,
			total,
			inc,
			relative_inc,
			P,
			utils.MiniCurrentTime()[2:],
		)
	}
}

func (p *ProfileTime) Checkpoint(t ITrainer) {
	// eval ?
	if p.shouldEval(t) {
		if p.Evaluator != nil {
			p.Evaluator()
		}
		p.lastEval = time.Now()
	}

	if p.SaveEvery == NEVER {
		return
	}

	if !p.shouldSave(t) {
		if !p.Silent {
			log.Println()
		}
		return
	}

	p.lastSave = time.Now() // <- seguro porque estoy con el Mu

	// iter 0 *hecha* -> t1
	// ok, debo salvar la estrategia actual
	iter_name := t.Get_t() + 1
	d := time.Since(p.start).Round(time.Minute).String()
	D := p.TotalRunningTime.Round(time.Minute).String()

	prunned := 0
	if p.IsPrunable(t) {
		prunned = 1
	}

	filename := fmt.Sprintf("%s/%s%s_d%s_D%s_t%d_p%d_%s_%s.model",
		p.SaveDir,
		p.SavePrefix,
		t.String(),
		d[:len(d)-2],
		D[:len(D)-2],
		iter_name,
		prunned,
		t.GetAbs().String(),
		utils.MiniCurrentTime(),
	)

	// iter 0 *hecha* -> t1
	// pero la strucy dice Current_Iter=0
	// entonces hago inc&desc

	log.Printf("Saving model at %s\n", filename)
	t.SaveModel(filename, 1_000_000, t.String(), nil)

	// completo
	// p := float32(iter_name) / float32(t.T) * 100
	// log.Printf("[%.f%%] - (%d/%d) %s\n", p, iter_name, t.T, filename)

	if p.PostSave != nil {
		p.PostSave()
	}
}

func (p *ProfileTime) CheckGC(t ITrainer) {
	if p.GCEvery == NEVER {
		return
	}

	if !p.shouldGC(t) {
		if !p.Silent {
			log.Println()
		}
		return
	}

	p.lastGC = time.Now() // <- seguro porque estoy con el Mu
	// log.Printf("GC...")
	runtime.GC()
	// log.Println(" [done]")
}

func (p *ProfileTime) Exploit() bot.Agent {
	return p.Exploiting
}
