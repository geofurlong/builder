package geocode

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/encoding/wkt"
)

func TestFindCalibrationSegment(t *testing.T) {
	cases := []struct {
		data     []CalibrationSegment
		tyTarget int
		expected CalibrationSegment
		found    bool
	}{
		{
			data: []CalibrationSegment{
				{
					TyFrom: 0,
					TyTo:   10},
				{
					TyFrom: 11,
					TyTo:   20},
				{
					TyFrom: 21,
					TyTo:   30},
			},
			tyTarget: 15,
			expected: CalibrationSegment{
				TyFrom: 11,
				TyTo:   20},
			found: true,
		},
		{
			data: []CalibrationSegment{
				{
					TyFrom: 0,
					TyTo:   10},
				{
					TyFrom: 11,
					TyTo:   20},
				{
					TyFrom: 21,
					TyTo:   30},
			},
			tyTarget: 35,
			expected: CalibrationSegment{},
			found:    false,
		},
	}

	for i, tc := range cases {
		result, found := findCalibrationSegment(tc.data, tc.tyTarget)
		if found != tc.found || result != tc.expected {
			t.Errorf("Test case %d failed: expected (%v, %v), got (%v, %v)", i+1, tc.expected, tc.found, result, found)
		}
	}
}

func TestPoint(t *testing.T) {
	cases := []struct {
		geometry      orb.LineString
		calibration   CalibrationSegment
		distance_m    int
		expectedPoint RailwayPoint
	}{
		{
			geometry: orb.LineString{{0, 0}, {20, 0}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     50,
				LoFrom:   0,
				LoTo:     20,
				Accuracy: 0},
			distance_m:    0,
			expectedPoint: RailwayPoint{orb.Point{0, 0}, 0},
		},
		{
			geometry: orb.LineString{{0, 0}, {20, 0}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     50,
				LoFrom:   0,
				LoTo:     20,
				Accuracy: 0},
			distance_m:    25,
			expectedPoint: RailwayPoint{orb.Point{10, 0}, 0},
		},
		{
			geometry: orb.LineString{{0, 0}, {20, 0}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     50,
				LoFrom:   0,
				LoTo:     20,
				Accuracy: 0},
			distance_m:    50,
			expectedPoint: RailwayPoint{orb.Point{20, 0}, 0},
		},
		{
			geometry: orb.LineString{{0, 0}, {0, 20}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     50,
				LoFrom:   0,
				LoTo:     20,
				Accuracy: 0},
			distance_m:    25,
			expectedPoint: RailwayPoint{orb.Point{0, 10}, 0},
		},
	}

	elr := "test"
	for _, c := range cases {
		testELR := ELR{
			TyFrom:              0,
			TyTo:                50,
			ShapeLen:            20,
			Metric:              false,
			Geometry:            c.geometry,
			CalibrationSegments: []CalibrationSegment{c.calibration},
		}

		gc := Geocoder{}
		gc.ELRs = make(map[string]ELR)
		gc.ELRs[elr] = testELR

		rp, err := gc.Point(elr, c.distance_m)
		if err != nil {
			t.Log("error, should not return an error but got:", err)
			t.Fail()
		}
		if rp != c.expectedPoint {
			t.Log("error, should be:", c.expectedPoint, "but got:", rp)
			t.Fail()
		}
	}
}

