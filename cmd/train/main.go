package main

import (
	"flag"
	"log"
	"sync"
	"time"

	"github.com/filevich/truco-ai/bot"
	"github.com/filevich/truco-ai/cfr"
	"github.com/filevich/truco-ai/eval"
	"github.com/filevich/truco-ai/eval/dataset"
	"github.com/filevich/truco-ai/utils"
)

// flags
var (
	modelPtr        = flag.String("model", "", "Filepath to .model file to continue training from")
	numPlayersPtr   = flag.Int("p", 2, "Number of players")
	trainerPtr      = flag.String("trainer", "esvmccfr", "CFR variant")
	hashPtr         = flag.String("h", "sha160", "Hash fn")                   // builder
	infoPtr         = flag.String("info", "InfosetRondaBase", "Infoset Impl") // builder
	absPtr          = flag.String("abs", "a1", "Abstraction")                 // builder
	threadsPtr      = flag.Int("threads", 1, "Threads")
	saveDirPtr      = flag.String("dir", "/tmp", "Save directory")
	tinyEvalPtr     = flag.Int("eval", 1_000, "Progress eval length")
	runPtr          = flag.String("run", "30m", "Total run time")
	prunningPtr     = flag.String("prunning", "", "Start prunning after")
	prunningProbPtr = flag.Float64("prunningProb", 0.01, "Pruning prob")
	saveEveryPtr    = flag.String("saveEvery", "10m", "Saving interval")
	evalEveryPtr    = flag.String("evalEvery", "1m", "Eval interval")
	silentPtr       = flag.Bool("silent", true, "Silent model")
	prefixPtr       = flag.String("prefix", "final_", "Model prefix")
	resetPtr        = flag.Bool("reset", false, "Reset strategy sum")
)

func init() {
	flag.Parse()
	if len(*modelPtr) > 0 {
		log.Println("model", *modelPtr)
	} else {
		log.Println("numPlayers", *numPlayersPtr)
		log.Println("trainer", *trainerPtr)
		log.Println("hash", *hashPtr)
		log.Println("info", *infoPtr)
		log.Println("abs", *absPtr)
	}
	log.Println("threads", *threadsPtr)
	log.Println("saveDir", *saveDirPtr)
	log.Println("tinyEval", *tinyEvalPtr)
	log.Println("run", *runPtr)
	log.Println("saveEvery", *saveEveryPtr)
	log.Println("evalEvery", *evalEveryPtr)
	log.Println("silent", *silentPtr)
	log.Println("prefix", *prefixPtr)
	log.Println("reset", *resetPtr)
}

func main() {

	var (
		saveDir          = *saveDirPtr
		threads          = *threadsPtr
		numPlayers       = *numPlayersPtr
		trainerID        = *trainerPtr
		tinyEval         = *tinyEvalPtr
		model            = *modelPtr
		totalRunningTime time.Duration
		prunningTreshold time.Duration
		saveEvery        time.Duration
		evalEvery        time.Duration
		err              error
	)

	if totalRunningTime, err = time.ParseDuration(*runPtr); err != nil {
		panic(err)
	}

	if *prunningPtr == "" {
		prunningTreshold = cfr.NEVER
		log.Println("prunning", "never")
	} else {
		prunningTreshold, err = time.ParseDuration(*prunningPtr)
		log.Println("prunning", prunningTreshold)
		log.Println("prunningProb", *prunningProbPtr)
		if err != nil {
			panic(err)
		}
	}

	if saveEvery, err = time.ParseDuration(*saveEveryPtr); err != nil {
		panic(err)
	}

	if evalEvery, err = time.ParseDuration(*evalEveryPtr); err != nil {
		panic(err)
	}

	var trainer cfr.ITrainer

	if len(model) == 0 {
		trainer = cfr.NewTrainer(
			cfr.Trainer_T(trainerID),
			numPlayers,
			*hashPtr,
			*infoPtr,
			*absPtr)
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
		mem := utils.GetMemUsage()
		log.Println(eval.Fmt(rr, agents), mem)
	}

	if *resetPtr {
		log.Printf("Resetting strategy sums")
		trainer.Reset()
	}

	trainer.Train(
		&cfr.ProfileTime{
			TotalRunningTime: totalRunningTime,
			// pruning
			PrunningTreshold: prunningTreshold,
			PrunningProb:     float32(*prunningProbPtr),
			// multi
			Threads: threads,
			Mu:      &sync.Mutex{},
			// io
			SaveEvery:  saveEvery,
			Silent:     *silentPtr,
			SaveDir:    saveDir,
			SavePrefix: *prefixPtr,
			PostSave:   nil,
			// tiny eval
			EvalEvery: evalEvery,
			Evaluator: evaluator,
			// GC
			GCEvery: 100 * time.Hour,
		},
	)
}
