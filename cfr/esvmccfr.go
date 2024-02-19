package cfr

import (
	"encoding/json"

	"github.com/filevich/truco-ai/utils"
	"github.com/truquito/truco/pdt"
)

type ESVMCCFR struct {
	*Trainer
}

func (trainer *ESVMCCFR) String() string {
	return string(ESVMCCFR_T) // "ESVMCCFR"
}

func (trainer *ESVMCCFR) regretUpdateEquation(

	t int,
	regret float32,
	cf_reach_prob float32,
	reg_acc float32,

) float32 {

	// vanilla cfr regret equation
	return reg_acc + (cf_reach_prob * regret)

}

func (trainer *ESVMCCFR) strategyUpdateEquation(

	t int,
	reach_prob float32,
	action_prob float32, // ~ strategy[a]
	strategy_acc float32,

) float32 {

	// vanilla cfr strategy update equation
	return strategy_acc + (reach_prob * action_prob)

}

func (trainer *ESVMCCFR) Train(profile IProfile) {
	profile.Init(trainer)

	for i := 0; i < profile.GetThreads(); i++ {
		go func() {
			// implementacion con simultaneous updates
			for ; profile.Continue(trainer); trainer.inc_t() {
				p := trainer.samplePartida()
				bs, _ := json.Marshal(p)
				for i := 0; i < 2; i++ {
					if i > 0 {
						p, _ = pdt.Parse(string(bs), true)
					}
					acc := make([]float32, trainer.getNumPlayers())
					reach_probabilities := utils.Ones(trainer.getNumPlayers())
					new_utils := trainer.run(profile, p, reach_probabilities, acc, i)
					trainer.addRootUtils(new_utils)
				}
				profile.Check(trainer)
			}
			// join
			trainer.Wg.Done()
		}()
	}

	trainer.Wg.Wait()
	trainer.FinalReport(profile)
}

func (trainer *ESVMCCFR) run(

	profile IProfile,
	p *pdt.Partida,
	reach_probabilities []float32,
	acc []float32,
	update_player int,

) []float32 {

	// pseudo jugador activo
	elMano := p.Ronda.GetElMano().Jugador.ID
	active_player := pdt.Rho(p)
	rix := utils.RIX(p, active_player) // su indice relativo al mano
	rix_mod2 := utils.Mod(rix, 2)      // lo hago 2p. (0 ~ los que son team el-mano, 1 ~ si no)

	// obtengo el infoset
	aixs := pdt.GetA(p, active_player)
	// i := MkInfoset1(p, active_player, aixs, trainer.GetAbs())
	// i := info.NewInfosetRondaBase(p, active_player, trainer.GetAbs(), nil)
	i := trainer.GetBuilder().Info(p, active_player, nil)
	hash, chi_len := i.Hash(trainer.GetBuilder().Hash), i.ChiLen()

	// obtengo el RNode
	rnode := trainer.GetRnode(hash, chi_len)
	trainer.Lock()
	strategy := rnode.GetStrategy()
	trainer.Unlock()
	raix := utils.Sample(strategy)

	// reach
	action_probability := strategy[raix]
	new_reach_probabilities := make([]float32, len(reach_probabilities))
	copy(new_reach_probabilities, reach_probabilities)
	new_reach_probabilities[rix_mod2] *= action_probability

	// obtengo el chi
	A := i.Iterable(p, active_player, aixs, trainer.GetAbs())
	bs, _ := p.MarshalJSON()

	// no es el `update_player` -> actualizo solo la estrategia
	if rix_mod2 != update_player {

		// acumulo la strategy (solo a raix ?)
		trainer.Lock()
		rnode.StrUpdates++
		t := rnode.StrUpdates

		// acumulacion de strategy (el vector entero)
		for aix := range A {

			// reach prob.
			reach_prob := reach_probabilities[rix_mod2]

			rnode.StrategySum[aix] = trainer.strategyUpdateEquation(
				t,                      // iter actual
				reach_prob,             // reach_prob
				strategy[aix],          // P(a)
				rnode.StrategySum[aix], // strategy_acc
			)
		}
		trainer.Unlock()

		////////////////////////////////////////////////////////////

		pkts := A[raix].Hacer(p)
		termino, pts_ganados, ganador := utils.IsDoneAndPts(pkts)
		if termino {

			// fin de la ronda, fin de la recursion:
			// no hace falta que vuelva a llamar recursivamente a cfr
			// ya se lo que deberia devolver
			new_pts := utils.Payoffs(p.Manojo(elMano), pts_ganados, p.Manojo(ganador))
			new_acc := utils.SumFloat32Slices(acc, new_pts)
			return new_acc

		} else {

			if pts_ganados > 0 {

				// acumulo los puntos (del envite)
				new_pts := utils.Payoffs(p.Manojo(elMano), pts_ganados, p.Manojo(ganador))
				new_acc := utils.SumFloat32Slices(acc, new_pts)
				return trainer.run(profile, p, new_reach_probabilities, new_acc, update_player)

			} else {

				// ni hay puntos nuevos, ni termino -> paso el acc intacto
				return trainer.run(profile, p, new_reach_probabilities, acc, update_player)

			}
		}

		//////////////////////////////////////////////////////////////
	}

	// sino todo "igual" que antes:
	counterfactual_values := make([][]float32, chi_len)

	for aix, j := range A {

		// prunning
		skip := profile.IsPrunable(trainer, strategy[aix])
		if skip {
			counterfactual_values[aix] = make([]float32, trainer.getNumPlayers())
			continue
		}

		// ejecuto la accion
		p, _ = pdt.Parse(string(bs), true)
		pkts := j.Hacer(p)

		////////////////////////////////////////////////////////////

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
				counterfactual_values[aix] = trainer.run(
					profile,
					p,
					new_reach_probabilities,
					new_acc,
					update_player,
				)

			} else {

				// ni hay puntos nuevos, ni termino -> paso el acc intacto
				counterfactual_values[aix] = trainer.run(
					profile,
					p,
					new_reach_probabilities,
					acc,
					update_player,
				)

			}
		}

		////////////////////////////////////////////////////////////

	}

	node_values := utils.Ndot(strategy, counterfactual_values)

	// actualizo los regrets
	trainer.Lock()
	rnode.RegUpdates++
	t := rnode.RegUpdates // alt: trainer.Get_t()+1

	for aix := range A {

		// prunning
		prunning := profile.IsPrunable(trainer, strategy[aix])
		if prunning {
			continue
		}

		// actualizacion de regrets
		regret := counterfactual_values[aix][rix_mod2] - node_values[rix_mod2]

		rnode.CumulativeRegrets[aix] = trainer.regretUpdateEquation(
			t,
			regret,
			new_reach_probabilities[rix_mod2],
			rnode.CumulativeRegrets[aix],
		)
	}

	trainer.Unlock()

	return node_values

}
