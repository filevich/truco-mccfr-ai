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

func run(ds dataset.Dataset, proAgent bot.Agent) {
	const (
		numPlayers = 2
		b          = "/Users/jp/Downloads/cluster/train-cfr/models/2p"
	)

	agent := &cfr.BotCFR{
		N: "esvmccfra2",
		F: b + "/a2/final_es-vmccfr_d70h0m_D70h0m_t7077536_p0_a2_2402052107.model",
		// F: b + "/a3/final_es-vmccfr_d70h0m_D70h0m_t3468734_p0_a3_2402052116.model",
	}

	againstThese := []bot.Agent{
		&bot.Simple{},
		proAgent,
	}

	var (
		s  = ""
		rr = eval.PlayMultipleDoubleGames(
			agent,
			againstThese,
			numPlayers,
			ds)
		delta time.Duration = 0
	)

	for i, r := range rr {
		s += fmt.Sprintf("%s=%s - ", againstThese[i].UID(), r)
		delta += r.Delta
	}

	log.Printf("%s: %s %s\n",
		agent.UID(),
		s,
		delta.Round(time.Second))

	agent.Free()
}

func main() {
	log.Println("loading t1k22")
	var ds dataset.Dataset = dataset.LoadDataset("/Users/jp/Workspace/go/truco-ai/truco-ai/t1k22.json")
	log.Println("done loading t1k22")

	run(
		ds,
		&bot.Pro{
			Alphas: []float32{2.9, 2.9, 2.9},
			Betas:  []float32{-0.9, -0.9, -0.9},
			Gammas: []float32{2, 2, 2},
			Ks:     []float32{0.01, 0.01, 0.01},
		},
	)

	// // alphas := []float32{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}  // 10
	// // betas := []float32{-2, -1.5, -1, 0, 0.5, 1, 1.5, 2} // 8

	// alphas := []float32{0, 1, 2, 3, 10}
	// betas := []float32{-2}
	// // gammas := []float32{2}
	// k := []float32{0, 0, 0}

	// for a0i := 0; a0i < len(alphas); a0i++ {

	// 	// a1i, a2i := a0i, a0i
	// 	for a1i := 0; a1i < len(alphas); a1i++ {
	// 		for a2i := 0; a2i < len(alphas); a2i++ {

	// 			for b0i := 0; b0i < len(betas); b0i++ {

	// 				b1i, b2i := b0i, b0i
	// 				// for b1i := 0; b1i < len(betas); b1i++ {
	// 				// 	for b2i := 0; b2i < len(betas); b2i++ {

	// 				a := []float32{alphas[a0i], alphas[a1i], alphas[a2i]}
	// 				b := []float32{betas[b0i], betas[b1i], betas[b2i]}
	// 				g := []float32{2, 2, 2}

	// 				log.Println("using", a, b)

	// 				run(
	// 					ds,
	// 					&bot.Pro{
	// 						Alphas: a,
	// 						Betas:  b,
	// 						Gammas: g,
	// 						Ks:     k,
	// 					},
	// 				)

	// 				// betas
	// 				// 	}
	// 				// }
	// 			}

	// 			// alphas
	// 		}
	// 	}
	// }
}
