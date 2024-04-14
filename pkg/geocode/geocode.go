// Primary API for geocoding railway ELR and mileage to geographic position.

package geocode

import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"sort"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
	"github.com/paulmach/orb/planar"
)

const (
	maxELRs             = 1_600 // Used to set initial slice sizes.
	calibrationNotFound = "No calibration segment found for ELR %s at TY %d\n"
)

// RailwayPoint represents a geographic position and associated linear accuracy.
type RailwayPoint struct {
	Point    orb.Point // Easting / Northing to EPSG:27700 (metres).
	Accuracy float64   // Calibrated linear accuracy along railway (metres).
}

// GeocoderConfig represents the production database and cache filenames.
type GeocoderConfig struct {
	ProductionDbFn string // Filename of the production database containing ELR and calibration.
	CacheFn        string // Filename of the serialised cache of ELR and calibration.
	VerboseOutput  bool   // Show logging output in event of no calibration segment being found.
}

// ELR represents a single ELR with its associated linear calibration segments.
type ELR struct {
	TyFrom              int                  // Total yards from (irrespective of reporting unit system).
	TyTo                int                  // Total yards to (irrespective of reporting unit system).
	ShapeLen            float64              // Geometry linestring length (metres).
	Metric              bool                 // Linear referencing reporting unit system.
	Geometry            orb.LineString       // Geometry of the centre-line 2D linestring.
	CalibrationSegments []CalibrationSegment // Calibration segments.
}

// Geocoder represents the primary interface offering railway mileage geocoding.
type Geocoder struct {
	ELRs    map[string]ELR  // ELRs with reported extents, geometry, and calibration.
	Metrics map[string]bool // Metric ELRs (reported extents in kilometres).
	config  GeocoderConfig  // Configuration settings.
}

// check aborts if an error is passed in.
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// NewGeocoder is a constructor function to return a Geocoder.
func NewGeocoder(cfg GeocoderConfig) (*Geocoder, error) {
	gc := &Geocoder{}
	gc.config = cfg

	err := gc.loadELRs()
	if err != nil {
		return nil, fmt.Errorf("failed to load ELRs and calibration: %w", err)
	}

	gc.Metrics = gc.MetricELRs()
	return gc, nil
}

// MetricELRs returns a map of ELRs that are reported in kilometres.
func (gc *Geocoder) MetricELRs() map[string]bool {
	metrics := make(map[string]bool)
	for elr, props := range gc.ELRs {
		if props.Metric {
			metrics[elr] = true
		}
	}

	return metrics
}

// IsMetric returns true if the ELR is reported in kilometres.
func (gc *Geocoder) IsMetric(elr string) bool {
	return gc.Metrics[elr]
}

// AllELRs returns all ELR codes in alphabetical order.
func (gc *Geocoder) AllELRs() []string {
	elrs := make([]string, 0, len(gc.ELRs))
	for elr := range gc.ELRs {
		elrs = append(elrs, elr)
	}

	sort.Strings(elrs)
	return elrs
}

// Point returns the point for a given distance (as total yards) on the ELR linestring,
// with linear offset accuracy reported by referring to the milepost calibration points.
func (gc *Geocoder) Point(elr string, ty int) (RailwayPoint, error) {
	elrSegment, err := gc.Find(elr, ty)
	if err != nil {
		if gc.config.VerboseOutput {
			log.Printf(calibrationNotFound, elr, ty)
		}
		return RailwayPoint{}, err
	}

	distance := interpolateSegment(ty, elrSegment.CalibrationSegments[0])
	return RailwayPoint{
			Point:    pointAtDistanceAlongLine(distance, elrSegment.Geometry),
			Accuracy: elrSegment.CalibrationSegments[0].Accuracy},
		nil
}

// Substring returns a portion of the ELR linestring based on the start and end distances (as total yards),
// interpolating linearly as necessary between linestring points.
func (gc *Geocoder) Substring(elr string, tyFrom, tyTo int) (orb.LineString, error) {
	// NOTE: Linear accuracy for either end of the substring is not currently returned.
	elrSegmentFrom, err := gc.Find(elr, tyFrom)
	if err != nil {
		if gc.config.VerboseOutput {
			log.Printf(calibrationNotFound, elr, tyFrom)
		}
		return orb.LineString{}, err
	}
	distanceFrom := interpolateSegment(tyFrom, elrSegmentFrom.CalibrationSegments[0])

	elrSegmentTo, err := gc.Find(elr, tyTo)
	if err != nil {
		if gc.config.VerboseOutput {
			log.Printf(calibrationNotFound, elr, tyTo)
		}

		return orb.LineString{}, err
	}
	distanceTo := interpolateSegment(tyTo, elrSegmentTo.CalibrationSegments[0])

	pts := make([]orb.Point, 0, 32) // Notional initial capacity.
	startPt := pointAtDistanceAlongLine(distanceFrom, elrSegmentFrom.Geometry)
	pts = append(pts, startPt)

	currentDistance := 0.0
	for i := 0; i < len(elrSegmentFrom.Geometry)-1; i++ {
		if currentDistance > distanceFrom && currentDistance < distanceTo {
			pts = append(pts, elrSegmentFrom.Geometry[i])
		} else if currentDistance >= distanceTo {
			break
		}
		currentDistance += planar.Distance(elrSegmentFrom.Geometry[i], elrSegmentFrom.Geometry[i+1])
	}

	endPt := pointAtDistanceAlongLine(distanceTo, elrSegmentFrom.Geometry) // Noting that Geometry To/From are the same ELR.
	pts = append(pts, endPt)
	return pts, err
}

