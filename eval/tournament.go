package eval

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/filevich/truco-ai/bot"
	"github.com/filevich/truco-ai/eval/dataset"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Tournament struct {
	NumPlayers     int
	NumDoubleGames int
	Agents         []bot.Agent
	Table          Table
}

func (tournament *Tournament) Start(ds dataset.Dataset, verbose bool) {

	tournament.NumDoubleGames = len(ds)

	if verbose {
		log.Printf("Tournament %dp:\n", tournament.NumPlayers)
		for ix, agent := range tournament.Agents {
			log.Printf("\t%2d. %s\n", ix+1, agent.UID())
		}
	}

	for i := 0; i < len(tournament.Agents)-1; i++ {
		agent1 := tournament.Agents[i]
		ops := tournament.Agents[i+1:]

		ress := PlayMultipleDoubleGames(
			agent1,
			ops,
			tournament.NumPlayers,
			ds)

		for j := 0; j < len(ops); j++ {
			tournament.Table.Add(
				agent1.UID(),
				ops[j].UID(),
				ress[j],
			)
		}

		if verbose {
			log.Printf("%s", agent1.UID())
			if i == len(tournament.Agents)-2 {
				log.Printf("\n\n")
			} else {
				log.Printf(", ")
			}
		}

		// Just finished `agent1` evaluation
		// No longer needed; Free it
		agent1.Free()
		runtime.GC()
	}
}

func (tournament *Tournament) Registered(subtitle string) []interface{} {
	// "0.857   8.437" ~> tiene len=20
	res := make([]interface{}, 0)
	for _, agent := range tournament.Agents {
		res = append(res, fmt.Sprintf("%s\n  %s ", agent.UID(), subtitle))
	}
	return res
}

func (tournament *Tournament) PrintWrTable(tabla Table) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		append(
			table.Row{"A\\B"},
			append(
				tournament.Registered("WR    ADP"),
				"B/A",
			)...,
		),
	)

	for i := 0; i < len(tournament.Agents); i++ {

		agent1 := tournament.Agents[i]
		row := []interface{}{
			agent1.UID(),
		}

		for j := 0; j < len(tournament.Agents); j++ {
			if i == j {
				row = append(row, "             ")
				continue
			}
			agent2 := tournament.Agents[j]
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

func (tournament *Tournament) PrintWaldTable(tabla Table) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		append(
			table.Row{"A\\B"},
			append(
				tournament.Registered("L    U"),
				"B/A",
			)...,
		),
	)

	for i := 0; i < len(tournament.Agents); i++ {

		agent1 := tournament.Agents[i]
		row := []interface{}{
			agent1.UID(),
		}

		for j := 0; j < len(tournament.Agents); j++ {
			if i == j {
				row = append(row, "             ")
				continue
			}
			agent2 := tournament.Agents[j]
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

func (tournament *Tournament) Report() {
	if tournament.NumDoubleGames > 0 {
		log.Printf("%d Double games:\n", tournament.NumDoubleGames)

		log.Printf("TABLE: WR (Win Rate) & ADP (Avg. Diff. Points) for A vs B\n\n")
		tournament.PrintWrTable(tournament.Table)

		log.Printf("TABLE: Adjusted Wald Intervals (90%%) for A vs B\n\n")
		tournament.PrintWaldTable(tournament.Table)
	}
}
