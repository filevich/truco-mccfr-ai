package eval

import (
	"fmt"
	"time"

	"github.com/filevich/truco-cfr/bot"
)

func TinyEvalFloat(

	agent bot.Agent,
	num_players int,
	ds Dataset,

) (

	// winrates
	wr_ale,
	wr_det float32,
	// d-index
	di_ale,
	di_det int,
	// intervalos de wald
	wu_ale, wd_ale float64,
	wu_det, wd_det float64,
	delta float64,

) {

	ops := []bot.Agent{
		&bot.Random{},
		&bot.Simple{},
	}

	tic := time.Now()
	// nota: SingleSimPartidasBin llama a op.Initialize() !
	res := SingleSimPartidasBin(agent, ops, num_players, ds)

	// winrates
	wr_ale = res[0].WP()
	wr_det = res[1].WP()

	// d-index
	di_ale = res[0].Dumbo1
	di_det = res[1].Dumbo1

	// walds
	wu_ale, wd_ale = res[0].WaldInterval(true)
	wu_det, wd_det = res[1].WaldInterval(true)

	return wr_ale,
		wr_det,
		// dumbos
		di_ale,
		di_det,
		// wald
		wu_ale, wd_ale,
		wu_det, wd_det,
		// elapsed time
		time.Since(tic).Seconds()
}

func FormatTinyEval(

	ale float32,
	det float32,
	di_ale,
	di_det int,
	wu_ale, wd_ale float64,
	wu_det, wd_det float64,
	delta float64,

) string {

	s := fmt.Sprintf("ale=%.3f (%.3f..%.3f) [di=%d] - det=%.3f (%.3f..%.3f) [di=%d]",
		ale, wd_ale, wu_ale, di_ale,
		det, wd_det, wu_det, di_det)
	s += fmt.Sprintf(" (%.0fs)", delta)
	return s

}
