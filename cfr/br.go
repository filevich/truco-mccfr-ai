package cfr

import (
	"encoding/json"

	"github.com/filevich/truco-mccfr-ai/utils"

	"github.com/truquito/gotruco/pdt"
)

type BestResponse struct {
	*Trainer
}

func (trainer *BestResponse) String() string {
	return string(BR_T)
}

func (trainer *BestResponse) regretUpdateEquation(
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

func (trainer *BestResponse) strategyUpdateEquation(
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

func (trainer *BestResponse) Train(profile IProfile) {
	profile.Init(trainer)

	for i := 0; i < profile.GetThreads(); i++ {
		go func() {
			for ; profile.Continue(trainer); trainer.inc_t() {
				// genera una partida Alice-Bob-Ariana-Ben (segun sea 1vs1 o 2vs2)
				// inicialmente:
				// 		jugador rix_mod2:0 ~ lo asocio a i=0
				// 		jugador rix_mod2:1 ~ lo asocio a i=1
				p := trainer.samplePartida()
				bs, _ := json.Marshal(p)

				// 2 veces; una ida y una vuelta
				// for i := 0; i < 2; i++ {
				for i := 0; i < 1; i++ {
					acc := make([]float32, trainer.getNumPlayers())
					var new_utils []float32

					// ojo: inicialmente tengo que asociar el jugador 0 con el modelo, 1 con agent

					if i == 0 {
						new_utils = trainer.run(profile, p, acc, 0)
					} else {
						p, _ = pdt.Parse(string(bs), true)
						// cambio las posiciones
						// p.Swap()
						// p.Ronda.Reset(0)
						new_utils = trainer.run(profile, p, acc, 1)
					}
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

func (trainer *BestResponse) run(
	profile IProfile,
	p *pdt.Partida,
	acc []float32,
	update_player int,
) []float32 {
	// pseudo jugador activo en el TRUCO
	elMano := p.Ronda.GetElMano().Jugador.ID
	active_player := pdt.Rho(p)
	rix := utils.RIX(p, active_player) // su indice relativo al mano
	rix_mod2 := utils.Mod(rix, 2)      // lo hago 2p. (0 ~ los que son team el-mano, 1 ~ si no)

	esTurnoDelAgente := rix_mod2 != update_player
	if esTurnoDelAgente {
		inGameID := active_player.Jugador.ID
		a, _ := profile.Exploit().Action(p, inGameID)
		pkts := a.Hacer(p)
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
	}

	// obtengo el infoset
	aixs := pdt.GetA(p, active_player)
	i := trainer.GetBuilder().Info(p, active_player, nil)
	hash, chi_len := i.Hash(trainer.GetBuilder().Hash()), i.ChiLen()

	// obtengo el RNode
	rnode := trainer.GetRnode(hash, chi_len)
	trainer.Lock()
	strategy := rnode.GetStrategy()
	trainer.Unlock()

	// obtengo el chi
	A := i.Iterable(p, active_player, aixs, trainer.GetAbs())
	bs, _ := p.MarshalJSON()

	// solo en este caso elige una estrategia al azar

	// acumulo la strategy
	trainer.Lock()
	rnode.StrUpdates++
	iter := rnode.StrUpdates

	// acumulacion de strategy (el vector entero)
	for aix := range A {
		rnode.StrategySum[aix] = trainer.strategyUpdateEquation(
			iter,                   // iter actual
			-1,                     // reach_prob <- ya no importa
			strategy[aix],          // P(a)
			rnode.StrategySum[aix], // strategy_acc
		)
	}
	trainer.Unlock()

	// sino todo "igual" que antes:
	counterfactual_values := make([][]float32, chi_len)

	for aix, j := range A {
		// prunning
		skip := profile.IsPrunable(trainer, strategy[aix]) // menor a 1% de prob
		if skip {
			counterfactual_values[aix] = make([]float32, trainer.getNumPlayers())
			continue
		}

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
				counterfactual_values[aix] = trainer.run(profile, p, new_acc, update_player)
			} else {
				// ni hay puntos nuevos, ni termino -> paso el acc intacto
				counterfactual_values[aix] = trainer.run(profile, p, acc, update_player)
			}
		}
	}

	node_values := utils.Ndot(strategy, counterfactual_values)

	// actualizo los regrets
	trainer.Lock()
	rnode.RegUpdates++
	t := rnode.RegUpdates // alt: trainer.Get_t()+1

	for aix := range A {
		// prunning
		prunning := profile.IsPrunable(trainer, strategy[aix]) // menor a 1% de prob
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
