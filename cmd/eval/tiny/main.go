package main

import (
	"fmt"
	"log"
	"time"

	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/eval/dataset"
	"github.com/filevich/truco-cfr/eval2"
)

func main() {
	const (
		tiny_eval   = 1_000
		num_players = 2
		b           = "/media/jp/6e5bdfb0-c84b-4144-8d6d-4688934f1afe/models/6p/48np-multi6/a1"
	)

	log.Println("loading t1k22")
	var ds dataset.Dataset = dataset.LoadDataset("t1k22.json")
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
		var (
			rr                  = eval2.PlayMultipleDoubleGames(agent, agents, num_players, ds)
			s                   = ""
			delta time.Duration = 0
		)

		for i, r := range rr {
			s += fmt.Sprintf("%s=%s - ", agents[i].UID(), r)
			delta += r.Delta
		}

		log.Printf("[%2d/%2d] %s: %s %s\n",
			i+1,
			len(agents),
			agent.UID(),
			s,
			delta.Round(time.Second))

		agent.Free()
	}
}
