package geocode

import (
	"math"
	"testing"

	"github.com/paulmach/orb"
)

func TestReproject(t *testing.T) {
	const Epsilon = 1e-5

	testPlaces := getTestPlaces()

	planarPoints := []orb.Point{}
	for _, testPlace := range testPlaces {
		planarPoints = append(planarPoints, orb.Point{float64(testPlace.easting), float64(testPlace.northing)})
	}

	geoPoints := ReprojectMulti(planarPoints, OSGBToLonLat())
	for i, testPlace := range testPlaces {
		deltaX := geoPoints[i].Point().X() - testPlace.lonLat.X()
		deltaY := geoPoints[i].Point().Y() - testPlace.lonLat.Y()

		if math.Abs(deltaX) > Epsilon || math.Abs(deltaY) > Epsilon {
			t.Log("ReprojectMulti error, should be:", testPlace.lonLat, "but got:", geoPoints[i])
			t.Fail()
		}
	}

	pjToGeo := OSGBToLonLat()
	for _, testPlace := range testPlaces {
		geoPoint := Reproject(orb.Point{float64(testPlace.easting), float64(testPlace.northing)}, pjToGeo)
		deltaX := geoPoint.X() - testPlace.lonLat.X()
		deltaY := geoPoint.Y() - testPlace.lonLat.Y()

		if math.Abs(deltaX) > Epsilon || math.Abs(deltaY) > Epsilon {
			t.Log("Reproject error, should be:", testPlace.lonLat, "but got:", geoPoint, "for location:", testPlace.lonLat, " - ", testPlace.name)
			t.Fail()
		}
	}

}
