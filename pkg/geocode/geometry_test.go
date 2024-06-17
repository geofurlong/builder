package geocode

import (
	"math"
	"testing"

	"github.com/paulmach/orb"
)

func TestInterpolatePoint(t *testing.T) {
	cases := []struct {
		name          string
		point1        orb.Point
		point2        orb.Point
		ratio         float64
		expectedPoint orb.Point
	}{
		{
			name:          "Horizontal line 1",
			point1:        orb.Point{0, 0},
			point2:        orb.Point{1, 0},
			ratio:         0.5,
			expectedPoint: orb.Point{0.5, 0},
		},
		{
			name:          "Vertical line",
			point1:        orb.Point{0, 0},
			point2:        orb.Point{0, 1},
			ratio:         0.5,
			expectedPoint: orb.Point{0, 0.5},
		},
		{
			name:          "Horizontal line 2",
			point1:        orb.Point{100, 999},
			point2:        orb.Point{500, 999},
			ratio:         0.25,
			expectedPoint: orb.Point{200, 999},
		},
		{
			name:          "Vertical line 2",
			point1:        orb.Point{999, 1_000},
			point2:        orb.Point{999, 11_000},
			ratio:         0.75,
			expectedPoint: orb.Point{999, 8_500},
		},
		{
			name:          "Diagonal line 1",
			point1:        orb.Point{0, 0},
			point2:        orb.Point{1000, 50},
			ratio:         0.5,
			expectedPoint: orb.Point{500, 25},
		},
		{
			name:          "Diagonal line 2",
			point1:        orb.Point{1_000, 2_000},
			point2:        orb.Point{5_000, 12_000},
			ratio:         0.1,
			expectedPoint: orb.Point{1_400, 3_000},
		},
	}

	for _, c := range cases {
		gotPoint := interpolatePoint(c.point1, c.point2, c.ratio)
		if gotPoint != c.expectedPoint {
			t.Log("error, should be:", c.expectedPoint, "but got:", gotPoint)
			t.Fail()
		}
	}

}

func TestNearestPointOnLine(t *testing.T) {
	tests := []struct {
		name         string
		line         orb.LineString
		point        orb.Point
		wantPoint    orb.Point
		wantDistance float64
	}{
		{
			name:         "Test 1",
			line:         orb.LineString{{0, 0}, {1, 1}, {2, 2}, {3, 3}},
			point:        orb.Point{2, 0},
			wantPoint:    orb.Point{1, 1},
			wantDistance: math.Sqrt(2),
		},
		{
			name:         "Test 2",
			line:         orb.LineString{{0, 0}, {10, 0}},
			point:        orb.Point{5, 5},
			wantPoint:    orb.Point{5, 0},
			wantDistance: 5.0},
		{
			name:         "Test 3",
			line:         orb.LineString{{0, 0}, {10, 0}},
			point:        orb.Point{-3, -4},
			wantPoint:    orb.Point{0, 0},
			wantDistance: 5.0,
		},
		{
			name:         "Test 4",
			line:         orb.LineString{{0, 0}, {10, 0}},
			point:        orb.Point{3, 4},
			wantPoint:    orb.Point{3, 0},
			wantDistance: 4.0,
		},
		{
			name:         "Test 5",
			line:         orb.LineString{{0, 0}, {10, 0}},
			point:        orb.Point{10, 0},
			wantPoint:    orb.Point{10, 0},
			wantDistance: 0.0,
		},
		{
			name:         "Test 6",
			line:         orb.LineString{{0, 0}, {10, 0}},
			point:        orb.Point{13, -4},
			wantPoint:    orb.Point{10, 0},
			wantDistance: 5.0,
		},
		{
			name:         "Test 7",
			line:         orb.LineString{{0, 0}, {10, 10}},
			point:        orb.Point{5, 5},
			wantPoint:    orb.Point{5, 5},
			wantDistance: 0.0,
		},
		{
			name:         "Test 8",
			line:         orb.LineString{{0, 0}, {10, 10}},
			point:        orb.Point{4, 6},
			wantPoint:    orb.Point{5, 5},
			wantDistance: math.Sqrt(2),
		},
		{
			name:         "Test 9",
			line:         orb.LineString{{0, 0}, {5, 0}, {10, 0}, {15, 0}},
			point:        orb.Point{7, 3},
			wantPoint:    orb.Point{7, 0},
			wantDistance: 3.0,
		},
		{
			name:         "Test 10",
			line:         orb.LineString{{0, 0}, {0, 5}, {0, 10}, {0, 15}},
			point:        orb.Point{-4, 8},
			wantPoint:    orb.Point{0, 8},
			wantDistance: 4.0,
		},
		{
			name:         "Test 11",
			line:         orb.LineString{{0, 0}, {5, 5}, {10, 10}, {15, 15}},
			point:        orb.Point{6, 8},
			wantPoint:    orb.Point{7, 7},
			wantDistance: math.Sqrt(2),
		},
		{
			name:         "Test 12",
			line:         orb.LineString{{0, 0}, {5, -5}, {10, -10}, {15, -15}},
			point:        orb.Point{8, -6},
			wantPoint:    orb.Point{7, -7},
			wantDistance: math.Sqrt(2),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPoint, gotDistance := NearestPointOnLine(&tt.line, tt.point)

			if gotPoint != tt.wantPoint {
				t.Errorf("Unexpected nearest point. Got %v, want %v", gotPoint, tt.wantPoint)
			}

			if math.Abs(gotDistance-tt.wantDistance) > 1e-9 {
				t.Errorf("Unexpected distance. Got %v, want %v", gotDistance, tt.wantDistance)
			}
		})
	}
}

