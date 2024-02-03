package eval2

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/filevich/truco-cfr/abs"
	"github.com/filevich/truco-cfr/bot"
	"github.com/filevich/truco-cfr/eval/dataset"
	"github.com/filevich/truco-cfr/eval/dumbo"
	"github.com/filevich/truco-cfr/utils"
	"github.com/truquito/truco/pdt"
)

func PlayMultipleDoubleGames(

	agent bot.Agent,
	ops []bot.Agent,
	numPlayers int,
	ds dataset.Dataset,

) []*Results {

	agent.Initialize()
	res := make([]*Results, 0, len(ops))

	for _, op := range ops {
		op.Initialize()

		// evaluar_bin.go
		r := PlayDoubleGames(ds, agent, op, numPlayers)
		res = append(res, r)

		// termino de jugar contra op -> ya no lo necesito
		op.Free()
		runtime.GC()
	}

	return res
}

// partidas dobles/bin (hasta el final)
// la primera mitad el agent1 empieza primero
// la otra mitad el agent2 empieza primero
func PlayDoubleGames(

	ds dataset.Dataset,
	agent1,
	agent2 bot.Agent,
	numPlayers int,

) *Results {

	num_partidas := 2 * len(ds)
	start := time.Now()

	res := &Results{
		Title:              fmt.Sprintf("Double games %s vs %s", agent1.UID(), agent2.UID()),
		TotalNumberOfGames: num_partidas,
		WonByACounter:      0,
		Wons:               make([]float64, num_partidas/2),
		PointsWonDiff:      make([]float32, 0, num_partidas),
		Dumbo1:             0,
		Dumbo2:             0,
	}

	// IDA
	// partidas simples (hasta el final)
	for i := 0; i < num_partidas/2; i++ {
		entries := ds[i]
		agent1Won, diffPtsWonByAgent1Acc, d1, d2 := playSingleGame(entries, agent1, agent2, numPlayers)
		res.Dumbo1 += d1
		res.Dumbo2 += d2
		if agent1Won {
			res.WonByACounter++
			res.Wons[i] += 0.5
		}
		res.PointsWonDiff = append(
			res.PointsWonDiff,
			float32(diffPtsWonByAgent1Acc),
		)
	}

	// VUELTA
	// ahora los cambio de posicion
	for i := 0; i < num_partidas/2; i++ {
		entries := ds[i]
		agent2Won, diffPtsWonAgent2Acc, d2, d1 := playSingleGame(entries, agent2, agent1, numPlayers)
		res.Dumbo1 += d1
		res.Dumbo2 += d2
		if !agent2Won {
			res.WonByACounter++
			res.Wons[i] += 0.5
		}
		res.PointsWonDiff = append(
			res.PointsWonDiff,
			float32(-diffPtsWonAgent2Acc),
		)
	}

	res.Delta = time.Since(start)

	return res
}

func playSingleGame(

	rows []*dataset.Row,
	agent1,
	agent2 bot.Agent,
	Num_players int,

) (

	agent1Won bool,
	diffPtsWonByAgent1 float32,
	di1 int,
	di2 int,

) {

	limEnvite := 4
	verbose := true

	A, B := generateNames(agent1, agent2, Num_players)
	p, _ := pdt.NuevaPartida(pdt.A20, A, B, limEnvite, verbose)

	d1Total, d2Total := 0, 0

	for i := 0; !p.Terminada(); i++ {
		// empieza ronda nueva
		// en cada ronda correponde resetear los catchers
		agent1.ResetCatch()
		agent2.ResetCatch()
		rows[i].Override(p)
		_, _, d1, d2, _, _ := PlayRound(agent1, agent2, Num_players, p)
		d1Total += d1
		d2Total += d2
	}
	// termino la partida

	// EXTRAIDO de `jugar_partida_hasta_el_final`
	gent1Won := p.ElQueVaGanando() == pdt.Azul
	diffPtsWonByAgent1 = float32(p.Puntajes[pdt.Azul] - p.Puntajes[pdt.Rojo])

	return gent1Won, diffPtsWonByAgent1, d1Total, d2Total
}

func generateNames(agent1, agent2 bot.Agent, numPlayers int) (A, B []string) {

	for i := 0; i < numPlayers; i++ {
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

func PlayRound(

	azul,
	rojo bot.Agent,
	numPlayers int,
	p *pdt.Partida,

) (

	agent1Won bool,
	diffPtsAgent1,
	di1, di2 int,
	primeraJugada string,
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
		// 	abs := cfrA.Model.GetAbs()
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
		dumboid = dumbo.IsDumbo(p, m, a, &abs.Null{})

		if dumboid {
			if m.Jugador.Equipo == pdt.Azul {
				di1++
			} else {
				di2++
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
	agent1Won = p.Manojo(ganador).Jugador.Equipo == pdt.Azul
	diffPtsAgent1 = p.Puntajes[pdt.Azul] - p.Puntajes[pdt.Rojo]

	return agent1Won, diffPtsAgent1, di1, di2, primera.String(), prob
}
