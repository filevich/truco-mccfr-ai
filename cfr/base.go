package cfr

import (
	"crypto/sha1"

	"github.com/filevich/truco-ai/info"
	"github.com/filevich/truco-ai/utils"
	"github.com/truquito/truco/pdt"
)

// implementacion con simultaneous updates
func _baseNonMcRun(

	trainer ITrainer,
	profile IProfile,
	p *pdt.Partida,
	reach_probabilities []float32,
	acc []float32,

) []float32 {

	// pseudo jugador activo
	elMano := p.Ronda.GetElMano().Jugador.ID
	active_player := pdt.Rho(p)
	rix := utils.RIX(p, active_player) // su indice relativo al mano
	rix_mod2 := utils.Mod(rix, 2)      // lo hago 2p. (0 ~ los que son team el-mano, 1 ~ si no)

	// obtengo el infoset
	aixs := pdt.GetA(p, active_player)
	// i := MkInfoset1(p, active_player, aixs, trainer.GetAbs())
	i := info.NewInfosetRondaBase(p, active_player, trainer.GetAbs(), nil)
	hash, chi_len := i.Hash(sha1.New()), i.ChiLen()

	// obtengo el RNode
	rnode := trainer.GetRnode(hash, chi_len)
	// capaz que ahora ya cambio la strategy
	trainer.Lock()
	strategy := rnode.GetStrategy()
	trainer.Unlock()
	counterfactual_values := make([][]float32, chi_len)

	// obtengo el chi
	A := i.Iterable(p, active_player, aixs, trainer.GetAbs())
	bs, _ := p.MarshalJSON()

	for aix, j := range A {

		// prunning
		skip := profile.IsPrunable(trainer) && strategy[aix] < 0.01 // menor a 1% de prob
		if skip {
			counterfactual_values[aix] = make([]float32, trainer.getNumPlayers())
			continue
		}

		action_probability := strategy[aix]
		new_reach_probabilities := make([]float32, len(reach_probabilities))
		copy(new_reach_probabilities, reach_probabilities)
		new_reach_probabilities[rix_mod2] *= action_probability

		// ejecuto la accion
		p, _ = pdt.Parse(string(bs), true)
		pkts := j.Hacer(p)

		// hemos llegado a un nodo terminal ?
		termino, pts_ganados, ganador := utils.IsDoneAndPts(pkts)
		if termino {

			// fin de la ronda, fin de la recursion:
			// no hace falta que vuelva a llamar recursivamente a cfr
			// ya se lo que deberia devolver
			new_pts := utils.Payoffs(p.Manojo(elMano), pts_ganados, p.Manojo(ganador))
			new_acc := utils.SumFloat32Slices(acc, new_pts)
			counterfactual_values[aix] = new_acc

		} else {
			if pts_ganados > 0 {

				// acumulo los puntos (del envite)
				new_pts := utils.Payoffs(p.Manojo(elMano), pts_ganados, p.Manojo(ganador))
				new_acc := utils.SumFloat32Slices(acc, new_pts)
				counterfactual_values[aix] = _baseNonMcRun(trainer, profile, p, new_reach_probabilities, new_acc)

			} else {

				// ni hay puntos nuevos, ni termino -> paso el acc intacto
				counterfactual_values[aix] = _baseNonMcRun(trainer, profile, p, new_reach_probabilities, acc)

			}
		}

	}

	node_values := utils.Ndot(strategy, counterfactual_values)

	trainer.Lock()
	rnode.RegUpdates++
	rnode.StrUpdates++
	t := rnode.RegUpdates // alt: trainer.Get_t()+1
	trainer.Unlock()

	for aix := range A {

		// prunning
		prunning := profile.IsPrunable(trainer) && strategy[aix] < 0.01 // menor a 1% de prob
		if prunning {
			continue
		}

		// actualizacion de regrets
		cf_reach_prob := counterfactualReachProbability(reach_probabilities, rix_mod2)
		regret := counterfactual_values[aix][rix_mod2] - node_values[rix_mod2]

		trainer.Lock() // abro qui, cierro en linea 127

		rnode.CumulativeRegrets[aix] = trainer.regretUpdateEquation(
			t,
			regret,
			cf_reach_prob,
			rnode.CumulativeRegrets[aix],
		)

		// acumulacion de strategy
		rnode.StrategySum[aix] = trainer.strategyUpdateEquation(
			t,                             // iter actual
			reach_probabilities[rix_mod2], // reach_prob
			strategy[aix],                 // P(a)
			rnode.StrategySum[aix],        // strategy_acc
		)

		trainer.Unlock() // cierro aqui lo que abri en linea 110

	}

	return node_values
}
