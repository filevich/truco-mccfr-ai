package utils

import "math/rand"

func Argmax(dist []float32) int {
	if len(dist) == 0 {
		return -1 // or panic with a meaningful message
	}

	maxIdx := 0
	maxVal := dist[0]

	for i := 1; i < len(dist); i++ {
		if dist[i] > maxVal {
			maxVal = dist[i]
			maxIdx = i
		}
	}

	return maxIdx
}

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

type DictDist struct {
	data  map[int]int
	total float32
}

func NewDictDist(total float32, data map[int]int) *DictDist {
	return &DictDist{
		data:  data,
		total: total,
	}
}

// cummulative density function
func (d *DictDist) CDF(key int) float32 {
	s := 0
	for i := 0; i <= key; i++ {
		if v, ok := d.data[i]; ok {
			s += v
		}
	}
	return float32(s) / d.total
}
