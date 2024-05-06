// Reprojecting functions for Easting / Northing points to and from Longitude / Latitude points.

package geocode

import (
	"github.com/paulmach/orb"
	"github.com/twpayne/go-proj/v10"
)

const (
	// Ordnance Survey National Grid (OSGB36) Co-ordinate Reference System (CRS).
	ProjectedCRS = "EPSG:27700"

	// World Geodetic System 1984 (WGS84) Co-ordinate Reference System (CRS).
	GeographicCRS = "EPSG:4326"
)

// OSGBtoLongLat returns a pointer to the transformer from projected OSGB36 (EPSG:27700) to geographic Longitude / Latitude (EPSG:4326).
func OSGBtoLongLat() *proj.PJ {
	pj, err := proj.NewCRSToCRS(ProjectedCRS, GeographicCRS, nil)
	Check(err)

	return pj
}

// Reproject takes a projected Easting / Northing point and returns the corresponding Longitude / Latitude point.
func Reproject(point orb.Point, pj *proj.PJ) orb.Point {
	latLon, err := pj.Forward(proj.Coord{point.X(), point.Y()})
	Check(err)

	// Note order of X / Y versus Longitude / Latitude is intentional (due to library utilising GDAL).
	return orb.Point{latLon.Y(), latLon.X()}
}

// ReprojectMulti takes a slice of projected Easting / Northing points and returns the corresponding Longitude / Latitude points slice.
func ReprojectMulti(points []orb.Point) []orb.Point {
	latLons := make([]orb.Point, len(points))
	pj, err := proj.NewCRSToCRS(ProjectedCRS, GeographicCRS, nil)
	Check(err)

	for i, point := range points {
		latLon, err := pj.Forward(proj.Coord{point.X(), point.Y()})
		Check(err)

		// Note order of X / Y versus Longitude / Latitude is intentional (due to library utilising GDAL).
		latLons[i] = orb.Point{latLon.Y(), latLon.X()}
	}

	return latLons
}
