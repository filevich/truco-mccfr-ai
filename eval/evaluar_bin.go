package eval

import (
	"fmt"
	"strconv"

	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/utils"
	"github.com/truquito/truco/pdt"
)

// partidas simples (hasta el final)
// la primera mitad el agent1 empieza primero
// la otra mitad el agent2 empieza primero
func SimPartidasBin(ds Dataset, agent1, agent2 bot.Agent, num_players int) *Resultados {

	num_partidas := 2 * len(ds)

	res := &Resultados{
		Titulo:          fmt.Sprintf("Partidas bin %s vs %s", agent1.UID(), agent2.UID()),
		Total:           num_partidas,
		Wins_by_A_count: 0,
		Wins:            make([]float64, num_partidas/2),
		Points_won_diff: make([]float32, 0),
		Dumbo1:          0,
		Dumbo2:          0,
	}

	// IDA
	// partidas simples (hasta el final)
	for i := 0; i < num_partidas/2; i++ {
		entries := ds[i]
		agent1_won, diff_pts_won_agent_1_acc, di_1, di_2 := partidaSegunEntradas(entries, agent1, agent2, num_players)
		res.Dumbo1 += di_1
		res.Dumbo2 += di_2
		if agent1_won {
			res.Wins_by_A_count++
			res.Wins[i] += 0.5
		}
		res.Points_won_diff = append(
			res.Points_won_diff,
			float32(diff_pts_won_agent_1_acc),
		)
	}

	// VUELTA
	// ahora los cambio de posicion
	for i := 0; i < num_partidas/2; i++ {
		entries := ds[i]
		agent2_won, diff_pts_won_agent_2_acc, di_2, di_1 := partidaSegunEntradas(entries, agent2, agent1, num_players)
		res.Dumbo1 += di_1
		res.Dumbo2 += di_2
		if !agent2_won {
			res.Wins[i] += 0.5
			res.Wins_by_A_count++
		}
		res.Points_won_diff = append(
			res.Points_won_diff,
			float32(-diff_pts_won_agent_2_acc),
		)
	}

	return res
}

func partidaSegunEntradas(

	entradas []*Row,
	agent1,
	agent2 bot.Agent,
	Num_players int,

) (

	agent1_won bool,
	diff_pts_won_agent_1_acc float32,
	di_1, di_2 int,

) {

	limEnvite := 4
	verbose := true

	// genero los nombres
	A, B := genNombres(agent1, agent2, Num_players)
	p, _ := pdt.NuevaPartida(pdt.A20, A, B, limEnvite, verbose)
	// p, _ := pdt.Parse(`{"puntuacion":20,"puntajes":{"Azul":0,"Rojo":0},"ronda":{"manoEnJuego":0,"cantJugadoresEnJuego":{"Azul":1,"Rojo":1},"elMano":0,"turno":0,"envite":{"estado":"noCantadoAun","puntaje":0,"cantadoPor":"","sinCantar":[]},"truco":{"cantadoPor":"","estado":"noCantado"},"manojos":[{"seFueAlMazo":false,"cartas":[{"palo":"Oro","valor":3},{"palo":"Espada","valor":2},{"palo":"Basto","valor":4}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Alice","equipo":"Azul"}},{"seFueAlMazo":false,"cartas":[{"palo":"Espada","valor":5},{"palo":"Copa","valor":2},{"palo":"Basto","valor":11}],"tiradas":[false,false,false],"ultimaTirada":0,"jugador":{"id":"Bob","equipo":"Rojo"}}],"mixs":{"Alice":0,"Bob":1},"muestra":{"palo":"Basto","valor":2},"manos":[{"resultado":"ganoRojo","ganador":"","cartasTiradas":null},{"resultado":"ganoRojo","ganador":"","cartasTiradas":null},{"resultado":"ganoRojo","ganador":"","cartasTiradas":null}]}}`)

	d1_acc, d2_acc := 0, 0

	for i := 0; !p.Terminada(); i++ {
		// empieza ronda nueva
		// en cada ronda correponde resetear los catchers
		agent1.ResetCatch()
		agent2.ResetCatch()
		entradas[i].Override(p)
		_, _, di1, di2, _, _ := JugarRondaHastaElFinal(agent1, agent2, Num_players, p)
		d1_acc += di1
		d2_acc += di2
	}
	// termino la partida

	// EXTRAIDO de `jugar_partida_hasta_el_final`
	gano_agent1 := p.ElQueVaGanando() == pdt.Azul
	diff_pts_agent1 := p.Puntajes[pdt.Azul] - p.Puntajes[pdt.Rojo]

	return gano_agent1, float32(diff_pts_agent1), d1_acc, d2_acc
}

func genNombres(agent1, agent2 bot.Agent, Num_players int) (A, B []string) {

	// genero los nombres
	// A := []string{} // equipo azul
	// B := []string{} // equipo rojo

	for i := 0; i < Num_players; i++ {
		if utils.Mod(i, 2) == 0 {
			given_name := agent1.UID() + strconv.Itoa(i+1)
			A = append(A, given_name)
		} else {
			given_name := agent2.UID() + strconv.Itoa(i+1)
			B = append(B, given_name)
		}
	}

	return A, B
}

// PRE: agent ya esta inicializado
// POST: ops seran incializados
func SingleSimPartidasBin(

	agent bot.Agent,
	ops []bot.Agent,
	num_players int,
	ds Dataset,

) []*Resultados {

	res := make([]*Resultados, len(ops))

	for i := 0; i < len(ops); i++ {
		op := ops[i]
		op.Initialize()
		res_partidas := SimPartidasBin(ds, agent, op, num_players)
		res[i] = res_partidas
	}

	return res
}
