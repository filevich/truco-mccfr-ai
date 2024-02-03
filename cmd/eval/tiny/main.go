package main

import (
	"log"

	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/eval"
)

func main() {
	const (
		tiny_eval   = 1_000
		num_players = 2
		b           = "/media/jp/6e5bdfb0-c84b-4144-8d6d-4688934f1afe/models/6p/48np-multi6/a1"
	)

	log.Println("loading t1k22")
	var ds eval.Dataset = eval.LoadDataset("eval/t1k22.json")
	log.Println("done loading t1k22")

	agents := []bot.Agent{
		&bot.Random{},
		&bot.Simple{},

		// &bot.BotCFR{
		// 	N: "final_es-lmccfr_d25h0m_D48h0m_t24878_p0_a1_2208092259",
		// 	F: b + "/final_es-lmccfr_d25h0m_D48h0m_t24878_p0_a1_2208092259.model",
		// },
	}

	for i, agent := range agents {
		agent.Initialize()
		res := eval.TinyEval(agent, num_players, ds[:tiny_eval])
		log.Printf("[%2d/%2d] %s: %s",
			i+1,
			len(agents),
			agent.UID(),
			res.String())
		agent.Free()
	}
}
