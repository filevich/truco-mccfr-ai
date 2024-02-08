package main

import (
	"fmt"
	"log"
	"time"

	"github.com/filevich/truco-ai/bot"
	"github.com/filevich/truco-ai/cfr"
	"github.com/filevich/truco-ai/eval"
	"github.com/filevich/truco-ai/eval/dataset"
)

func main() {
	const (
		tinyEval   = 1_000
		numPlayers = 2
		b          = "/Users/jp/Downloads/cluster/train-cfr/models/2p"
	)

	log.Println("loading t1k22")
	var ds dataset.Dataset = dataset.LoadDataset("t1k22.json")
	log.Println("done loading t1k22")

	testThese := []bot.Agent{
		&cfr.BotCFR{
			N: "esvmccfra2",
			F: b + "/a2/final_es-vmccfr_d70h0m_D70h0m_t7077536_p0_a2_2402052107.model",
		},
		&cfr.BotCFR{
			N: "esvmccfra3",
			F: b + "/a3/final_es-vmccfr_d70h0m_D70h0m_t3468734_p0_a3_2402052116.model",
		},
	}

	againstThese := []bot.Agent{
		&bot.Random{},
		&bot.Simple{},
		&bot.Pro{
			Alphas: []float32{6, 4, 6},
			Betas:  []float32{-1, 0, -0.9},
			Gammas: []float32{2, 2, 2},
			Ks:     []float32{0, 0, 0},
		},
	}

	for i, agent := range testThese {
		var (
			rr = eval.PlayMultipleDoubleGames(
				agent,
				againstThese,
				numPlayers,
				ds[:tinyEval])
			s                   = ""
			delta time.Duration = 0
		)

		for i, r := range rr {
			s += fmt.Sprintf("%s=%s - ", againstThese[i].UID(), r)
			delta += r.Delta
		}

		log.Printf("[%2d/%2d] %s: %s %s\n",
			i+1,
			len(testThese),
			agent.UID(),
			s,
			delta.Round(time.Second))

		agent.Free()
	}
}
