package api

import (
	"math"
	"testing"
)

var hexTests = []struct {
	in  string
	out float64
}{
	{"3", float64(3)},
	{"7", float64(7)},
	{"F", float64(15)},
}

func TestHexToFloat(t *testing.T) {
	for _, tt := range hexTests {
		t.Run(tt.in, func(t *testing.T) {
			got := hexToFloat(tt.in)
			if got != tt.out {
				t.Errorf("got %f, want %f", got, tt.out)
			}
		})
	}
}

var boolTests = []struct {
	in  interface{}
	out bool
}{
	{"0", false},
	{"1", true},
	{"bbq", false},
	{3, false},
}

func TestGetBool(t *testing.T) {
	for _, tt := range boolTests {
		got := getBool(tt.in)
		if got != tt.out {
			t.Errorf("got %t, want %t", got, tt.out)
		}
	}
}

var floatTests = []struct {
	in  interface{}
	out float64
}{
	{float32(1), float64(1)},
	{"bbq", math.NaN()},
	{true, math.NaN()},
	{float64(3), float64(3)},
	{12, float64(12)},
}

func TestGetFloat(t *testing.T) {
	for _, tt := range floatTests {
		got := getFloat(tt.in)

		// NaN != NaN, so NaN expecting NaN fails without the and check.
		if got != tt.out && !(math.IsNaN(got) && math.IsNaN(tt.out)) {
			t.Errorf("got %f, want %f", got, tt.out)
		}
	}
}
