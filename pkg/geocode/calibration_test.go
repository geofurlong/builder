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
		{0, CalibrationSegment{0, 1_000, 0, 914.4, 0}, 0},
		{625, CalibrationSegment{500, 1_000, 0, 914.4, 0}, 914.4 / 4.0},
		{750, CalibrationSegment{500, 1_000, 0, 914.4, 0}, 914.4 / 2.0},
		{15_000, CalibrationSegment{10_000, 20_000, 0, 1, -999}, 0.5},
	}

	for _, c := range cases {
		i := interpolateSegment(c.ty, c.c)
		if i != c.expected {
			t.Log("error, should be:", c.expected, "but got:", i)
			t.Fail()
		}
	}
}
