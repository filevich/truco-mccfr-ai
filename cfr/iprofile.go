package cfr

import "github.com/filevich/truco-cfr/bot"

// silent:
// solo los [*]
// no silent:
// todos

type IProfile interface {
	//
	Init(trainer ITrainer)
	Continue(trainer ITrainer) bool
	// multi
	IsMulti() bool
	GetThreads() int
	// algoritmo
	IsPrunable(trainer ITrainer) bool
	// reporte
	IsSilent() bool
	IsFullySilent() bool
	PercentageDone(t int) float32 // [0 .. 100]
	//
	Checkpoint(ITrainer)
	Check(ITrainer)
	// eplotabilidad
	Exploit() bot.Agent
}
