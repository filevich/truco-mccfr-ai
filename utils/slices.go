package utils

// removeElements returns a copy of xs without any element that may be present in ys
func CopyWithoutThese(xs []int, ys ...int) []int {
	// create a map to store the elements of ys
	m := make(map[int]bool)
	for _, y := range ys {
		m[y] = true
	}

	// create a slice to store the result
	result := make([]int, 0, len(xs))

	// iterate over xs and append only the elements that are not in m
	for _, x := range xs {
		if !m[x] {
			result = append(result, x)
		}
	}

	// return the result
	return result
}
