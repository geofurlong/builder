// Computes statistics for a set of numbers.

package geocode

import (
	"math"
	"sort"
)

// Statistics represents the combined statistics of a set of numbers.
type Statistics struct {
	Count  int
	Min    float64
	Max    float64
	Mean   float64
	Median float64
	StdDev float64
}

// calcMinMax calculates the minimum and maximum of a set of numbers.
func calcMinMax(nums []float64) (float64, float64) {
	min := math.MaxFloat64
	max := -min

	for _, num := range nums {
		if num < min {
			min = num
		}
		if num > max {
			max = num
		}
	}

	return min, max
}

// calcMean calculates the mean of a set of numbers.
func calcMean(nums []float64, count int) float64 {
	sum := 0.0

	for _, num := range nums {
		sum += num
	}

	return sum / float64(count)
}

// / calcMedian calculates the median of a set of numbers.
func calcMedian(nums []float64, count int) float64 {
	sort.Float64s(nums)
	middle := count / 2

	if count%2 == 0 {
		return (nums[middle-1] + nums[middle]) / 2
	} else {
		return nums[middle]
	}
}

// calcStdDev calculates the standard deviation of a set of numbers.
func calcStdDev(nums []float64, count int, mean float64) float64 {
	sumOfSquares := 0.0

	for _, num := range nums {
		sumOfSquares += math.Pow(num-mean, 2)
	}

	// count-1 used as estimating population mean from the sample.
	// This is compatible with computation of other data analysis libraries.
	return math.Sqrt(sumOfSquares / float64(count-1))
}

// Stats calculates the combined statistics of a set of numbers.
func Stats(nums []float64) Statistics {
	count := len(nums)
	min, max := calcMinMax(nums)
	mean := calcMean(nums, count)
	median := calcMedian(nums, count)
	stddev := calcStdDev(nums, count, mean)

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
