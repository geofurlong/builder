package geocode

import (
	"testing"

	"github.com/paulmach/orb"
)

func TestPointToOSGR(t *testing.T) {
	testPlaces := getTestPlaces()

	for _, tc := range testPlaces {
		res := PointToOSGR(orb.Point{float64(tc.easting), float64(tc.northing)})

		if res != tc.osgr {
			t.Errorf("PointToOSGR(%d, %d) = %s; want %s", tc.easting, tc.northing, res, tc.osgr)
		}
	}

}

func TestOSGRToPoint(t *testing.T) {
	testPlaces := getTestPlaces()

	for _, tc := range testPlaces {
		pt, err := OSGRToPoint(tc.osgr)
		if err != nil {
			t.Errorf("OSGRToPoint(%q) returned error: %v", tc.osgr, err)
		} else if int(pt.X()) != tc.easting || int(pt.Y()) != tc.northing {
			t.Errorf("OSGRToPoint(%q) = %d, %d; want %d, %d", tc.osgr, int(pt.X()), int(pt.Y()), tc.easting, tc.northing)
		}
	}
}
