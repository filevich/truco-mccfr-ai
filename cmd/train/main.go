package main

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/filevich/truco-cfr/abs"
	"github.com/filevich/truco-cfr/cfr"
	"github.com/filevich/truco-cfr/eval"
)

const (
	save_dir = "/tmp"
)

func main() {

	threads := 1
	num_players := 2
	tiny_eval := 1_000

	trainer := cfr.New_Trainer(cfr.ESVMCCFR_T, num_players, &abs.A1{})

	// trainer := cfr.Load(
	// 	cfr.CFR_T,
	// 	"/media/jp/DATA/models/2p/models-24h+48p-multi-core/extension-2d/a2/final_CFR_d24h2m_D24h0m_t8435_p0_a2_2205262321.json")

	// trainer := cfr.Load_model(
	// 	"/media/jp/DATA/models/2p/models-24h+48p-multi-core/extension-2d/a2/final_cfr_d48h4m_D48h0m_t26190_p0_a2_2210141233.model",
	// 	true,
	// 	1_000_000)

	// tiny eval
	log.Println("loading t1k22")
	var ds eval.Dataset = eval.Load_dataset("eval/t1k22.json")
	log.Println("done loading t1k22")

	post_save := func() {
		agent := &cfr.BotCFR{
			N:     trainer.String(),
			Model: trainer,
		}
		log.Println("tiny evaluating")
		ale, det, di1, di2, wu_ale, wd_ale, wu_det, wd_det, delta := eval.Tiny_eval_float(agent, num_players, ds[:tiny_eval])
		log.Printf("%s\n\n", eval.Format_Tiny_eval(ale, det, di1, di2, wu_ale, wd_ale, wu_det, wd_det, delta))
		runtime.GC()
	}

	post_save()

	// trainer.Train(
	// 	&cfr.ProfileTime{
	// 		TotalRunningTime:  24 * time.Hour,
	// 		Prunning_treshold: cfr.NEVER,
	// 		// multi
	// 		Threads: threads,
	// 		Mu:      &sync.Mutex{},
	// 		// io
	// 		Save_every:  25 * time.Hour,
	// 		Silent:      true,
	// 		Save_dir:    save_dir,
	// 		Save_prefix: "pre_",
	// 		// tiny eval
	// 		PostSave: post_save,
	// 		// GC
	// 		GC_every: 1 * time.Hour,
	// 	},
	// )

	// log.Printf("Resetting strategy sums")
	// trainer.Reset()

	trainer.Train(
		&cfr.ProfileTime{
			TotalRunningTime:  25 * time.Minute,
			Prunning_treshold: cfr.NEVER,
			// multi
			Threads: threads,
			Mu:      &sync.Mutex{},
			// io
			Save_every:  2 * time.Minute,
			Silent:      true,
			Save_dir:    save_dir,
			Save_prefix: "final_",
			// tiny eval
			PostSave:   post_save,
			Eval_every: 1 * time.Minute,
			// GC
			GC_every: 100 * time.Hour,
		},
	)

	// trainer.Train(
	// 	&cfr.ProfileTime{
	// 		TotalRunningTime:  4 * 24 * time.Hour,
	// 		Prunning_treshold: time.Nanosecond,
	// 		// multi
	// 		Threads: threads,
	// 		Mu:      &sync.Mutex{},
	// 		// io
	// 		Save_every:  24 * time.Hour,
	// 		Silent:      true,
	// 		Save_dir:    save_dir,
	// 		Save_prefix: "final_",
	// 		// tiny eval
	// 		PostSave: post_save,
	// 		// GC
	// 		GC_every: 1 * time.Hour,
	// 	},
	// )

}