func TestNearestPointOnSegment(t *testing.T) {
	tests := []struct {
		name             string
		startPoint       orb.Point
		endPoint         orb.Point
		targetPoint      orb.Point
		expectedPoint    orb.Point
		expectedDistance float64
	}{
		{
			name:             "Test 1",
			startPoint:       orb.Point{0, 0},
			endPoint:         orb.Point{10, 0},
			targetPoint:      orb.Point{5, 5},
			expectedPoint:    orb.Point{5, 0},
			expectedDistance: 5.0,
		},
		{
			name:             "Test 2",
			startPoint:       orb.Point{0, 0},
			endPoint:         orb.Point{0, 10},
			targetPoint:      orb.Point{5, 5},
			expectedPoint:    orb.Point{0, 5},
			expectedDistance: 5.0,
		},
		{
			name:             "Test 3 - Zero Length Segment",
			startPoint:       orb.Point{55, 33},
			endPoint:         orb.Point{55, 33},
			targetPoint:      orb.Point{55, 3},
			expectedPoint:    orb.Point{55, 33},
			expectedDistance: 30.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPoint, gotDistance := nearestPointOnSegment(tt.startPoint, tt.endPoint, tt.targetPoint)

			if gotPoint != tt.expectedPoint {
				t.Errorf("Expected point %v, got %v", tt.expectedPoint, gotPoint)
			}

			if math.Abs(gotDistance-tt.expectedDistance) > 1e-9 {
				t.Errorf("Expected distance %v, got %v", tt.expectedDistance, gotDistance)
			}
		})
	}
}

func TestPointOnLineSegment(t *testing.T) {
	tests := []struct {
		name        string
		startPoint  orb.Point
		endPoint    orb.Point
		targetPoint orb.Point
		expectedHit bool
	}{
		{
			name:        "Test 1",
			startPoint:  orb.Point{0, 0},
			endPoint:    orb.Point{1, 1},
			targetPoint: orb.Point{0.5, 0.5},
			expectedHit: true,
		},
		{
			name:        "Test 2",
			startPoint:  orb.Point{0, 0},
			endPoint:    orb.Point{1, 1},
			targetPoint: orb.Point{1, 1},
			expectedHit: true,
		},
		{
			name:        "Test 3",
			startPoint:  orb.Point{0, 0},
			endPoint:    orb.Point{1, 1},
			targetPoint: orb.Point{2, 2},
			expectedHit: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pointOnLineSegment(tt.startPoint, tt.endPoint, tt.targetPoint)

			if got != tt.expectedHit {
				t.Errorf("Expected %v, got %v", tt.expectedHit, got)
			}
		})
	}
}

