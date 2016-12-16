package models

import (
    "math"
    "strconv"
)

// DOC: Round is a helper function for cleaning up calculated results. The
//      results is rounded to the fourth decimal.
func Round(input float64) float64 {
	rounded := math.Floor((input*10000.0)+0.5) / 10000.0
	return rounded
}

// DOC: ToString is a helper function to translate float64 into string.
func ToString(number float64) string {
	return strconv.FormatFloat(number, 'f', 4, 64)
}
