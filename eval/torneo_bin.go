package eval

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/eval/dataset"
	"github.com/jedib0t/go-pretty/v6/table"
)

type TBinomial struct {
	Num_players          int
	Cant_partidas_dobles int
	Agents               []bot.Agent
	Partidas             Tabla
}

func (t *TBinomial) Start(ds dataset.Dataset, verbose bool) {

	t.Cant_partidas_dobles = len(ds)

	if verbose {
		log.Printf("\nTorneo Binomial de %dp:\n", t.Num_players)
		for ix, agent := range t.Agents {
			log.Printf("\t%2d. %s\n", ix+1, agent.UID())
		}

		log.Printf("\nDone: ")
	}

	for i := 0; i < len(t.Agents)-1; i++ {
		agent1 := t.Agents[i]
		agent1.Initialize()
		for j := i + 1; j < len(t.Agents); j++ {
			agent2 := t.Agents[j]
			agent2.Initialize()
			res_partidas := SimPartidasBin(ds, agent1, agent2, t.Num_players)
			t.Partidas.Add(
				agent1.UID(),
				agent2.UID(),
				res_partidas,
			)
			// termino de jugar contra agent2 -> ya no lo necesito
			agent2.Free()
			runtime.GC()
		}

		if verbose {
			log.Printf("%s", agent1.UID())
			if i == len(t.Agents)-2 {
				log.Printf("\n\n")
			} else {
				log.Printf(", ")
			}
		}

		// termino de jugar contra agent2 -> ya no lo necesito
		agent1.Free()
		runtime.GC()
	}
}

func (t *TBinomial) Inscriptos(subtitle string) []interface{} {
	// "0.857   8.437" ~> tiene len=20
	res := make([]interface{}, 0)
	for _, agent := range t.Agents {
		res = append(res, fmt.Sprintf("%s\n  %s ", agent.UID(), subtitle))
	}
	return res
}

func (torneo *TBinomial) PrintTablaMedia(tabla Tabla) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		append(
			table.Row{"A\\B"},
			append(
				torneo.Inscriptos("WP    ADP"),
				"B/A",
			)...,
		),
	)

	for i := 0; i < len(torneo.Agents); i++ {

		agent1 := torneo.Agents[i]
		row := []interface{}{
			agent1.UID(),
		}

		for j := 0; j < len(torneo.Agents); j++ {
			if i == j {
				row = append(row, "             ")
				continue
			}
			agent2 := torneo.Agents[j]
			wp, adp := tabla.Metrics(agent1.UID(), agent2.UID())
			row = append(row, fmt.Sprintf("%.3f  %.3f", wp, adp))
		}

		row = append(row, agent1.UID())

		t.AppendRow(row)
		t.AppendSeparator()
	}

	// t.AppendFooter(table.Row{"", "", "Total", 10000})
	t.Render()
}

func (torneo *TBinomial) PrintTablaWaldInterval(tabla Tabla) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		append(
			table.Row{"A\\B"},
			append(
				torneo.Inscriptos("L    U"),
				"B/A",
			)...,
		),
	)

	for i := 0; i < len(torneo.Agents); i++ {

		agent1 := torneo.Agents[i]
		row := []interface{}{
			agent1.UID(),
		}

		for j := 0; j < len(torneo.Agents); j++ {
			if i == j {
				row = append(row, "             ")
				continue
			}
			agent2 := torneo.Agents[j]
			// wp, adp := tabla.Metrics(agent1.UID(), agent2.UID())
			// row = append(row, fmt.Sprintf("%.3f  %.3f", wp, adp))
			u, l := tabla.WaldInterval(agent1.UID(), agent2.UID())
			row = append(row, fmt.Sprintf("%.3f %.3f", l, u))
		}

		row = append(row, agent1.UID())

		t.AppendRow(row)
		t.AppendSeparator()
	}

	// t.AppendFooter(table.Row{"", "", "Total", 10000})
	t.Render()
}

// func (torneo *TBinomial) PrintTablaNormalInterval(tabla Tabla) {
// 	t := table.NewWriter()
// 	t.SetOutputMirror(os.Stdout)
// 	t.AppendHeader(
// 		append(
// 			table.Row{"A\\B"},
// 			append(
// 				torneo.Inscriptos("l    u"),
// 				"B/A",
// 			)...,
// 		),
// 	)

// 	for i := 0; i < len(torneo.Agents); i++ {

// 		agent1 := torneo.Agents[i]
// 		row := []interface{}{
// 			agent1.UID(),
// 		}

// 		for j := 0; j < len(torneo.Agents); j++ {
// 			if i == j {
// 				row = append(row, "             ")
// 				continue
// 			}
// 			agent2 := torneo.Agents[j]
// 			u, l := tabla.NormalInterval(agent1.UID(), agent2.UID())
// 			row = append(row, fmt.Sprintf("%.3f %.3f", l, u))
// 		}

// 		row = append(row, agent1.UID())

// 		t.AppendRow(row)
// 		t.AppendSeparator()
// 	}

// 	// t.AppendFooter(table.Row{"", "", "Total", 10000})
// 	t.Render()
// }

func (torneo *TBinomial) Report() {
	if torneo.Cant_partidas_dobles > 0 {
		log.Printf("%d Partidas dobles:\n", torneo.Cant_partidas_dobles)

		log.Println()
		log.Printf("\nTABLA: WP (win percentage) & ADP (Avg. Diff. Points) para A vs B\n\n")
		torneo.PrintTablaMedia(torneo.Partidas)

		log.Println()
		log.Printf("\nTABLA: Intervalos Ajustados de Wald al 90%% para A vs B\n\n")
		torneo.PrintTablaWaldInterval(torneo.Partidas)

		// log.Println()
		// log.Printf("\nTABLA: Intervalos Normales al 90%% para A vs B\n\n")
		// torneo.PrintTablaNormalInterval(torneo.Partidas)
	}
}
