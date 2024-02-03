package eval

import (
	"fmt"
	"time"

	"github.com/filevich/truco-cfr/bot"
)

type TinyEvalResult struct {
	// winrates
	WinRateRandom,
	WinRateSimple float32
	// d-index
	DumboIndexRandom,
	DumboIndexSimple int
	// intervalos de wald
	WaldUpperRandom, WaldLowerRandom float64
	WaldUpperSimple, WaldLowerSimple float64
	Delta                            float64
}

func (r TinyEvalResult) String() string {
	s := fmt.Sprintf("random=%.3f [%.3f,%.3f] (di=%d) - simple=%.3f [%.3f,%.3f] (di=%d)",
		r.WinRateRandom,
		r.WaldLowerRandom,
		r.WaldUpperRandom,
		r.DumboIndexRandom,
		r.WinRateSimple,
		r.WaldLowerSimple,
		r.WaldUpperSimple,
		r.DumboIndexSimple)
	s += fmt.Sprintf(" (%.0fs)", r.Delta)
	return s
}

func TinyEvalFloat(

	agent bot.Agent,
	num_players int,
	ds Dataset,

) *TinyEvalResult {

	ops := []bot.Agent{
		&bot.Random{},
		&bot.Simple{},
	}

	tic := time.Now()
	// nota: SingleSimPartidasBin llama a op.Initialize() !
	res := SingleSimPartidasBin(agent, ops, num_players, ds)

	// winrates
	wr_ale := res[0].WP()
	wr_det := res[1].WP()

	// d-index
	di_ale := res[0].Dumbo1
	di_det := res[1].Dumbo1

	// walds
	wu_ale, wd_ale := res[0].WaldInterval(true)
	wu_det, wd_det := res[1].WaldInterval(true)

	return &TinyEvalResult{
		WinRateRandom: wr_ale,
		WinRateSimple: wr_det,
		// dumbos
		DumboIndexRandom: di_ale,
		DumboIndexSimple: di_det,
		// wald
		WaldUpperRandom: wu_ale, WaldLowerRandom: wd_ale,
		WaldUpperSimple: wu_det, WaldLowerSimple: wd_det,
		// elapsed time
		Delta: time.Since(tic).Seconds(),
	}
}
