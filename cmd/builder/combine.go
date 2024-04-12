// Build the production database by combining ELR and Calibration data.

package main

import (
	"fmt"
	"geofurlong/pkg/geocode"
	"log"
)

// buildProductionDb combines ELR (attributes, geometry) and Calibration into production database.
func buildProductionDb(cfg GeofurlongConfig) {
	log.Print("Building production database")
	deleteFile(cfg["production_db"])

	input := fmt.Sprintf(`
BEGIN;

ATTACH DATABASE '%s' AS ext_cl;
ATTACH DATABASE '%s' AS ext_calib;

CREATE TABLE elr_tmp (elr TEXT NOT NULL, route TEXT NOT NULL, section TEXT, remarks TEXT, quail_book TEXT, grouping TEXT, neighbours TEXT, PRIMARY KEY (elr));
.mode csv
.headers on
.import %s elr_tmp

CREATE TABLE elr (elr TEXT, l_system TEXT, shape_length_m FLOAT, total_yards_from INTEGER, total_yards_to INTEGER,
                  route TEXT NOT NULL, section TEXT, remarks TEXT, quail_book TEXT NOT NULL, grouping TEXT, neighbours TEXT,
                  geometry BLOB NOT NULL, PRIMARY KEY (elr));


-- Join manually maintained non-geospatial ELR attributes with geospatial ELR centre-line data.
INSERT INTO elr
    SELECT cl.elr, cl.l_system, cl.shape_length_m, cl.total_yards_from, cl.total_yards_to,
           elr_tmp.route, elr_tmp.section, elr_tmp.remarks, elr_tmp.quail_book, elr_tmp.grouping, elr_tmp.neighbours,
           cl.geometry
    FROM ext_cl.cl AS cl
    LEFT OUTER JOIN elr_tmp ON cl.elr = elr_tmp.elr;

DROP TABLE elr_tmp;

CREATE VIEW elr_metric AS SELECT elr FROM elr WHERE l_system="K" ORDER BY elr;

-- Subset of calibration stored.
-- For external GIS systems (e.g. PostGIS), use the normalised linear offset values for point/substring operations.
CREATE TABLE calibration AS SELECT elr, total_yards_from, total_yards_to, linear_offset_from_m, linear_offset_to_m, CAST(accuracy AS INT) AS accuracy FROM ext_calib.calibration;
CREATE UNIQUE INDEX ix_calibration ON calibration (elr, total_yards_from, total_yards_to);

COMMIT;

DETACH DATABASE ext_cl;
DETACH DATABASE ext_calib;

ANALYZE;
VACUUM;
`, cfg["cl_db"], cfg["calib_db"], cfg["elr_csv"])

	runSQLiteCommand(cfg["production_db"], input)
	log.Print("Production database built")

	// initialise the cache and serialise to disk.
	gcCfg := geocode.GeocoderConfig{
		ProductionDbFn: cfg["production_db"],
		CacheFn:        cfg["cache_fn"],
		VerboseOutput:  false,
	}

	_, err := geocode.NewGeocoder(gcCfg)
	geocode.Check(err)
}
