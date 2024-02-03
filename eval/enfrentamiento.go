package eval

import (
	"encoding/json"
	"strconv"

	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/utils"
	"github.com/truquito/truco/pdt"
)

func rondaDoble(

	agent1,
	agent2 bot.Agent,
	Num_players int,

) (

	agent1_win_count int,
	diff_pts_won_agent_1_acc float32,
	dumbo1, dumbo2 int,

) {

	limEnvite := 2
	verbose := true

	// genero los nombres
	A := []string{} // equipo azul
	B := []string{} // equipo rojo

	for i := 0; i < Num_players; i++ {
		if utils.Mod(i, 2) == 0 {
			given_name := agent1.UID() + strconv.Itoa(i+1)
			A = append(A, given_name)
		} else {
			given_name := agent2.UID() + strconv.Itoa(i+1)
			B = append(B, given_name)
		}
	}

	gano_agent1_count := 0
	diff_pts_agent1_acc := 0
	dumbo_agent1 := 0
	dumbo_agent2 := 0

	p, _ := pdt.NuevaPartida(pdt.A20, A, B, limEnvite, verbose)

	// guardo la config para luego shiftear la posicion de los jugadores
	data_muestra, _ := json.Marshal(p.Ronda.Muestra)
	data_manojos, _ := json.Marshal(p.Ronda.Manojos)

	// 1. ronda de ida
	if gano_agent1, diff_pts, di_1, di_2, _, _ := JugarRondaHastaElFinal(agent1, agent2, Num_players, p); gano_agent1 {
		gano_agent1_count += 1
		diff_pts_agent1_acc += diff_pts
		dumbo_agent1 += di_1
		dumbo_agent2 += di_2
	} else {
		// gano_agent1_count -= 1 // <- no resto! sino no es el porcentaje
		diff_pts_agent1_acc += diff_pts
		dumbo_agent1 += di_1
		dumbo_agent2 += di_2
	}

	// 2. ronda de vuelta
	// invierto la posicion de los jugadores
	p, _ = pdt.NuevaPartida(pdt.A20, B, A, limEnvite, verbose)

	var (
		manojos []pdt.Manojo
		muestra pdt.Carta
	)

	json.Unmarshal(data_muestra, &muestra)
	json.Unmarshal(data_manojos, &manojos)

	p.Ronda.SetManojos(manojos)
	p.Ronda.SetMuestra(muestra)

	if gano_agent2, diff_pts, di_2, di_1, _, _ := JugarRondaHastaElFinal(agent2, agent1, Num_players, p); gano_agent2 {
		// gano_agent1_count -= 1 // no resto! sino no es el porcentaje
		diff_pts_agent1_acc -= diff_pts
		dumbo_agent1 += di_1
		dumbo_agent2 += di_2
	} else {
		gano_agent1_count += 1
		diff_pts_agent1_acc -= diff_pts
		dumbo_agent1 += di_1
		dumbo_agent2 += di_2
	}

	return gano_agent1_count, float32(diff_pts_agent1_acc), dumbo_agent1, dumbo_agent2
	// lo divido entre 2 porque fueron 2 partidos: el de ida y el de vuelta
}

func partidaUnica(

	agent1,
	agent2 bot.Agent,
	Num_players int,

) (agent1_won bool, diff_pts_won_agent_1_acc float32) {

	limEnvite := 2
	verbose := true

	// genero los nombres
	A := []string{} // equipo azul
	B := []string{} // equipo rojo

	for i := 0; i < Num_players; i++ {
		if utils.Mod(i, 2) == 0 {
			given_name := agent1.UID() + strconv.Itoa(i+1)
			A = append(A, given_name)
		} else {
			given_name := agent2.UID() + strconv.Itoa(i+1)
			B = append(B, given_name)
		}
	}

	p, _ := pdt.NuevaPartida(pdt.A20, A, B, limEnvite, verbose)

	agent1_won, diff_pts := JugarPartidaHastaElFinal(agent1, agent2, Num_players, p)

	return agent1_won, float32(diff_pts)
}
