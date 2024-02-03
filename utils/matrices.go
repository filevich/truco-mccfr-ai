package utils

// manipulacion de matrices

func Prod(xs []float32) float32 {
	var res float32 = 1.0
	for _, x := range xs {
		res *= x
	}
	return res
}

func Sum_float32_slices(xs, ys []float32) []float32 {
	res := make([]float32, len(xs))
	for i := 0; i < len(xs); i++ {
		res[i] = xs[i] + ys[i]
	}
	return res
}

func Ones(n int) []float32 {
	res := make([]float32, n)
	for i := 0; i < n; i++ {
		res[i] = 1
	}
	return res
}

func Ndot(xs []float32, yss [][]float32) []float32 {
	n := len(yss[0])
	res := make([]float32, n)

	for ix := 0; ix < n; ix++ {
		for i := 0; i < len(xs); i++ {
			// row := xs[i]
			// col := yss[i][ix]
			// res[ix] += row * col
			res[ix] += xs[i] * yss[i][ix]
		}
	}

	return res
}
