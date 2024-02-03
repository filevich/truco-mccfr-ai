package cfr

import (
	"github.com/filevich/truco-cfr/abs"
	"github.com/truquito/truco/pdt"
)

type ITrainer interface {
	String() string
	Train(IProfile)
	Reset()

	Get_t() int
	inc_t()

	get_T() int
	set_T(int)
	inc_T()

	Count_infosets() int
	get_num_players() int
	Get_abs() abs.IAbstraction
	Get_rnode(hash string, chi_len int) *RNode
	sample_partida() *pdt.Partida
	Max_Avg_Game_Value() float32

	// multi
	Lock()
	Unlock()
	SetWorkers(int)
	WorkerDone()
	AllDones() bool

	regret_update_equation(
		t int,
		regret float32,
		cf_reach_prob float32,
		reg_acc float32,
	) float32

	strategy_update_equation(
		t int,
		reach_prob float32,
		action_prob float32,
		strategy_acc float32,
	) float32

	// eval
	Get_avg_strategy(hash string, chi_len int) []float32

	// io
	Save(filename string)
	Load(filename string)

	// new io
	Save_model(filename string, report_interval int, id string, extras []string)
}
