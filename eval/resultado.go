package eval

import (
	"fmt"

	"github.com/filevich/truco-cfr/utils"
)

type Resultados struct {
	Titulo          string
	Total           int
	Wins_by_A_count int
	Wins            []float64
	Points_won_diff []float32
	Dumbo1          int
	Dumbo2          int
}

// Winning Percentage The number of the games won by A divided by the total
// number of games
func (r *Resultados) WP() float32 {
	return float32(r.Wins_by_A_count) / float32(r.Total)
}

// Average Difference in Points: The average difference of points scored per
// game between A and B
func (r *Resultados) ADP() float32 {
	var sum float32 = 0.0
	for _, diff := range r.Points_won_diff {
		sum += diff
	}
	return sum / float32(len(r.Points_won_diff))
}

func (r *Resultados) WaldInterval(agent1 bool) (upper, lower float64) {
	trials, success := 0, 0

	if agent1 {
		trials, success = r.Total, r.Wins_by_A_count
	} else {
		// debo devolver el inverso
		trials, success = r.Total, r.Total-r.Wins_by_A_count
	}

	return utils.WaldAdjusted(success, trials)
}

func (r *Resultados) String() string {
	// WP ~ Winning Percentage
	// ADP ~ Average Difference in Points
	return fmt.Sprintf("%s\nWP: %.1f%%\nADP: %.3f",
		r.Titulo, r.WP()*float32(100), r.ADP())
}
