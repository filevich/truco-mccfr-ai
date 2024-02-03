package main

import (
	"fmt"

	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/eval"
	"github.com/filevich/truco-cfr/utils"
)

func main() {

	const (
	// base
	// b = "/media/jp/6e5bdfb0-c84b-4144-8d6d-4688934f1afe/models/6p/48np-multi6/a1"
	)

	ds := eval.Load_dataset("eval/t1k22.json")

	// un tournament reune a varios agentes, y los hace pelear a todos contra todos
	torneo := &eval.TBinomial{
		Num_players: 6,
		Partidas:    eval.New_Tabla(),
		Agents: []bot.Agent{
			// &eval.BotCFR{
			// 	N: "ESL-A1",
			// 	F: b + "/final_es-lmccfr_d48h1m_D48h0m_t398613799_p1_a1_2210061829.model",
			// },

			// baselines
			&bot.BotSimple{},
			&bot.BotRandom{},
		},
	}

	torneo.Start(ds[:], true)

	torneo.Report()
	fmt.Println()

	// guardo el resultado
	t := utils.Mini_Current_time()
	utils.Write(torneo.Partidas, "/tmp/res-"+t+".json", true)
	fmt.Printf("resultado guardado en %s\n\n", "/tmp/res-"+t+".json")

}
