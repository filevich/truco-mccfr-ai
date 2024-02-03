package cfr

import "github.com/filevich/truco-cfr/utils"

type RNode struct {
	Cumulative_regrets []float32
	Strategy_sum       []float32
	Reg_Updates        int
	Str_Updates        int
	// usado solo por LMCCFR - Linear (External-Sampling) Monte Carlo Counterfactual Regret Minimization
}

func NewRNode(n int) *RNode {
	return &RNode{
		Cumulative_regrets: make([]float32, n),
		Strategy_sum:       make([]float32, n),
	}
}

func (i *RNode) Reset() {
	i.Strategy_sum = make([]float32, len(i.Strategy_sum))
	// i.Reg_Updates = 0
	i.Str_Updates = 0
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

func (i *RNode) Get_strategy() []float32 {
	// Return regret-matching strategy
	r_plus := make([]float32, len(i.Cumulative_regrets))
	for ix, r := range i.Cumulative_regrets {
		r_plus[ix] = utils.Max(0, r)
	}
	strategy := i.normalize(r_plus)
	return strategy
}

func (i *RNode) Get_average_strategy() []float32 {
	return i.normalize(i.Strategy_sum)
}
