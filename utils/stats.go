package utils

import "math/rand"

func Sample(dist []float32) int {
	var (
		r                     float32 = rand.Float32()
		ix                    int     = 0
		cumulativeProbability float32 = 0.0
	)

	for ix < len(dist)-1 {
		cumulativeProbability += dist[ix]
		if r < cumulativeProbability {
			break
		}
		ix++
	}

	return ix
}
