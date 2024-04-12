// SQL statements used by the calibration functions.

package main

const (
	SQLCreateTableCalibration = `
CREATE TABLE calibration (
	elr TEXT NOT NULL,
	total_yards_from INTEGER NOT NULL,
	total_yards_to INTEGER NOT NULL,
	linear_offset_from_m REAL NOT NULL,
	linear_offset_to_m REAL NOT NULL,
	linear_offset_from_norm REAL NOT NULL,
	linear_offset_to_norm REAL NOT NULL,
	accuracy REAL NOT NULL,
	quarter_mile_norm_y REAL NOT NULL
)
`

	SQLCreateIndexCalibration = "CREATE UNIQUE INDEX ix_calibration ON calibration (elr, total_yards_from, total_yards_to)"

	SQLInsertCalibration = "INSERT INTO calibration(elr, total_yards_from, total_yards_to, linear_offset_from_m, linear_offset_to_m, " +
		"linear_offset_from_norm, linear_offset_to_norm, accuracy, quarter_mile_norm_y) values(?,?,?,?,?,?,?,?,?)"

	SQLCreateTableStatistics = `
CREATE TABLE statistics (
  elr TEXT NOT NULL,
  accuracy_count INTEGER NOT NULL,
  accuracy_min REAL NOT NULL,
  accuracy_max REAL NOT NULL,
  accuracy_mean REAL NOT NULL,
  accuracy_median REAL NOT NULL,
  accuracy_std REAL NOT NULL,
  seg_len_count INTEGER NOT NULL,
  seg_len_min INTEGER NOT NULL,
  seg_len_max INTEGER NOT NULL,
  seg_len_mean REAL NOT NULL,
  seg_len_median REAL NOT NULL,
  seg_len_std REAL NOT NULL,
  quarter_mile_norm_count INTEGER NOT NULL,
  quarter_mile_norm_min REAL NOT NULL,
  quarter_mile_norm_max REAL NOT NULL,
  quarter_mile_norm_mean REAL NOT NULL,
  quarter_mile_norm_median REAL NOT NULL,
  quarter_mile_norm_std REAL NOT NULL
)
`

	SQLCreateIndexStatistics = "CREATE UNIQUE INDEX ix_statistics_elr ON statistics (elr)"

	SQLInsertStatistics = "INSERT INTO statistics(elr, accuracy_count, accuracy_min, accuracy_max, accuracy_mean, accuracy_median, accuracy_std, " +
		"seg_len_count, seg_len_min, seg_len_max, seg_len_mean, seg_len_median, seg_len_std, " +
		"quarter_mile_norm_count, quarter_mile_norm_min, quarter_mile_norm_max, quarter_mile_norm_mean, quarter_mile_norm_median, quarter_mile_norm_std) " +
		"values(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"

	QryAllELRs = "SELECT elr, total_yards_from, total_yards_to, shape_length_m, geometry FROM cl"

	QryAllMPsInELR = "SELECT total_yards_from, geometry FROM mp WHERE elr=? ORDER BY total_yards_from"

	SQLVacuumAnalyze = "ANALYZE; VACUUM"
)
