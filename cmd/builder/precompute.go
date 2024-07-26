// Precomputes geocoded railway positions at a specified yardage resolution for all ELRs.

package main

import (
	"bytes"
	"fmt"
	"geofurlong/pkg/geocode"
	"log"
	"os"
)

// precompute generates a CSV file of geocoded railway positions at defined yardage resolution, and
// always including the start and end points of each ELR.
func precompute(cfg GeofurlongConfig, resolution int) { //
	log.Printf("Precomputing geocoded railway positions at %d yard resolution\n", resolution)

	gcCfg := geocode.GeocoderConfig{
		ProductionDbFn: cfg["production_db"],
		CacheFn:        cfg["cache_fn"],
		VerboseOutput:  false,
	}

	gc, err := geocode.NewGeocoder(gcCfg)
	geocode.Check(err)

	// Set up projection conversion from OSGB projected (EPSG:27700) to geographic longitude / latitude (EPSG:4326).
	pj := geocode.OSGBToLonLat()

	file, err := os.Create(fmt.Sprintf("%s/geofurlong_precomputed_%.4dy.csv", cfg["precompute_dir"], resolution))
	geocode.Check(err)
	defer file.Close()

	fmt.Fprintln(file, "elr,total_yards,mileage,easting,northing,longitude,latitude,osgr,accuracy")

	// buffer 1,000 records before printing to output file to improve performance.
	const BatchBufferLen = 1_000
	var buffer bytes.Buffer
	count := 0

	for _, elr := range gc.AllELRs() {
		prop := gc.ELRs[elr]
		for ty := prop.TyFrom; ty <= prop.TyTo; ty++ {
			if ty%resolution != 0 && ty != prop.TyFrom && ty != prop.TyTo {
				// Position is not at a resolution point or the start or end point of the ELR, so skip.
				continue
			}

			if elr == "FTC" && ty >= 21450 {
				// NOTE: Hardcoded to skip CTRL (HS1) beyond mainland limits.
				continue
			}

			pt, err := gc.Point(elr, ty)
			geocode.Check(err)

			osgr := geocode.PointToOSGR(pt.Point)
			lonLat := geocode.Reproject(pt.Point, pj)

			// 6 decimal places for latitude / longitude is approximately 0.11 metre precision,
			// notionally equivalent to the 0.1 metre precision of the OSGB Easting / Northing.
			// Linear accuracy is rounded to nearest metre.
			buffer.WriteString(fmt.Sprintf("%s,%d,%s,%.1f,%.1f,%.6f,%.6f,%s,%d\n",
				elr,
				ty,                                     // Total yards.
				geocode.FmtTotalYards(ty, prop.Metric), // Formatted mileage.
				pt.Point[0],                            // OS Easting.
				pt.Point[1],                            // OS Northing.
				lonLat.X(),                             // Longitude (decimal degrees).
				lonLat.Y(),                             // Latitude (decimal degrees).
				osgr,                                   // Ordnance Survey Grid Reference.
				int(pt.Accuracy+0.5)))                  // Railway linear accuracy (metres).

			count++
			if count >= BatchBufferLen {
				fmt.Fprint(file, buffer.String())
				buffer.Reset()
				count = 0
			}

		}
	}

	// Print residual buffer data to output file.
	if buffer.Len() > 0 {
		fmt.Fprint(file, buffer.String())
	}
}
