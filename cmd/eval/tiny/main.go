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
		// &bot.Pro{
		// 	Alphas: []float32{6, 4, 6},
		// 	Betas:  []float32{-1, 0, -0.9},
		// 	Gammas: []float32{2, 2, 2},
		// 	Ks:     []float32{0, 0, 0},
		// },
		&bot.Lineal{
			LowerBounds: []float32{
				0.99, // 0 flor/noquiero
				0.99, // 1 contraflor/flor
				0.99, // 2 contrafloralrest/contraflor
				0.99, // 3 quiero_contraflor/noquiero_contraflor
				0.99, // 4 contraflor_alrest/quiero
				0.99, // 5 quiero_contrafloralresto/noquiero_contrafloralrest

				0.99, // 6 envido/noquiero_envido
				0.99, // 7 realenvido/envido
				0.99, // 8 faltaenvido/realenvido

				0.99, // 9 quiero_envido/noquiero_envido
				0.99, // 10 real_envido/quiero_envido
				0.99, // 11 falta_envid/real_envido

				0.99, // 12 quiero_realenvido/noquiero_realenvido
				0.99, // 13 faltaenvido/quiero_realevnido

				0.99, // 14 quiero_faltaenvid/noquiero_faltaenvido

				0.99, // 15 retruco/ (
				0.99, // 16 quiero_truco/noquiero_truco)
				0.99, // 17 vale4/ (
				0.99, // 18 quiero_retruco/noquiero_retruco)
				0.99, // 19 quiero_vale4/noquiero_vale4

				0.99, // 20 truco/nada
				0.99, // 21 retruco/nada
				0.99, // 22 vale4/nada

				0.99, // 23 mazo/seguir
			},
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
