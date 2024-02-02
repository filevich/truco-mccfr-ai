package eval

import (
	"fmt"

	"github.com/filevich/truco-cfr/bot"
)

// partidas simples (hasta el final)
// la primera mitad el agent1 empieza primero
// la otra mitad el agent2 empieza primero
func SimPartidas(num_partidas int, agent1, agent2 bot.Agent, num_players int) *Resultados {

	res := &Resultados{
		Titulo:          fmt.Sprintf("Partidas únicas %s vs %s", agent1.UID(), agent2.UID()),
		Total:           num_partidas,
		Wins_by_A_count: 0,
		Points_won_diff: make([]float32, 0),
	}

	// partidas simples (hasta el final)
	for i := 0; i < num_partidas/2; i++ {
		agent1_won, diff_pts_won_agent_1_acc := partida_unica(agent1, agent2, 2)
		if agent1_won {
			res.Wins_by_A_count++
		}
		res.Points_won_diff = append(
			res.Points_won_diff,
			float32(diff_pts_won_agent_1_acc),
		)
	}

	// ahora los cambio de posicion
	for i := 0; i < num_partidas/2; i++ {
		agent2_won, diff_pts_won_agent_2_acc := partida_unica(agent2, agent1, 2)
		if !agent2_won {
			res.Wins_by_A_count++
		}
		res.Points_won_diff = append(
			res.Points_won_diff,
			float32(-diff_pts_won_agent_2_acc),
		)
	}

	return res
}

// rondas dobles (hasta el primer `Pkt` de `NuevaRonda`)
// la primera mitad el agent1 empieza primero
// la otra mitad el agent2 empieza primero
// es exactamente la misma ronda
func SimRondas(num_rondas int, agent1, agent2 bot.Agent, Num_players int) *Resultados {

	res := &Resultados{
		Titulo:          fmt.Sprintf("Rondas dobles %s vs %s", agent1.UID(), agent2.UID()),
		Total:           num_rondas * 2, // porque sino da 200% de "ganacion" cuando deberia ser 100% (ya que se juga a doble ronda: ida-y-vuelta)
		Wins_by_A_count: 0,
		Points_won_diff: make([]float32, 0),
	}

	// rondas dobles
	for i := 0; i < num_rondas; i++ {
		agent1_win_count, diff_pts_won_agent_1_acc, di_1, di_2 := ronda_doble(agent1, agent2, 2)
		res.Wins_by_A_count += agent1_win_count // {2,1,0}
		// res.points_won_delta[i] = diff_pts_won_agent_1_acc
		res.Points_won_diff = append(res.Points_won_diff, diff_pts_won_agent_1_acc)
		res.Dumbo1 += di_1
		res.Dumbo2 += di_2
	}

	return res
}
