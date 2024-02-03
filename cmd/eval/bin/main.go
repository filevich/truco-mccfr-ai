package main

import (
	"log"

	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/eval/dataset"
	"github.com/filevich/truco-cfr/eval2"
	"github.com/filevich/truco-cfr/utils"
)

func main() {

	const (
	// base
	// b = "/media/jp/6e5bdfb0-c84b-4144-8d6d-4688934f1afe/models/6p/48np-multi6/a1"
	)

	ds := dataset.LoadDataset("t1k22.json")

	// un tournament reune a varios agentes, y los hace pelear a todos contra todos
	torneo := &eval2.Tournament{
		NumPlayers: 4,
		Table:      eval2.NewTable(),
		Agents: []bot.Agent{
			// &eval2.BotCFR{
			// 	N: "ESL-A1",
			// 	F: b + "/final_es-lmccfr_d48h1m_D48h0m_t398613799_p1_a1_2210061829.model",
			// },

			// baselines
			&bot.Simple{},
			&bot.Random{},
		},
	}

	torneo.Start(ds[:], true)

	torneo.Report()
	log.Println()

	// guardo el resultado
	t := utils.MiniCurrentTime()
	utils.Write(torneo.Table, "/tmp/res-"+t+".json", true)
	log.Printf("resultado guardado en %s\n\n", "/tmp/res-"+t+".json")

}
