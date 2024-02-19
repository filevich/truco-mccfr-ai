package cfr

import (
	"github.com/filevich/truco-ai/abs"
	"github.com/filevich/truco-ai/info"
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

	CountInfosets() int
	getNumPlayers() int
	GetBuilder() *info.Builder
	GetAbs() abs.IAbstraction
	GetRnode(hash string, chi_len int) *RNode
	samplePartida() *pdt.Partida
	MaxAvgGameValue() float32

	// multi
	Lock()
	Unlock()
	SetWorkers(int)
	WorkerDone()
	AllDones() bool

	regretUpdateEquation(
		t int,
		regret float32,
		cf_reach_prob float32,
		reg_acc float32,
	) float32

	strategyUpdateEquation(
		t int,
		reach_prob float32,
		action_prob float32,
		strategy_acc float32,
	) float32

	// eval
	GetAvgStrategy(hash string, chi_len int) []float32

	// new io
	SaveModel(filename string, report_interval int, id string, extras []string)
}
