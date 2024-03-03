package vector

import "math"

type Vector struct {
	Id     int
	Size   string
	Vector []int
}

func Distance(v1 []int, v2 []int) float64 {
	if len(v1) != len(v2) {
		return math.Inf(1) // or any other error handling
	}
	sum := 0.0
	for i := 0; i < len(v1); i++ {
		diff := float64(v1[i] - v2[i])
		sum += diff * diff
	}
	return math.Sqrt(sum)
}
