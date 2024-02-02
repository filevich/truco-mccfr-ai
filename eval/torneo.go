package eval

import (
	"fmt"
	"os"

	"github.com/filevich/truco-cfr/bot"
	"github.com/jedib0t/go-pretty/v6/table"
)

type Torneo struct {
	Cant_partidas_simples int
	Cant_rondas_dobles    int
	Agents                []bot.Agent
	Partidas              Tabla
	Rondas                Tabla
}

func (t *Torneo) Start(cant_partidas_simples, cant_rondas_dobles int) {

	t.Cant_partidas_simples = cant_partidas_simples
	t.Cant_rondas_dobles = cant_rondas_dobles

	fmt.Println("\nTorneo:")
	for ix, agent := range t.Agents {
		fmt.Printf("\t%2d. %s\n", ix+1, agent.UID())
	}

	fmt.Printf("\nDone: ")

	for i := 0; i < len(t.Agents)-1; i++ {
		agent1 := t.Agents[i]
		for j := i + 1; j < len(t.Agents); j++ {
			agent2 := t.Agents[j]

			res_partidas := SimPartidas(cant_partidas_simples, agent1, agent2, 2)
			t.Partidas.Add(agent1.UID(), agent2.UID(), res_partidas)

			if cant_rondas_dobles > 0 {
				res_rondas := SimRondas(cant_rondas_dobles, agent1, agent2, 2)
				t.Rondas.Add(agent1.UID(), agent2.UID(), res_rondas)
			}

		}
		fmt.Printf("%s", agent1.UID())
		if i == len(t.Agents)-2 {
			fmt.Printf("\n\n")
		} else {
			fmt.Printf(", ")
		}
	}
}

func (t *Torneo) Inscriptos() []interface{} {
	// "0.857   8.437" ~> tiene len=20
	res := make([]interface{}, 0)
	for _, agent := range t.Agents {
		res = append(res, fmt.Sprintf("%s\n  WP    ADP ", agent.UID()))
	}
	return res
}

func (torneo *Torneo) PrintTabla(tabla Tabla) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(
		append(
			table.Row{"A\\B"},
			append(
				torneo.Inscriptos(),
				"A\\B",
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

func (torneo *Torneo) Report() {
	if torneo.Cant_partidas_simples > 0 {
		fmt.Printf("%d Partidas simples:\n", torneo.Cant_partidas_simples)
		torneo.PrintTabla(torneo.Partidas)
	}

	if torneo.Cant_rondas_dobles > 0 {
		fmt.Printf("\n%d Rondas simples:\n", torneo.Cant_rondas_dobles)
		torneo.PrintTabla(torneo.Rondas)
	}
}
