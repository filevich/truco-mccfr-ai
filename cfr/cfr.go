package cfr

import "github.com/filevich/truco-cfr/utils"

type CFR struct {
	*Trainer
}

func (trainer *CFR) String() string {
	return string(CFR_T) // "CFR"
}

func (trainer *CFR) regret_update_equation(

	t int,
	regret float32,
	cf_reach_prob float32,
	reg_acc float32,

) float32 {

	// vanilla cfr regret equation
	return reg_acc + (cf_reach_prob * regret)

}

func (trainer *CFR) strategy_update_equation(

	t int,
	reach_prob float32,
	action_prob float32, // ~ strategy[a]
	strategy_acc float32,

) float32 {

	// vanilla cfr strategy update equation
	return strategy_acc + (reach_prob * action_prob)

}

func (trainer *CFR) Train(profile IProfile) {
	profile.Init(trainer)

	for i := 0; i < profile.GetThreads(); i++ {
		go func() {
			// implementacion con simultaneous updates
			for ; profile.Continue(trainer); trainer.inc_t() {
				p := trainer.sample_partida()
				reach_probabilities := utils.Ones(trainer.get_num_players())
				acc := make([]float32, trainer.get_num_players())
				new_utils := _base_non_mc_run(trainer, profile, p, reach_probabilities, acc)
				trainer.add_root_utils(new_utils)
				profile.Check(trainer)
			}
			// join
			trainer.Wg.Done()
		}()
	}

	trainer.Wg.Wait()
	trainer.FinalReport(profile)
}
