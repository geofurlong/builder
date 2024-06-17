package geocode

import (
	"math"
	"reflect"
	"testing"
)

const epsilon = 1e-4

func TestMinMax(t *testing.T) {
	samples := []float64{1.5, 2.3, 0.8, 4.2, 3.1}
	expectedMin := 0.8
	expectedMax := 4.2

	min, max := minMax(samples)

	if min != expectedMin {
		t.Errorf("Expected min to be %f, but got %f", expectedMin, min)
	}

	if max != expectedMax {
		t.Errorf("Expected max to be %f, but got %f", expectedMax, max)
	}
}
func TestMean(t *testing.T) {
	samples := []float64{1.0, 5.0, 3.0, 4.0, 2.0}
	expected := 3.0
	result := mean(samples, len(samples))
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestMedian(t *testing.T) {
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
			got := median(tt.samples, tt.count)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calcMedian() = %v, want %v", got, tt.want)
			}
		})
	}
}
func TestStdDev(t *testing.T) {
	numbers := []float64{1.0, 4.0, 2.0, 3.0, 5.0}
	count := len(numbers)
	mean := (1.0 + 4.0 + 2.0 + 3.0 + 5.0) / float64(count)
	expected := math.Sqrt((math.Pow(1.0-mean, 2) + math.Pow(4.0-mean, 2) + math.Pow(2.0-mean, 2) + math.Pow(3.0-mean, 2) + math.Pow(5.0-mean, 2)) / float64(count-1))

	result := stdDev(numbers, count, mean)
	if result != expected {
		t.Errorf("Expected %f, but got %f", expected, result)
	}
}

func TestStats(t *testing.T) {
	samples := []float64{2.0, 6.0, 8.0, 12.0, 7.0}
	expected := Statistics{
		Count:  5,
		Min:    2.0,
		Max:    12.0,
		Mean:   7.0,
		Median: 7.0,
		StdDev: 3.605551,
	}

	result := Stats(samples)

	if result.Count != expected.Count ||
		result.Min != expected.Min ||
		result.Max != expected.Max ||
		result.Mean != expected.Mean ||
		result.Median != expected.Median ||
		math.Abs(result.StdDev-expected.StdDev) > epsilon {
		t.Errorf("Expected %v, but got %v", expected, result)
	}
}

func TestCollateStats(t *testing.T) {
	const notUsedInTests = -999

	tests := []struct {
		name                 string
		samples              []CalibrationSegmentNormalised
		expectedAccuracy     Statistics
		expectedSegmentLen   Statistics
		expectedQmNormalised Statistics
	}{
		{
			name: "Test CollateStats",
			samples: []CalibrationSegmentNormalised{
				{
					TyFrom:           0,
					TyTo:             100,
					LoMetresFrom:     notUsedInTests,
					LoMetresTo:       notUsedInTests,
					LoNormalisedFrom: notUsedInTests,
					LoNormalisedTo:   notUsedInTests,
					Accuracy:         5,
					QmNormalised:     420,
				},
				{
					TyFrom:           100,
					TyTo:             200,
					LoMetresFrom:     notUsedInTests,
					LoMetresTo:       notUsedInTests,
					LoNormalisedFrom: notUsedInTests,
					LoNormalisedTo:   notUsedInTests,
					Accuracy:         -20,
					QmNormalised:     470,
				},
				{
					TyFrom:           200,
					TyTo:             300,
					LoMetresFrom:     notUsedInTests,
					LoMetresTo:       notUsedInTests,
					LoNormalisedFrom: notUsedInTests,
					LoNormalisedTo:   notUsedInTests,
					Accuracy:         30,
					QmNormalised:     430,
				},
			},
			expectedAccuracy:     Statistics{Count: 3, Min: -20, Max: 30, Mean: 5, Median: 5, StdDev: 25.0},
			expectedSegmentLen:   Statistics{Count: 3, Min: 100, Max: 100, Mean: 100, Median: 100, StdDev: 0},
			expectedQmNormalised: Statistics{Count: 3, Min: 420, Max: 470, Mean: 440.0, Median: 430, StdDev: 26.45751},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accuracy, segmentLen, qmNormalised := CollateStats(tt.samples)

			if accuracy.Count != tt.expectedAccuracy.Count ||
				accuracy.Min != tt.expectedAccuracy.Min ||
				accuracy.Max != tt.expectedAccuracy.Max ||
				accuracy.Mean != tt.expectedAccuracy.Mean ||
				accuracy.Median != tt.expectedAccuracy.Median ||
				math.Abs(accuracy.StdDev-tt.expectedAccuracy.StdDev) > epsilon {
				t.Errorf("Expected accuracy %v, but got %v", tt.expectedAccuracy, accuracy)
			}

			if segmentLen.Count != tt.expectedSegmentLen.Count ||
				segmentLen.Min != tt.expectedSegmentLen.Min ||
				segmentLen.Max != tt.expectedSegmentLen.Max ||
				segmentLen.Mean != tt.expectedSegmentLen.Mean ||
				segmentLen.Median != tt.expectedSegmentLen.Median ||
				math.Abs(segmentLen.StdDev-tt.expectedSegmentLen.StdDev) > epsilon {
				t.Errorf("Expected segmentLen %v, but got %v", tt.expectedSegmentLen, segmentLen)
			}

			if qmNormalised.Count != tt.expectedQmNormalised.Count ||
				qmNormalised.Min != tt.expectedQmNormalised.Min ||
				qmNormalised.Max != tt.expectedQmNormalised.Max ||
				qmNormalised.Mean != tt.expectedQmNormalised.Mean ||
				qmNormalised.Median != tt.expectedQmNormalised.Median ||
				math.Abs(qmNormalised.StdDev-tt.expectedQmNormalised.StdDev) > epsilon {
				t.Errorf("Expected qmNormalised %v, but got %v", tt.expectedQmNormalised, qmNormalised)
			}
		})
	}

}
