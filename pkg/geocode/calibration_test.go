package geocode

import (
	"testing"
)

func TestInterpolateSegment(t *testing.T) {
	cases := []struct {
		ty       int
		c        CalibrationSegment
		expected float64
	}{
		{
			ty: 0,
			c: CalibrationSegment{
				TyFrom:   0,
				TyTo:     1_000,
				LoFrom:   0,
				LoTo:     914.4,
				Accuracy: 0},
			expected: 0},
		{
			ty: 625,
			c: CalibrationSegment{
				TyFrom:   500,
				TyTo:     1_000,
				LoFrom:   0,
				LoTo:     914.4,
				Accuracy: 0},
			expected: 914.4 / 4.0},
		{
			ty: 750,
			c: CalibrationSegment{
				TyFrom:   500,
				TyTo:     1_000,
				LoFrom:   0,
				LoTo:     914.4,
				Accuracy: 0},
			expected: 914.4 / 2.0},
		{
			ty: 15_000,
			c: CalibrationSegment{
				TyFrom:   10_000,
				TyTo:     20_000,
				LoFrom:   0,
				LoTo:     1,
				Accuracy: -999},
			expected: 0.5},
	}

	for _, c := range cases {
		i := interpolateSegment(c.ty, c.c)
		if i != c.expected {
			t.Log("error, should be:", c.expected, "but got:", i)
			t.Fail()
		}
	}
}
