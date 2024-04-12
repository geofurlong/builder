// Gazetteer aggregator groups gazetteer information at specific yardage intervals into a compact "to" and "from" format.

package main

import (
	"database/sql"
	"fmt"
	"geofurlong/pkg/geocode"
	"log"
	"math"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Constants for the Gazetteer database grouping.
// May rationalise / remove in future, depending on whether individual tables are referenced in client programs.
const (
	NRRegionID         = 1 // Network Rail region grouping code.
	CountryAdminAreaID = 2 // Country and admin area grouping code.
	DistrictPlaceID    = 3 // District and place grouping code.
)

// GazetteerRow represents an unaggregated row in the Gazetteer database.
type GazetteerRow struct {
	ty        int    // Total yards.
	region    string // Network Rail Region railway point is contained within.
	country   string // Country railway point is contained within.
	adminArea string // Administrative Area railway point is contained within.
	district  string // County / District of Nearest Populated Place.
	place     string // Nearest Populated Place name.
	distance  int    // Distance from railway point to nearest Populated Place (metres).
}

// AggregateGroup represents an aggregated group of Gazetteer rows.
type AggregateGroup struct {
	tyFrom       int    // Total yards from.
	tyTo         int    // Total yards to.
	value        string // Group value.
	minDistance  int    // Minimum distance (metres).
	maxDistance  int    // Maximum distance (metres).
	meanDistance int    // Mean distance (metres).
}

// AggregatorConfig represents the Gazetteer Aggregator configuration.
type AggregatorConfig struct {
	unaggregatedDb string                 // Filename of the unaggregated gazetteer database.
	aggregatedCSV  string                 // Filename of the aggregated gazetteer CSV file.
	gcConfig       geocode.GeocoderConfig // Geocoder configuration.
}

// Aggregator represents the Gazetteer Aggregator.
type Aggregator struct {
	dbGaz   *sql.DB         // Unaggregated gazetteer database.
	stmtGaz *sql.Stmt       // Prepared statement for unaggregated gazetteer database.
	elrs    []string        // All ELR codes.
	metrics map[string]bool // Metric ELRs.
	buf     strings.Builder // Output buffer.
}

// NewAggregator creates a new Gazetteer Aggregator.
func NewAggregator(config AggregatorConfig) *Aggregator {
	dbGaz, err := sql.Open("sqlite3", config.unaggregatedDb)
	geocode.Check(err)
	stmtGaz, err := dbGaz.Prepare("SELECT total_yards, nr_region, country, admin_area, county_district, place_name, distance_m FROM gazetteer_summary WHERE elr=? ORDER BY total_yards")
	geocode.Check(err)

	gc, err := geocode.NewGeocoder(config.gcConfig)
	geocode.Check(err)

	elrs := gc.AllELRs()
	metrics := gc.MetricELRs()

	buf := strings.Builder{}
	return &Aggregator{dbGaz, stmtGaz, elrs, metrics, buf}
}

// Close closes the database and prepared statement.
func (a *Aggregator) close() {
	a.stmtGaz.Close()
	a.dbGaz.Close()
}

// tyToStr converts total yardages to corresponding formatted strings.
func (a *Aggregator) tyToStr(elr string, tyFrom int, tyTo int) (string, string) {
	metric := a.metrics[elr]
	return geocode.FmtTotalYards(tyFrom, metric), geocode.FmtTotalYards(tyTo, metric)
}

// aggregate aggregates (groups) the gazetteer for a given ELR.
func (a *Aggregator) aggregate(elr string) {
	rows, err := a.stmtGaz.Query(elr)
	geocode.Check(err)
	defer rows.Close()

	// Build a slice of gazetteer rows for the given ELR.
	// At 22y resolution, min=6, max=19,710, mean=557, median=115 unaggregated gazetteer rows per ELR.
	gazetteerRows := make([]GazetteerRow, 0, 600)
	for rows.Next() {
		row := GazetteerRow{}
		err := rows.Scan(&row.ty, &row.region, &row.country, &row.adminArea, &row.district, &row.place, &row.distance)
		gazetteerRows = append(gazetteerRows, row)
		geocode.Check(err)
	}

	groupsNRRegion := aggregateDataText(gazetteerRows, func(r GazetteerRow) string { return r.region })
	for _, gNRRegion := range groupsNRRegion {
		mileageFrom, mileageTo := a.tyToStr(elr, gNRRegion.tyFrom, gNRRegion.tyTo)
		a.buf.WriteString(fmt.Sprintf("%s,%d,%d,%d,%s,%s,%s,,,,\n",
			elr, NRRegionID, gNRRegion.tyFrom, gNRRegion.tyTo, mileageFrom, mileageTo, gNRRegion.value))
	}

	const Delimiter = "|"

	groupsCountryAdminArea := aggregateDataText(gazetteerRows, func(r GazetteerRow) string { return r.country + Delimiter + r.adminArea })
	for _, gCountryAdminArea := range groupsCountryAdminArea {
		countryAdminArea := strings.Split(gCountryAdminArea.value, Delimiter)
		country := countryAdminArea[0]
		adminArea := countryAdminArea[1]
		mileageFrom, mileageTo := a.tyToStr(elr, gCountryAdminArea.tyFrom, gCountryAdminArea.tyTo)
		a.buf.WriteString(fmt.Sprintf("%s,%d,%d,%d,%s,%s,%s,\"%s\",,,\n",
			elr, CountryAdminAreaID, gCountryAdminArea.tyFrom, gCountryAdminArea.tyTo, mileageFrom, mileageTo, country, adminArea))
	}

	groupsPlace := aggregateDataNumeric(gazetteerRows, func(r GazetteerRow) string { return r.district + Delimiter + r.place })
	for _, gPlace := range groupsPlace {
		districtPlace := strings.Split(gPlace.value, Delimiter)
		district := districtPlace[0]
		place := districtPlace[1]
		mileageFrom, mileageTo := a.tyToStr(elr, gPlace.tyFrom, gPlace.tyTo)
		a.buf.WriteString(fmt.Sprintf("%s,%d,%d,%d,%s,%s,\"%s\",\"%s\",%d,%d,%d\n",
			elr, DistrictPlaceID, gPlace.tyFrom, gPlace.tyTo, mileageFrom, mileageTo,
			district, place, gPlace.minDistance, gPlace.maxDistance, gPlace.meanDistance))
	}

}

// outputCSV outputs the aggregated gazetteer as a CSV file.
func (a *Aggregator) outputCSV(fn string) {
	file, err := os.Create(fn)
	geocode.Check(err)
	defer file.Close()

	_, err = file.WriteString(a.buf.String())
	geocode.Check(err)

	err = file.Sync()
	geocode.Check(err)
}

// csvToDb builds a SQLite gazetteer database for 22y resolution, including helper tables.
func csvToDb(csvFn string, dbFn string, sqlFn string) {
	deleteFile(dbFn)
	sql, err := os.ReadFile(sqlFn)
	geocode.Check(err)
	modifiedSql := strings.Replace(string(sql), "gazetteer_aggregated.csv", csvFn, -1)
	runSQLiteCommand(dbFn, modifiedSql)
}

// aggregateDataText aggregates the data for a given text value.
func aggregateDataText(data []GazetteerRow, valueFunc func(GazetteerRow) string) []AggregateGroup {
	var groups []AggregateGroup
	var currentGroup *AggregateGroup

	for _, r := range data {
		value := valueFunc(r)
		if currentGroup == nil || currentGroup.value != value {
			if currentGroup != nil {
				meanOffset := (currentGroup.tyTo + r.ty) / 2
				// Subtract 1 from meanOffset to avoid overlapping groups.
				currentGroup.tyTo = meanOffset - 1
				groups = append(groups, *currentGroup)
				currentGroup = &AggregateGroup{tyFrom: meanOffset, tyTo: r.ty, value: value}
			} else {
				currentGroup = &AggregateGroup{tyFrom: r.ty, tyTo: r.ty, value: value}
			}
		} else {
			if r.ty > currentGroup.tyTo {
				currentGroup.tyTo = r.ty
			}
		}
	}

	if currentGroup != nil {
		groups = append(groups, *currentGroup)
	}

	return groups
}

// aggregateDataNumeric aggregates the data for a given numeric value.
func aggregateDataNumeric(data []GazetteerRow, valueFunc func(GazetteerRow) string) []AggregateGroup {
	var groups []AggregateGroup
	var currentGroup *AggregateGroup
	var sumDistance int
	var count int

	for _, r := range data {
		value := valueFunc(r)
		if currentGroup == nil || currentGroup.value != value {
			if currentGroup != nil {
				meanOffset := (currentGroup.tyTo + r.ty) / 2
				// Subtract 1 from meanOffset to avoid overlapping groups.
				currentGroup.tyTo = meanOffset - 1
				currentGroup.meanDistance = int(math.Round(float64(sumDistance) / float64(count)))
				groups = append(groups, *currentGroup)
				currentGroup = &AggregateGroup{tyFrom: meanOffset, tyTo: r.ty, value: value,
					minDistance: r.distance, maxDistance: r.distance, meanDistance: r.distance}
				sumDistance = r.distance
				count = 1
			} else {
				currentGroup = &AggregateGroup{tyFrom: r.ty, tyTo: r.ty, value: value,
					minDistance: r.distance, maxDistance: r.distance, meanDistance: r.distance}
				sumDistance = r.distance
				count = 1
			}
		} else {
			if r.ty > currentGroup.tyTo {
				currentGroup.tyTo = r.ty
			}
			if r.distance < currentGroup.minDistance {
				currentGroup.minDistance = r.distance
			}
			if r.distance > currentGroup.maxDistance {
				currentGroup.maxDistance = r.distance
			}
			sumDistance += r.distance
			count++
		}
	}

	if currentGroup != nil {
		currentGroup.meanDistance = int(math.Round(float64(sumDistance) / float64(count)))
		groups = append(groups, *currentGroup)
	}

	return groups
}

func aggregateGazetteer(cfg GeofurlongConfig) {
	// Future consideration for full database normalisation (esp. NR Region, Country).
	log.Println("Gazetteer aggregator started")
	unaggregatedDb := cfg["gazetteer_dir"] + "/geofurlong_gazetteer_0022y.sqlite"
	aggregatedCSV := cfg["gazetteer_dir"] + "/geofurlong_gazetteer_aggregated.csv"

	config := AggregatorConfig{
		gcConfig: geocode.GeocoderConfig{
			ProductionDbFn: cfg["production_db"],
			CacheFn:        cfg["cache_fn"],
			VerboseOutput:  false},
		unaggregatedDb: unaggregatedDb,
		aggregatedCSV:  aggregatedCSV}

	aggregator := NewAggregator(config)
	defer aggregator.close()

	counter := 0
	aggregator.buf.WriteString("elr,group_id,offset_from,offset_to,mileage_from,mileage_to,value_1,value_2,min_distance,max_distance,mean_distance\n")
	for _, elr := range aggregator.elrs {
		if counter%50 == 0 {
			fmt.Printf("\r%d", counter)
			os.Stdout.Sync()
		}
		counter++
		aggregator.aggregate(elr)
	}

	fmt.Printf("\r%d\n", counter)
	log.Printf("saving aggregated gazetteer as CSV")
	aggregator.outputCSV(aggregatedCSV)

	log.Printf("saving aggregated gazetteer as database")
	aggregatedDb := cfg["gazetteer_aggregated_db"]
	sqlFn := cfg["scripts_dir"] + "/gazetteer_aggregate.sql"
	csvToDb(aggregatedCSV, aggregatedDb, sqlFn)

	log.Println("Gazetteer aggregator completed")
}
