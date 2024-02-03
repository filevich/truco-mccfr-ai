package eval2

import (
	"fmt"
	"time"

	"github.com/filevich/truco-cfr/utils"
)

type Results struct {
	Title              string        `json:"title"`
	TotalNumberOfGames int           `json:"totalNumberOfGames"`
	WonByACounter      int           `json:"wonByACounter"`
	Wons               []float64     `json:"wons"`
	PointsWonDiff      []float32     `json:"pointsWonDiff"`
	Dumbo1             int           `json:"dumbo1"`
	Dumbo2             int           `json:"dumbo2"`
	Delta              time.Duration `json:"delta"`
}

// Winning Percentage The number of the games won by A divided by the total
// number of games
func (r *Results) WP() float32 {
	return float32(r.WonByACounter) / float32(r.TotalNumberOfGames)
}

// Average Difference in Points: The average difference of points scored per
// game between A and B
func (r *Results) ADP() float32 {
	var sum float32 = 0.0
	for _, diff := range r.PointsWonDiff {
		sum += diff
	}
	return sum / float32(len(r.PointsWonDiff))
}

func (r *Results) WaldInterval(agent1 bool) (upper, lower float64) {
	trials, success := 0, 0

	if agent1 {
		trials, success = r.TotalNumberOfGames, r.WonByACounter
	} else {
		// debo devolver el inverso
		trials, success = r.TotalNumberOfGames, r.TotalNumberOfGames-r.WonByACounter
	}

	return utils.WaldAdjusted(success, trials)
}

func (r *Results) String() string {
	u, l := r.WaldInterval(true)
	s := fmt.Sprintf("%.3f [%.3f, %.3f] (di=%d)",
		r.WP(),
		l,
		u,
		r.Dumbo1)
	return s
}
