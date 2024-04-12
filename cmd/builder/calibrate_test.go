package main

import (
	"geofurlong/pkg/geocode"
	"math"
	"testing"
)

func TestCalibrationPointsToSegments(t *testing.T) {
	const Epsilon = 1e-6

	calibrationPoints := []geocode.CalibrationPoint{
		{Ty: 100, LoMetres: 1_000, LoNormalised: 444},
		{Ty: 200, LoMetres: 1_105, LoNormalised: 555},
		{Ty: 8_888, LoMetres: 9_999, LoNormalised: 666},
	}

	calibrationSegments := calibrationPointsToSegments(calibrationPoints)

	if len(calibrationPoints)-1 != len(calibrationSegments) {
		t.Errorf("Expected length %v, but got %v", len(calibrationPoints)-1, len(calibrationSegments))
	}

	expected := geocode.CalibrationSegmentNormalised{
		TyFrom:           100,          // Total yards from.
		TyTo:             200,          // Total yards to.
		LoMetresFrom:     1_000,        // Linear offset distance from (metres).
		LoMetresTo:       1_105,        // Linear offset distance to (metres).
		LoNormalisedFrom: 444,          // Normalised offset distance from.
		LoNormalisedTo:   555,          // Normalised offset distance to.
		Accuracy:         13.56,        // Linear accuracy (metres).
		QmNormalised:     505.24934372, // Normalised quarter mile length.
	}

	ms0 := calibrationSegments[0]
	if ms0.TyFrom != expected.TyFrom {
		t.Errorf("Expected %v, but got %v", expected.TyFrom, ms0.TyFrom)
	}

	if ms0.TyTo != expected.TyTo {
		t.Errorf("Expected %v, but got %v", expected.TyTo, ms0.TyTo)
	}

	if ms0.LoMetresFrom != expected.LoMetresFrom {
		t.Errorf("Expected %v, but got %v", expected.LoMetresFrom, ms0.LoMetresFrom)
	}

	if ms0.LoMetresTo != expected.LoMetresTo {
		t.Errorf("Expected %v, but got %v", expected.LoMetresTo, ms0.LoMetresTo)
	}

	if ms0.LoNormalisedFrom != expected.LoNormalisedFrom {
		t.Errorf("Expected %v, but got %v", expected.LoNormalisedFrom, ms0.LoNormalisedFrom)
	}

	if ms0.LoNormalisedTo != expected.LoNormalisedTo {
		t.Errorf("Expected %v, but got %v", expected.LoNormalisedTo, ms0.LoNormalisedTo)
	}

	if math.Abs(ms0.Accuracy-expected.Accuracy) > Epsilon {
		t.Errorf("Expected %v, but got %v", expected.Accuracy, ms0.Accuracy)
	}

	if math.Abs(ms0.QmNormalised-expected.QmNormalised) > Epsilon {
		t.Errorf("Expected %v, but got %v", expected.QmNormalised, ms0.QmNormalised)
	}

}