// findCalibrationSegment searches for the target yardage within the calibration slice.
func findCalibrationSegment(calibrationSegments []CalibrationSegment, tyTarget int) (CalibrationSegment, bool) {
	var result CalibrationSegment

	// Binary search; calibration slice is sorted by total yards from.
	lo, hi := 0, len(calibrationSegments)-1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if calibrationSegments[mid].TyFrom <= tyTarget && tyTarget <= calibrationSegments[mid].TyTo {
			result = calibrationSegments[mid]
			return result, true
		} else if tyTarget < calibrationSegments[mid].TyFrom {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}

	// Target yardage not found in calibration slice.
	return result, false
}

// Find searches for the calibration segment that contains the target yardage on the ELR.
// Future enhancement may add option to "clamp" to start or end of ELR limits.
func (gc *Geocoder) Find(elr string, ty int) (ELR, error) {
	e, ok := (gc.ELRs)[elr]
	if !ok {
		return ELR{}, fmt.Errorf("no ELR found: %s", elr)
	}

	// Search the calibration slice and return the row where total_yards_from and total_yards_to contain the target yardage.
	calib, ok := findCalibrationSegment(e.CalibrationSegments, ty)
	if !ok {
		return ELR{}, fmt.Errorf("no calibration found for ELR %s at total yards %d", elr, ty)
	}

	e.CalibrationSegments = []CalibrationSegment{calib}
	return e, nil
}

// loadELRs returns the principal properties, geometry, and calibration of ELRs.
func (gc *Geocoder) loadELRs() error {
	if _, err := os.Stat(gc.config.CacheFn); os.IsNotExist(err) {
		// Cache file doesn't exist, so build and serialise.
		log.Printf("Building cache from production database")
		if !gc.buildCache() || !gc.serialiseCache() {
			return fmt.Errorf("failed to import data / serialise cache")
		}
		return nil
	}

	if !gc.deserialiseCache() {
		return fmt.Errorf("failed to deserialise cache")
	}

	return nil
}

// buildCache reads the production database and builds the ELR cache, returning true if successful.
func (gc *Geocoder) buildCache() bool {
	prodDb, err := sql.Open("sqlite3", fmt.Sprintf("%s?mode=ro", gc.config.ProductionDbFn))
	Check(err)
	defer prodDb.Close()

	const elrSQL = "SELECT elr, total_yards_from, total_yards_to, shape_length_m, l_system, geometry FROM elr"
	elrRows, err := prodDb.Query(elrSQL)
	Check(err)
	defer elrRows.Close()

	calibration := make(map[string][]CalibrationSegment, maxELRs)

	const calibSQL = "SELECT elr, total_yards_from, total_yards_to, linear_offset_from_m, linear_offset_to_m, accuracy " +
		"FROM calibration ORDER BY elr, total_yards_from"
	calibRows, err := prodDb.Query(calibSQL)
	Check(err)
	defer calibRows.Close()

	for calibRows.Next() {
		var elr string
		var c CalibrationSegment
		err := calibRows.Scan(&elr, &c.TyFrom, &c.TyTo, &c.LoFrom, &c.LoTo, &c.Accuracy)
		Check(err)
		calibration[elr] = append(calibration[elr], c)
	}

	gc.ELRs = make(map[string]ELR, maxELRs)

	for elrRows.Next() {
		var e ELR
		var elr string
		var lSystem string
		err := elrRows.Scan(&elr, &e.TyFrom, &e.TyTo, &e.ShapeLen, &lSystem, wkb.Scanner(&e.Geometry))
		Check(err)
		e.Metric = lSystem == "K"
		e.CalibrationSegments = calibration[elr]
		gc.ELRs[elr] = e
	}

	return true
}

// serialiseCache writes the ELR cache to disk, returning true if successful.
func (gc *Geocoder) serialiseCache() bool {
	file, err := os.Create(gc.config.CacheFn)
	Check(err)
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(gc.ELRs)
	Check(err)

	return true
}

// deserialise reads the ELR cache from disk, returning true if successful.
func (gc *Geocoder) deserialiseCache() bool {
	file, err := os.Open(gc.config.CacheFn)
	Check(err)
	defer file.Close()

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(&gc.ELRs)
	Check(err)

	return true
}
