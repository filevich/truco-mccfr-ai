package main

import (
	"flag"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/filevich/truco-mccfr-ai/bot"
	"github.com/filevich/truco-mccfr-ai/cfr"
	"github.com/filevich/truco-mccfr-ai/eval"
	"github.com/filevich/truco-mccfr-ai/eval/dataset"
	"github.com/filevich/truco-mccfr-ai/internal/tournamentclient"
	"github.com/filevich/truco-mccfr-ai/utils"
	"github.com/truquito/gotruco"
)

// flags
var (
	modelPtr          = flag.String("model", "", "Filepath to .model file to continue training from")
	numPlayersPtr     = flag.Int("p", 2, "Number of players")
	trainerPtr        = flag.String("trainer", "esvmccfr", "CFR variant")
	hashPtr           = flag.String("hash", "sha160", "Hash fn")                // builder
	infoPtr           = flag.String("info", "InfosetRondaBase", "Infoset Impl") // builder
	absPtr            = flag.String("abs", "a1", "Abstraction")                 // builder
	threadsPtr        = flag.Int("threads", 1, "Threads")
	saveDirPtr        = flag.String("dir", "/tmp", "Save directory")
	tinyEvalPtr       = flag.Int("eval", 1_000, "Progress eval length")
	runPtr            = flag.String("run", "30m", "Total run time")
	prunningPtr       = flag.String("prunning", "", "Start prunning after")
	prunningProbPtr   = flag.Float64("prunning_prob", 0.01, "Pruning prob")
	saveEveryPtr      = flag.String("save_every", "10m", "Saving interval")
	evalEveryPtr      = flag.String("eval_every", "1m", "Eval interval")
	silentPtr         = flag.Bool("silent", true, "Silent model")
	prefixPtr         = flag.String("prefix", "final_", "Model prefix")
	fmtPtr            = flag.String("fmt", "", "Model name format")
	resetPtr          = flag.Bool("reset", false, "Reset strategy sum")
	tournamentAddrPtr = flag.String("tournament-addr", "", "Tournament server address (empty to skip tournament eval)")
	tournamentNPtr    = flag.Int("tournament-n", 1000, "Number of (doubled) games for tournament eval")
	tournamentNamePtr = flag.String("tournament-name", "", "Bot name for tournament eval")
)

func init() {
	flag.Parse()
	slog.Info(
		"START",
		"model", *modelPtr,
		"numPlayers", *numPlayersPtr,
		"trainer", *trainerPtr,
		"hash", *hashPtr,
		"info", *infoPtr,
		"abs", *absPtr,
		"threads", *threadsPtr,
		"saveDir", *saveDirPtr,
		"tinyEval", *tinyEvalPtr,
		"run", *runPtr,
		"prunning", *prunningPtr,
		"prunningProb", *prunningProbPtr,
		"saveEvery", *saveEveryPtr,
		"evalEvery", *evalEveryPtr,
		"silent", *silentPtr,
		"prefix", *prefixPtr,
		"fmt", *fmtPtr,
		"reset", *resetPtr,
		"tournamentAddr", *tournamentAddrPtr,
		"tournamentN", *tournamentNPtr,
		"tournamentName", *tournamentNamePtr,
		"gotruco", gotruco.VERSION,
	)
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
	} else {
		prunningTreshold, err = time.ParseDuration(*prunningPtr)
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
		trainer = cfr.LoadModel(model, true, 1_000_000, false)
	}

	// tiny eval
	slog.Info("LOADING_t1k22")
	tic := time.Now()
	var ds dataset.Dataset = dataset.LoadDataset("t1k22.json")
	slog.Info("FINISHED_LOADING_t1k22", "delta", time.Since(tic))

	agents := []eval.Agent{
		&bot.Random{},
		&bot.Simple{},
	}

	evaluator := func() {
		agent := &cfr.BotCFR{
			ID:    trainer.String(),
			Model: trainer,
		}
		rr := eval.PlayMultipleDoubleGames(agent, agents, numPlayers, ds[:tinyEval])
		infos := trainer.CountInfosets()

		heapAlloc, totalAlloc, sys := utils.GetMemUsageMiB()

		var delta time.Duration = 0

		// general progress info
		slog.Info("REPORT", "infos", infos, "iters", trainer.Get_t())

		for i, r := range rr {
			delta += r.Delta
			u, l := r.WaldInterval(true)
			slog.Info(
				"RESULTS",
				"opponent", agents[i].UID(),
				"wr", r.WP(),
				"wald_interval_upper", u,
				"wald_interval_lower", l,
				"di", r.Dumbo1,
			)
		}
		slog.Info("EVAL_DONE", "delta", delta)
		slog.Info(
			"MEMORY",
			"heapAlloc", heapAlloc,
			"totalAlloc", totalAlloc,
			"sys", sys,
		)

		// Tournament evaluation (if configured)
		if *tournamentAddrPtr != "" {
			slog.Info("TOURNAMENT_EVAL_START", "addr", *tournamentAddrPtr, "n", *tournamentNPtr)
			tic := time.Now()
			err := tournamentclient.RunChallenge(
				*tournamentAddrPtr,
				*tournamentNamePtr,
				*tournamentNPtr,
				trainer,
			)
			if err != nil {
				slog.Error("TOURNAMENT_EVAL_ERROR", "error", err)
			} else {
				slog.Info("TOURNAMENT_EVAL_DONE", "delta", time.Since(tic))
			}
		}
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
			SaveFormat: *fmtPtr,
			PostSave:   nil,
			// tiny eval
			EvalEvery: evalEvery,
			Evaluator: evaluator,
			// GC
			GCEvery: 100 * time.Hour,
		},
	)
}
