package cfr

import (
	"crypto/sha1"
	"encoding/json"

	"github.com/filevich/truco-ai/info"
	"github.com/filevich/truco-ai/utils"

	"github.com/truquito/truco/pdt"
)

type ESLMCCFR struct {
	*Trainer
}

func (trainer *ESLMCCFR) String() string {
	return string(ESLMCCFR_T) // "ESLMCCFR"
}

func (trainer *ESLMCCFR) regretUpdateEquation(

	t int,
	regret float32,
	cf_reach_prob float32,
	reg_acc float32,

) float32 {

	// LCFR strategy update equation se diferencia de vanilla en cuanto a que
	// pondera las actualizaciones
	// pero a diferencia de cfr+ NO ignora los regrets negativos
	z := float32(t)
	weight := z / (z + 1.0)
	return (weight * reg_acc) + (1 * regret)

}

func (trainer *ESLMCCFR) strategyUpdateEquation(

	t int,
	reach_prob float32,
	action_prob float32, // ~ strategy[a]
	strategy_acc float32,

) float32 {

	// LCFR strategy update equation se diferencia de vanilla en cuanto a que
	// pondera las actualizaciones
	z := float32(t)
	weight := z / (z + 1.0)
	return (weight * strategy_acc) + (1 * action_prob)

}

func (trainer *ESLMCCFR) Train(profile IProfile) {
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
					new_utils := trainer.run(profile, p, acc, i)
					trainer.add_root_utils(new_utils)
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

func (trainer *ESLMCCFR) run(

	profile IProfile,
	p *pdt.Partida,
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
	// i := MkInfoset1(p, active_player, aixs, trainer.Get_abs())
	// hash, chi_len := i.Hash(), i.Chi_len()
	i := info.NewInfosetRondaBase(p, active_player, trainer.GetAbs(), nil)
	hash, chi_len := i.Hash(sha1.New()), i.ChiLen()

	// obtengo el RNode
	rnode := trainer.GetRnode(hash, chi_len)
	trainer.Lock()
	strategy := rnode.GetStrategy()
	trainer.Unlock()
	raix := utils.Sample(strategy)

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
			rnode.StrategySum[aix] = trainer.strategyUpdateEquation(
				t,                      // iter actual
				-1,                     // reach_prob <- ya no importa
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
				return trainer.run(profile, p, new_acc, update_player)

			} else {

				// ni hay puntos nuevos, ni termino -> paso el acc intacto
				return trainer.run(profile, p, acc, update_player)

			}
		}

		//////////////////////////////////////////////////////////////

		// return aqui
	}

	// sino todo "igual" que antes:
	counterfactual_values := make([][]float32, chi_len)

	for aix, j := range A {

		// prunning
		skip := profile.IsPrunable(trainer) && strategy[aix] < 0.01 // menor a 1% de prob
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
				counterfactual_values[aix] = trainer.run(profile, p, new_acc, update_player)

			} else {

				// ni hay puntos nuevos, ni termino -> paso el acc intacto
				counterfactual_values[aix] = trainer.run(profile, p, acc, update_player)

			}
		}

		////////////////////////////////////////////////////////////

	}

	// solo de debug:
	// if len(strategy) == 0 || len(counterfactual_values) == 0 {
	// 	fmt.Println(strategy, counterfactual_values)
	// 	fmt.Println(pdt.Renderizar(p))
	// 	aixs := pdt.GetA(p, active_player)
	// 	i := MkInfoset1(p, active_player, aixs, trainer.Get_abs())
	// 	hash, chi_len := i.Hash(), i.Chi_len()
	// 	fmt.Println(hash, chi_len)
	// }
	node_values := utils.Ndot(strategy, counterfactual_values)

	// actualizo los regrets
	trainer.Lock()
	rnode.RegUpdates++
	t := rnode.RegUpdates // alt: trainer.Get_t()+1

	for aix := range A {

		// prunning
		prunning := profile.IsPrunable(trainer) && strategy[aix] < 0.01 // menor a 1% de prob
		if prunning {
			continue
		}

		// actualizacion de regrets
		regret := counterfactual_values[aix][rix_mod2] - node_values[rix_mod2]

		rnode.CumulativeRegrets[aix] = trainer.regretUpdateEquation(
			t,
			regret,
			-1, // <- ya no hay CR_prob
			rnode.CumulativeRegrets[aix],
		)
	}

	trainer.Unlock()

	return node_values

}
