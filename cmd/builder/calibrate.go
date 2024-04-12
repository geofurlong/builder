// Calibration functions to support the linear distance calculation along an ELR linestring for a given mileage.

package main

import (
	"database/sql"
	"fmt"
	"geofurlong/pkg/geocode"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkb"
)

// ELRFeature represents a single ELR feature.
type ELRFeature struct {
	elr      string         // ELR code.
	tyFrom   int            // Total yards from.
	tyTo     int            // Total yards to.
	length   float64        // Geometry linestring length (metres).
	geometry orb.LineString // Geometry of the centre-line linestring.
}

// calibrationPointsToSegments pairwise transforms calibration points to calibration segments (and normalises).
func calibrationPointsToSegments(calibPoints []geocode.CalibrationPoint) []geocode.CalibrationSegmentNormalised {
	calibSegments := make([]geocode.CalibrationSegmentNormalised, 0, len(calibPoints)-1)

	for i := 0; i < len(calibPoints)-1; i++ {
		current := calibPoints[i]
		next := calibPoints[i+1]

		lenReported := float64(next.Ty-current.Ty) * geocode.YardsToMetres
		lenMeasured := next.LoMetres - current.LoMetres
		accuracy := lenMeasured - lenReported
		qmNormalised := (geocode.QuarterMileYards / float64(next.Ty-current.Ty)) * (lenMeasured / geocode.YardsToMetres)

		multi := geocode.CalibrationSegmentNormalised{
			TyFrom:           current.Ty,
			TyTo:             next.Ty,
			LoMetresFrom:     current.LoMetres,
			LoMetresTo:       next.LoMetres,
			LoNormalisedFrom: current.LoNormalised,
			LoNormalisedTo:   next.LoNormalised,
			Accuracy:         accuracy,
			QmNormalised:     qmNormalised,
		}

		calibSegments = append(calibSegments, multi)
	}

	return calibSegments
}

// Calibrator represents the database connections and prepared statements for the calibration process.
type Calibrator struct {
	dbELR                 *sql.DB   // ELR database.
	dbMilepost            *sql.DB   // Milepost database.
	dbCalibration         *sql.DB   // Calibration database.
	stmtMilepost          *sql.Stmt // Prepared statement for milepost query.
	rowsELR               *sql.Rows // Rows for the ELR query.
	tx                    *sql.Tx   // Calibration database transaction.
	stmtInsertCalibration *sql.Stmt // Prepared statement for inserting calibration rows.
	stmtInsertStatistics  *sql.Stmt // Prepared statement for inserting calibration statistics rows.
}

// initialise opens the centre-line and milepost databases, creates the calibration database and prepares the SQL statements.
func (c *Calibrator) initialise(ELRFn, MilepostFn, CalibrationFn string) error {
	var err error
	c.dbELR, err = sql.Open("sqlite3", fmt.Sprintf("%s?mode=ro", ELRFn))
	geocode.Check(err)

	c.dbMilepost, err = sql.Open("sqlite3", fmt.Sprintf("%s?mode=ro", MilepostFn))
	geocode.Check(err)

	c.stmtMilepost, err = c.dbMilepost.Prepare(QryAllMPsInELR)
	geocode.Check(err)

	c.rowsELR, err = c.dbELR.Query(QryAllELRs)
	geocode.Check(err)

	_, err = os.Stat(CalibrationFn) // Delete calibration database if it exists.
	if err == nil {
		geocode.Check(os.Remove(CalibrationFn))
	}

	c.dbCalibration, err = sql.Open("sqlite3", CalibrationFn)
	geocode.Check(err)

	c.tx, err = c.dbCalibration.Begin() // Begin a database transaction.
	geocode.Check(err)

	_, err = c.tx.Exec(SQLCreateTableCalibration)
	geocode.Check(err)

	_, err = c.tx.Exec(SQLCreateTableStatistics)
	geocode.Check(err)

	c.stmtInsertCalibration, err = c.tx.Prepare(SQLInsertCalibration)
	geocode.Check(err)

	c.stmtInsertStatistics, err = c.tx.Prepare(SQLInsertStatistics)
	geocode.Check(err)

	return nil
}

// close closes the database connections and prepared statements.
func (c *Calibrator) close() {
	geocode.Check(c.dbELR.Close())
	geocode.Check(c.dbMilepost.Close())
	geocode.Check(c.dbCalibration.Close())
	geocode.Check(c.stmtMilepost.Close())
	geocode.Check(c.rowsELR.Close())
	geocode.Check(c.stmtInsertCalibration.Close())
	geocode.Check(c.stmtInsertStatistics.Close())
}

