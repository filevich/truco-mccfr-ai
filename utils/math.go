package utils

func Mod(a, b int) int {
	c := a % b
	if c < 0 {
		c += b
	}
	return c
}

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func Max(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}
