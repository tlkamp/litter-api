package util

import (
	"math"
	"strconv"
)

func Bool(i interface{}) bool {
	switch unk := i.(type) {
	case bool:
		return unk
	case string:
		b, err := strconv.ParseBool(unk)
		if err != nil {
			return false
		}
		return b
	default:
		return false
	}
}

func Float(i interface{}) float64 {
	switch i := i.(type) {
	case string:
		f, err := strconv.ParseFloat(i, 32)
		if err != nil {
			return math.NaN()
		}
		return f
	case int:
		return float64(i)
	case float32:
		return float64(i)
	case float64:
		return i
	default:
		return math.NaN()
	}
}

func HexToFloat(s string) float64 {
	f, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return math.NaN()
	}
	return float64(f)
}
