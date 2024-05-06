package geocode

import (
	"math"
	"reflect"
	"testing"
)

func TestCalcMinMax(t *testing.T) {
	samples := []float64{1.5, 2.3, 0.8, 4.2, 3.1}
	expectedMin := 0.8
	expectedMax := 4.2

	min, max := calcMinMax(samples)

	if min != expectedMin {
		t.Errorf("Expected min to be %f, but got %f", expectedMin, min)
	}

	if max != expectedMax {
		t.Errorf("Expected max to be %f, but got %f", expectedMax, max)
	}
}
func TestCalcMean(t *testing.T) {
	samples := []float64{1.0, 5.0, 3.0, 4.0, 2.0}
	expected := 3.0
	result := calcMean(samples, len(samples))
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestCalcMedian(t *testing.T) {
	tests := []struct {
		name    string
		samples []float64
		count   int
		want    float64
	}{
		{
			name:    "Odd number of samples",
			samples: []float64{1, 5, 4, 3, 2},
			count:   5,
			want:    3,
		},
		{
			name:    "Even number of samples",
			samples: []float64{4, 3, 2, 1},
			count:   4,
			want:    2.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calcMedian(tt.samples, tt.count)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calcMedian() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestCalcStdDev(t *testing.T) {
	numbers := []float64{1.0, 4.0, 2.0, 3.0, 5.0}
	count := len(numbers)
	mean := (1.0 + 4.0 + 2.0 + 3.0 + 5.0) / float64(count)
	expected := math.Sqrt((math.Pow(1.0-mean, 2) + math.Pow(4.0-mean, 2) + math.Pow(2.0-mean, 2) + math.Pow(3.0-mean, 2) + math.Pow(5.0-mean, 2)) / float64(count-1))

	result := calcStdDev(numbers, count, mean)
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

// TODO implement TestCollateStats