func TestDistanceAlongLine(t *testing.T) {
	line := &orb.LineString{{0, 0}, {0, 1}, {1, 1}}

	cases := []struct {
		name             string
		point            orb.Point
		expectedDistance float64
	}{
		{
			name:             "Test 1",
			point:            orb.Point{0, 0.5},
			expectedDistance: 0.5,
		},
		{
			name:             "Test 2",
			point:            orb.Point{0, 1},
			expectedDistance: 1.0,
		},
		{
			name:             "Test 3",
			point:            orb.Point{0.5, 1},
			expectedDistance: 1.5,
		},
		{
			name:             "Test 4",
			point:            orb.Point{1, 1},
			expectedDistance: 2,
		},
	}

	for _, c := range cases {
		got := DistanceAlongLine(line, c.point)
		if !almostEqual(got, c.expectedDistance) {
			t.Errorf("DistanceAlongLine(%v, %v) == %v, want %v", line, c.point, got, c.expectedDistance)
		}
	}
}

// almostEqual checks if two float64s are approximately equal, considering floating point precision.
func almostEqual(a, b float64) bool {
	const Epsilon = 1e-9
	return math.Abs(a-b) <= Epsilon
}

func TestPointAtDistanceAlongLine(t *testing.T) {
	cases := []struct {
		name          string
		line          orb.LineString
		distance      float64
		expectedPoint orb.Point
	}{
		{
			name:          "Test 1",
			line:          orb.LineString{{0, 0}, {10, 0}},
			distance:      3,
			expectedPoint: orb.Point{3, 0},
		},
		{
			name:          "Test 2",
			line:          orb.LineString{{0, 0}, {0, 100}},
			distance:      99,
			expectedPoint: orb.Point{0, 99},
		},
		{
			name:          "Test 3",
			line:          orb.LineString{{0, 0}, {10, 0}, {10, 110}},
			distance:      70,
			expectedPoint: orb.Point{10, 60},
		},
		{
			name:          "Test 4",
			line:          orb.LineString{{3, 4}, {0, 0}, {-3, -4}},
			distance:      0,
			expectedPoint: orb.Point{3, 4},
		},
		{
			name:          "Test 5",
			line:          orb.LineString{{3, 4}, {0, 0}, {-3, -4}},
			distance:      5,
			expectedPoint: orb.Point{0, 0},
		},
		{
			name:          "Test 6",
			line:          orb.LineString{{3, 4}, {0, 0}, {-3, -4}},
			distance:      10,
			expectedPoint: orb.Point{-3, -4},
		},
		{
			name:          "Test 7",
			line:          orb.LineString{{3, 4}, {0, 0}, {-3, -4}, {987, -4}},
			distance:      1000,
			expectedPoint: orb.Point{987, -4},
		},
		{
			name:          "Test 8 - Negative distance",
			line:          orb.LineString{{3, 4}, {0, 0}, {-3, -4}},
			distance:      -999,
			expectedPoint: orb.Point{3, 4},
		},
	}

	const Epsilon = 1e-6

	for _, c := range cases {
		gotPoint := pointAtDistanceAlongLine(c.distance, c.line)
		deltaX := gotPoint.X() - c.expectedPoint.X()
		deltaY := gotPoint.Y() - c.expectedPoint.Y()
		if math.Abs(deltaX) > Epsilon || math.Abs(deltaY) > Epsilon {
			t.Log("error, should be:", c.expectedPoint, "but got:", gotPoint)
			t.Fail()
		}
	}

}
