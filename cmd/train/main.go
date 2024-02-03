package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/filevich/truco-cfr/abs"
	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/cfr"
	"github.com/filevich/truco-cfr/eval"
	"github.com/filevich/truco-cfr/eval/dataset"
)

// flags
var (
	modelPtr      = flag.String("model", "", "Filepath to .model file to continue training")
	numPlayersPtr = flag.Int("p", 2, "Number of players")
	trainerPtr    = flag.String("trainer", "es-vmccfr", "CFR variant")
	absPtr        = flag.String("abs", "a1", "Abstraction")
	threadsPtr    = flag.Int("threads", 1, "Threads")
	saveDirPtr    = flag.String("dir", "/tmp", "Save directory")
	tinyEvalPtr   = flag.Int("eval", 1_000, "Progress eval length.")
)

func init() {
	flag.Parse()
	if len(*modelPtr) > 0 {
		log.Println("model", *modelPtr)
	} else {
		log.Println("numPlayers", *numPlayersPtr)
		log.Println("algo", *trainerPtr)
		log.Println("abs", *absPtr)
	}
	log.Println("threads", *threadsPtr)
	log.Println("saveDir", *saveDirPtr)
	log.Println("tinyEval", *tinyEvalPtr)
}

func main() {

	var (
		saveDir     = *saveDirPtr
		threads     = *threadsPtr
		numPlayers  = *numPlayersPtr
		algo        = *trainerPtr
		tinyEval    = *tinyEvalPtr
		model       = *modelPtr
		abstraction = abs.ParseAbs(*absPtr)
	)

	var trainer cfr.ITrainer

	if len(model) == 0 {
		trainer = cfr.NewTrainer(
			cfr.Trainer_T(algo),
			numPlayers,
			abstraction)
	} else {
		trainer = cfr.LoadModel(model, true, 1_000_000)
	}

	// tiny eval
	log.Println("Loading t1k22")
	var ds dataset.Dataset = dataset.LoadDataset("t1k22.json")
	log.Println("Done loading t1k22")

	agents := []bot.Agent{
		&bot.Random{},
		&bot.Simple{},
		// &bot.BotCFR{
		// 	N: "final_es-lmccfr_d25h0m_D48h0m_t24878_p0_a1_2208092259",
		// 	F: b + "/final_es-lmccfr_d25h0m_D48h0m_t24878_p0_a1_2208092259.model",
		// },
	}

	evaluator := func() {
		agent := &cfr.BotCFR{
			N:     trainer.String(),
			Model: trainer,
		}
		rr := eval.PlayMultipleDoubleGames(agent, agents, numPlayers, ds[:tinyEval])
		log.Println(eval.Fmt(rr, agents))
	}

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
			TotalRunningTime: 25 * time.Minute,
			PrunningTreshold: cfr.NEVER,
			// multi
			Threads: threads,
			Mu:      &sync.Mutex{},
			// io
			SaveEvery:  2 * time.Minute,
			Silent:     true,
			SaveDir:    saveDir,
			SavePrefix: "final_",
			PostSave:   nil,
			// tiny eval
			EvalEvery: 1 * time.Minute,
			Evaluator: evaluator,
			// GC
			GCEvery: 100 * time.Hour,
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