// appendDB appends the calibration data to the database.
func (c *Calibrator) appendDB(elr string, calibPoints []geocode.CalibrationPoint) error {
	calibSegments := calibrationPointsToSegments(calibPoints)

	// Save rows to calibration table.
	for _, cm := range calibSegments {
		_, err := c.stmtInsertCalibration.Exec(elr, cm.TyFrom, cm.TyTo, cm.LoMetresFrom, cm.LoMetresTo,
			cm.LoNormalisedFrom, cm.LoNormalisedTo, cm.Accuracy, cm.QmNormalised)
		geocode.Check(err)
	}

	accuracy, segLen, qmNormalised := geocode.CollateStats(calibSegments)

	// Save rows to calibration statistics table.
	_, err := c.stmtInsertStatistics.Exec(elr,
		accuracy.Count, accuracy.Min, accuracy.Max, accuracy.Mean, accuracy.Median, accuracy.StdDev,
		segLen.Count, segLen.Min, segLen.Max, segLen.Mean, segLen.Median, segLen.StdDev,
		qmNormalised.Count, qmNormalised.Min, qmNormalised.Max, qmNormalised.Mean, qmNormalised.Median, qmNormalised.StdDev)
	geocode.Check(err)

	return nil
}

// computeAndSaveCalibration computes and saves the calibration for each ELR.
func (c *Calibrator) computeAndSaveCalibration() error {
	for c.rowsELR.Next() {
		// Loop through all ELR centre-line records.
		var (
			ef      ELRFeature
			tyMP    int
			pointMP orb.Point
		)

		err := c.rowsELR.Scan(&ef.elr, &ef.tyFrom, &ef.tyTo, &ef.length, wkb.Scanner(&ef.geometry))
		geocode.Check(err)
		rowsMP, err := c.stmtMilepost.Query(ef.elr)
		geocode.Check(err)
		defer rowsMP.Close()
		gotFirstMP := false

		// Initial size based on 99% of ELRs having 300 or less mileposts in total.
		cs := make([]geocode.CalibrationPoint, 0, 300)

		for rowsMP.Next() {
			// Loop through all milepost records for the current ELR.
			err = rowsMP.Scan(&tyMP, wkb.Scanner(&pointMP))
			geocode.Check(err)

			if !gotFirstMP {
				gotFirstMP = true
				if tyMP > ef.tyFrom {
					// The mileage of the first milepost is greater than the low mileage end of the ELR,
					// so record a quasi-milepost at the low mileage end of the ELR.
					csStart := geocode.CalibrationPoint{Ty: ef.tyFrom, LoMetres: 0.0, LoNormalised: 0.0}
					cs = append(cs, csStart)
				}
			}

			// Record the milepost projected against the ELR geometry.
			nearestPt, _ := geocode.NearestPointOnLine(&ef.geometry, pointMP)
			lo := geocode.DistanceAlongLine(&ef.geometry, nearestPt)
			loNormalised := lo / ef.length
			csNormalised := geocode.CalibrationPoint{Ty: tyMP, LoMetres: lo, LoNormalised: loNormalised}
			cs = append(cs, csNormalised)
		}

		if tyMP < ef.tyTo {
			// The mileage of the last milepost is less than the high mileage end of the ELR,
			// so record a quasi-milepost at the high mileage end of the ELR.
			csEnd := geocode.CalibrationPoint{Ty: ef.tyTo, LoMetres: ef.length, LoNormalised: 1.0}
			cs = append(cs, csEnd)
		}

		c.appendDB(ef.elr, cs)
	}

	return nil
}

// finalise commits the database transaction and performs optimisation.
func (c *Calibrator) finalise() error {
	_, err := c.tx.Exec(SQLCreateIndexCalibration)
	geocode.Check(err)

	_, err = c.tx.Exec(SQLCreateIndexStatistics)
	geocode.Check(err)

	geocode.Check(c.tx.Commit())

	_, err = c.dbCalibration.Exec(SQLVacuumAnalyze)
	geocode.Check(err)

	return nil
}

// calibrate performs the calibration process, referencing mileposts against ELR centre-lines, and saving to a database.
func calibrate(cfg GeofurlongConfig) {
	log.Print("Calibration started")
	c := Calibrator{}
	c.initialise(cfg["cl_db"], cfg["mp_db"], cfg["calib_db"])
	defer c.close()
	geocode.Check(c.computeAndSaveCalibration())
	geocode.Check(c.finalise())
	log.Print("Calibration completed")
}
