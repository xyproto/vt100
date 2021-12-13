package main

import "math"

func abs(a int) int {
	if a >= 0 {
		return a
	}
	return -a
}

func distance(x1, x2, y1, y2 int) float64 {
	return math.Sqrt((float64(x1)*float64(x1) - float64(x2)*float64(x2)) + (float64(y1)*float64(y1) - float64(y2)*float64(y2)))
}