func TestSubstring(t *testing.T) {
	cases := []struct {
		geometry    orb.LineString
		calibration CalibrationSegment
		tyFrom      int
		tyTo        int
		expectedWKT string
	}{
		// Simple horizontal line segments with two points.
		{
			geometry: orb.LineString{{0, 0}, {100, 0}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     50,
				LoFrom:   0,
				LoTo:     100,
				Accuracy: 0},
			tyFrom:      0,
			tyTo:        0,
			expectedWKT: "LINESTRING(0 0,0 0)"},
		{
			geometry: orb.LineString{{0, 0}, {100, 0}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     50,
				LoFrom:   0,
				LoTo:     100,
				Accuracy: 0},
			tyFrom:      50,
			tyTo:        50,
			expectedWKT: "LINESTRING(100 0,100 0)"},
		{
			geometry: orb.LineString{{0, 0}, {100, 0}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     10,
				LoFrom:   0,
				LoTo:     100,
				Accuracy: 0},
			tyFrom:      2,
			tyTo:        5,
			expectedWKT: "LINESTRING(20 0,50 0)"},

		// Simple vertical line segments with two points.
		{
			geometry: orb.LineString{{0, 0}, {0, 100}},
			calibration: CalibrationSegment{
				TyFrom: 0, TyTo: 10,
				LoFrom:   0,
				LoTo:     100,
				Accuracy: 0},
			tyFrom:      0,
			tyTo:        0,
			expectedWKT: "LINESTRING(0 0,0 0)"},
		{
			geometry: orb.LineString{{0, 0}, {0, 100}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     10,
				LoFrom:   0,
				LoTo:     100,
				Accuracy: 0},
			tyFrom:      10,
			tyTo:        10,
			expectedWKT: "LINESTRING(0 100,0 100)"},
		{
			geometry: orb.LineString{{0, 0}, {0, 100}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     10,
				LoFrom:   0,
				LoTo:     100,
				Accuracy: 0},
			tyFrom:      2,
			tyTo:        5,
			expectedWKT: "LINESTRING(0 20,0 50)"},

		// Simple diagonal line segments.
		{
			geometry: orb.LineString{{0, 0}, {3, 4}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     5,
				LoFrom:   0,
				LoTo:     5,
				Accuracy: 0},
			tyFrom: 0, tyTo: 0,
			expectedWKT: "LINESTRING(0 0,0 0)"},
		{
			geometry: orb.LineString{{0, 0}, {3, 4}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     5,
				LoFrom:   0,
				LoTo:     5,
				Accuracy: 0},
			tyFrom:      5,
			tyTo:        5,
			expectedWKT: "LINESTRING(3 4,3 4)"},
		{
			geometry: orb.LineString{{0, 0}, {3, 4}},
			calibration: CalibrationSegment{
				TyFrom:   0,
				TyTo:     5,
				LoFrom:   0,
				LoTo:     5,
				Accuracy: 0},
			tyFrom:      0,
			tyTo:        5,
			expectedWKT: "LINESTRING(0 0,3 4)"},
		{
			geometry: orb.LineString{{2, 3}, {5, 7}},
			calibration: CalibrationSegment{
				TyFrom:   666,
				TyTo:     777,
				LoFrom:   0,
				LoTo:     5,
				Accuracy: 0},
			tyFrom:      777,
			tyTo:        777,
			expectedWKT: "LINESTRING(5 7,5 7)"},

		// Two diagonal line segments.
		{
			geometry: orb.LineString{{0, 0}, {3, 4}, {18, 4}},
			calibration: CalibrationSegment{
				TyFrom:   2_000,
				TyTo:     2_020,
				LoFrom:   0,
				LoTo:     20,
				Accuracy: 0},
			tyFrom:      2_000,
			tyTo:        2_000,
			expectedWKT: "LINESTRING(0 0,0 0)"},
		{
			geometry: orb.LineString{{0, 0}, {3, 4}, {18, 4}},
			calibration: CalibrationSegment{
				TyFrom:   2_000,
				TyTo:     2_020,
				LoFrom:   0,
				LoTo:     20,
				Accuracy: 0},
			tyFrom:      2_020,
			tyTo:        2_020,
			expectedWKT: "LINESTRING(18 4,18 4)"},
		{
			geometry: orb.LineString{{0, 0}, {3, 4}, {18, 4}},
			calibration: CalibrationSegment{
				TyFrom:   2_000,
				TyTo:     2_020,
				LoFrom:   0,
				LoTo:     20,
				Accuracy: 0},
			tyFrom:      2_000,
			tyTo:        2_020,
			expectedWKT: "LINESTRING(0 0,3 4,18 4)"},
		{
			geometry: orb.LineString{{0, 0}, {3, 4}, {18, 4}},
			calibration: CalibrationSegment{
				TyFrom:   2_000,
				TyTo:     2_020,
				LoFrom:   0,
				LoTo:     20,
				Accuracy: 0},
			tyFrom:      2_010,
			tyTo:        2_020,
			expectedWKT: "LINESTRING(8 4,18 4)"},
	}

	elrCode := "not used"

	for _, c := range cases {
		testELR := ELR{
			Geometry:            c.geometry,
			CalibrationSegments: []CalibrationSegment{c.calibration},
		}

		gc := Geocoder{}
		gc.ELRs = make(map[string]ELR)
		gc.ELRs[elrCode] = testELR

		ls, _ := gc.Substring(elrCode, c.tyFrom, c.tyTo)
		lsWKT := wkt.MarshalString(ls)
		if lsWKT != c.expectedWKT {
			t.Log("error, should be:", c.expectedWKT, "but got:", lsWKT)
			t.Fail()
		}
	}
}
