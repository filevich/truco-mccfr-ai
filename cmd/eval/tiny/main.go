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
		b          = "/Users/jp/Downloads/output/truco/go-cfr-models/models"
	)

	log.Println("loading t1k22")
	var ds dataset.Dataset = dataset.LoadDataset("t1k22.json")
	log.Println("done loading t1k22")

	testThese := []eval.Agent{
		&cfr.BotCFR{
			ID:       "botcfr",
			Filepath: b + "/2p/irb-a3/pruned_esvmccfr_d70h0m_D70h0m_t288652014_p1_a3_2402151230.model",
		},
		&cfr.BotCFR{
			ID:       "twin",
			Filepath: b + "/2p/irb-a3/pruned_esvmccfr_d70h0m_D70h0m_t288652014_p1_a3_2402151230.model",
		},
	}

	againstThese := []eval.Agent{
		&bot.Random{},
		&bot.Simple{},
		&bot.SimpleX{},
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
