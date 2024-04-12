// GeoFurlong builder.
// https://www.github.com/geofurlong
// https://www.geofurlong.com

package main

import (
	"geofurlong/pkg/geocode"
	"log"
	"strconv"
	"sync"
	"time"
)

// main is the entry point for the GeoFurlong builder.
func main() {
	startTime := time.Now()

	config, err := readConfig()
	geocode.Check(err)
	log.Printf("GeoFurlong builder (version %s) started", config["version"])

	// Convert the source geospatial files from Shapefile to SQLite format.
	runPython(config, "convert.py", "")

	// Compute the linear calibration, referencing mileposts against the ELR centre-line.
	calibrate(config)

	// Build the production database.
	buildProductionDb(config)

	resolutions := []int{
		8800, // 5 miles ~ 8.045 km
		1760, // 1 mile ~ 1.609 km
		440,  // 1/4 mile ~ 402 m
		220,  // 1/8 mile ~ 201 m
		110,  // 1/16 mile ~ 101 m
		22,   // 1/80 mile (1 chain) ~ 20 m
	}

	var wg sync.WaitGroup

	// Compute geographic positions of railway points at multiple yardage resolutions concurrently.
	for _, resolution := range resolutions {
		wg.Add(1)
		go func(resolution int) {
			defer wg.Done()
			precompute(config, resolution)
		}(resolution)
	}
	wg.Wait()

	// Build gazetteer of railway points at multiple yardage resolutions concurrently.
	for _, resolution := range resolutions {
		wg.Add(1)
		go func(resolution int) {
			defer wg.Done()
			runPython(config, "gazetteer.py", strconv.Itoa(resolution))
		}(resolution)
	}
	wg.Wait()

	// Compact the highest resolution gazetteer to offer concise look-up.
	aggregateGazetteer(config)

	elapsedTime := time.Since(startTime)
	log.Printf("build complete - duration: %s", elapsedTime)
}
