package cfr

import "github.com/filevich/truco-ai/utils"

func counterfactualReachProbability(probs []float32, player int) float32 {
	return utils.Prod(probs[:player]) * utils.Prod(probs[player+1:])
}
