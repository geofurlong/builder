// Support types and function for linear offset calibration of ELR geometry relative to Milepost points.

package geocode

// CalibrationPoint represents linear calibration values at a railway point.
type CalibrationPoint struct {
	Ty           int     // Total yards.
	LoMetres     float64 // Linear offset (metres).
	LoNormalised float64 // Linear offset (normalised 0 -> 1).
}

// CalibrationSegment represents linear calibration values between two railway points.
type CalibrationSegment struct {
	TyFrom   int     // Low mileage end of calibration segment (as total yards).
	TyTo     int     // High mileage end of calibration segment (as total yards).
	LoFrom   float64 // Linear offset (metres) at low mileage end.
	LoTo     float64 // Linear offset (metres) at high mileage end.
	Accuracy float64 // Accuracy of calibration segment, comparing reported versus measured length (metres).
}

// CalibrationSegmentNormalised represents linear calibration values (including normalised values) between two railway points.
type CalibrationSegmentNormalised struct {
	TyFrom           int     // Low mileage end of calibration segment (as total yards).
	TyTo             int     // High mileage end of calibration segment (as total yards).
	LoMetresFrom     float64 // Linear offset (metres) at low mileage end.
	LoMetresTo       float64 // Linear offset (metres) at high mileage end.
	LoNormalisedFrom float64 // Linear offset (normalised 0 -> 1) at low mileage end.
	LoNormalisedTo   float64 // Linear offset (normalised 0 -> 1) at high mileage end.
	Accuracy         float64 // Accuracy (metres).
	QmNormalised     float64 // "Normalised" quarter mile length (relative to 440 yards).
}

// interpolateSegment returns the linear interpolated offset value within a given linear calibration segment.
func interpolateSegment(tyTarget int, c CalibrationSegment) float64 {
	return c.LoFrom + (float64(tyTarget)-float64(c.TyFrom))/(float64(c.TyTo)-float64(c.TyFrom))*(c.LoTo-c.LoFrom)
}
