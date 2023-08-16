package main

func StdDev(highs []float64, lows []float64, volume []int64, closes []float64) (int, int, int, int) {
	var sum float64
	var sum2 float64
	for _, val := range highs {
		sum += val
		sum2 += val * val
	}
	mean := sum / float64(len(highs))
	std := (sum2/float64(len(highs)) - (mean * mean))
	// sum x 2 / num elements - (mean^2)
	return 0, 0, 0, 0
}
