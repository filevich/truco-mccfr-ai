package eval

import (
	"github.com/filevich/truco-cfr/abs"
	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/utils"
	"github.com/truquito/truco/pdt"
)

func Jugar_partida_hasta_el_final(

	agent1,
	agent2 bot.Agent,
	Num_players int,
	p *pdt.Partida,

) (gano_agent1 bool, diff_pts_agent1 int) {

	// importante
	// "agent1 siempre es equipo azul"
	// "agent2 siempre es equipo rojo"

	// juega hasta terminar la partida, no la ronda
	for !p.Terminada() {

		// log.Println(pdt.Renderizar(p))

		m := pdt.Rho(p)
		agent := agent1
		if m.Jugador.Equipo == pdt.Rojo {
			agent = agent2
		}
		a, _ := agent.Action(p, m.Jugador.ID)

		// log.Println(a)

		_, ok := a.Ok(p)
		if !ok {
			aa, _ := agent.Action(p, m.Jugador.ID)
			aa.Ok(p)
			panic("LA ACCION NO ES VALIDA")
		}

		// pkts := a.Hacer(p)
		a.Hacer(p)

	}

	// termino la partida
	// log.Println(pdt.Renderizar(p))
	gano_agent1 = p.ElQueVaGanando() == pdt.Azul
	diff_pts_agent1 = p.Puntajes[pdt.Azul] - p.Puntajes[pdt.Rojo]
	// log.Println(gano_agent1, diff_pts_agent1)

	return gano_agent1, diff_pts_agent1
}

func Jugar_ronda_hasta_el_final(

	azul,
	rojo bot.Agent,
	Num_players int,
	p *pdt.Partida,

) (

	gano_agent1 bool,
	diff_pts_agent1,
	di_1, di_2 int,
	primera_jugada string,
	prob float32,

) {

	// importante
	// "agent1 siempre es equipo azul"
	// "agent2 siempre es equipo rojo"

	var (
		termino_ronda bool        = false
		ganador       string      = ""
		primera       pdt.IJugada = nil
	)

	prob = 1.0

	// juega hasta terminar la partida, no la ronda
	for !termino_ronda {

		m := pdt.Rho(p)
		agent := azul
		if m.Jugador.Equipo == pdt.Rojo {
			agent = rojo
		}
		a, prob_a := agent.Action(p, m.Jugador.ID)

		if primera == nil {
			primera = a
		}

		prob *= prob_a

		// dumbo ?
		dumboid := false
		// if cfrA, isCFR := agent.(*BotCFR); isCFR {
		// 	abs := cfrA.Model.Get_abs()
		// 	dumboid = Is_dumbo(p, m, a, abs)
		// } else {
		// 	// _, isAle := agent.(*BotAleatorio)
		// 	// _, isDet := agent.(*BotDeterminista)
		// 	// _, isDet2 := agent.(*BotDeterministaMax)
		// 	// if isAle || isDet || isDet2 {
		// 	// 	dumboid = is_dumbo(p, m, a, cfr.Z0{})
		// 	// }
		// 	dumboid = Is_dumbo(p, m, a, cfr.Z0{})
		// }
		dumboid = Is_dumbo(p, m, a, &abs.Null{})

		if dumboid {
			if m.Jugador.Equipo == pdt.Azul {
				di_1++
			} else {
				di_2++
			}
		}

		_, ok := a.Ok(p)
		if !ok {
			aa, _ := agent.Action(p, m.Jugador.ID)
			// bs, _ := json.Marshal(p)
			// log.Println(string(bs))
			// log.Println("-------------")
			// log.Println(a)
			aa.Ok(p)
			panic("LA ACCION NO ES VALIDA")
		}

		aunEnPrimeraMano := p.Ronda.ManoEnJuego == pdt.Primera
		pkts := a.Hacer(p)
		if aunEnPrimeraMano {
			azul.Catch(p, pkts)
			rojo.Catch(p, pkts)
		}

		termino_ronda, _, ganador = utils.IsDoneAndPts(pkts)
	}

	// termino la ronda
	// log.Println(pdt.Renderizar(p))
	gano_agent1 = p.Manojo(ganador).Jugador.Equipo == pdt.Azul
	diff_pts_agent1 = p.Puntajes[pdt.Azul] - p.Puntajes[pdt.Rojo]

	return gano_agent1, diff_pts_agent1, di_1, di_2, primera.String(), prob
}
