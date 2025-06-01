package cfr

import "github.com/filevich/truco-mccfr-ai/utils"

type RNode struct {
	CumulativeRegrets []float32 `json:"cumulativeRegrets"`
	StrategySum       []float32 `json:"strategySum"`
	RegUpdates        int       `json:"regUpdates"`
	StrUpdates        int       `json:"strUpdates"`
}

func NewRNode(n int) *RNode {
	return &RNode{
		CumulativeRegrets: make([]float32, n),
		StrategySum:       make([]float32, n),
	}
}

func (i *RNode) Reset() {
	i.StrategySum = make([]float32, len(i.StrategySum))
	// i.Reg_Updates = 0
	i.StrUpdates = 0
}

func (i *RNode) normalize(xs []float32) []float32 {
	normalized := make([]float32, len(xs))
	var sum float32 = 0.0
	for _, x := range xs {
		sum += x
	}
	if sum == 0 {
		for ix := range xs {
			normalized[ix] = 1 / float32(len(xs))
		}
	} else {
		for ix, x := range xs {
			normalized[ix] = x / sum
		}
	}

	return normalized
}

func (i *RNode) GetStrategy() []float32 {
	// Return regret-matching strategy
	r_plus := make([]float32, len(i.CumulativeRegrets))
	for ix, r := range i.CumulativeRegrets {
		r_plus[ix] = utils.Max(0, r)
	}
	strategy := i.normalize(r_plus)
	return strategy
}

func (i *RNode) GetAverageStrategy() []float32 {
	return i.normalize(i.StrategySum)
}
