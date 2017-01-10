package models

import (
	"math"
	"strconv"
)

func Round(input float64) float64 {
	rounded := math.Floor((input*10000.0)+0.5) / 10000.0
	return rounded
}

func ToString(number float64) string {
	return strconv.FormatFloat(number, 'f', 4, 64)
}
