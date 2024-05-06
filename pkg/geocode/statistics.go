// Computes statistics for a sample of numbers.

package geocode

import (
	"math"
	"sort"
)

// Statistics represents the combined statistics of a sample of numbers.
type Statistics struct {
	Count  int
	Min    float64
	Max    float64
	Mean   float64
	Median float64
	StdDev float64
}

// calcMinMax calculates the minimum and maximum of a sample of numbers.
func calcMinMax(samples []float64) (float64, float64) {
	min := math.MaxFloat64
	max := -min

	for _, sample := range samples {
		if sample < min {
			min = sample
		}
		if sample > max {
			max = sample
		}
	}

	return min, max
}

// calcMean calculates the mean of a sample of numbers.
func calcMean(samples []float64, count int) float64 {
	sum := 0.0

	for _, sample := range samples {
		sum += sample
	}

	return sum / float64(count)
}

// / calcMedian calculates the median of a sample of numbers.
func calcMedian(samples []float64, count int) float64 {
	sort.Float64s(samples)
	middle := count / 2

	if count%2 == 0 {
		return (samples[middle-1] + samples[middle]) / 2
	} else {
		return samples[middle]
	}
}

// calcStdDev calculates the standard deviation of a sample of numbers.
func calcStdDev(samples []float64, count int, mean float64) float64 {
	sumOfSquares := 0.0

	for _, sample := range samples {
		sumOfSquares += math.Pow(sample-mean, 2)
	}

	// count-1 used as estimating population mean from the sample.
	// This is compatible with computation of other data analysis libraries.
	return math.Sqrt(sumOfSquares / float64(count-1))
}

// Stats calculates the combined statistics of a sample of numbers.
func Stats(samples []float64) Statistics {
	count := len(samples)
	min, max := calcMinMax(samples)
	mean := calcMean(samples, count)
	median := calcMedian(samples, count)
	stddev := calcStdDev(samples, count, mean)

	return Statistics{
		Count:  count,
		Min:    min,
		Max:    max,
		Mean:   mean,
		Median: median,
		StdDev: stddev,
	}
}

// CollateStats returns the statistics for accuracy, segment length, and normalised Quarter Mile length for an ELR.
func CollateStats(calibSegments []CalibrationSegmentNormalised) (Statistics, Statistics, Statistics) {
	accuracyValues := make([]float64, len(calibSegments))
	for i, cm := range calibSegments {
		accuracyValues[i] = cm.Accuracy
	}
	accuracy := Stats(accuracyValues)

	segmentLenValues := make([]float64, len(calibSegments))
	for i, calibSegment := range calibSegments {
		segmentLenValues[i] = float64(calibSegment.TyTo - calibSegment.TyFrom)
	}
	segmentLen := Stats(segmentLenValues)

	// Quarter Mile statistics.
	qmNormalisedValues := make([]float64, len(calibSegments))
	for i, calibSegment := range calibSegments {
		qmNormalisedValues[i] = calibSegment.QmNormalised
	}
	qmNormalised := Stats(qmNormalisedValues)

	return accuracy, segmentLen, qmNormalised
}
