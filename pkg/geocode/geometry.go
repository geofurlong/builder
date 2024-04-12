// 2D Euclidean plane geometry functions to support geocoding.

package geocode

import (
	"math"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

// interolatePoint returns a point between two points, interpolated by the fractional ratio.
func interpolatePoint(point1, point2 orb.Point, ratio float64) orb.Point {
	x := point1[0] + (point2[0]-point1[0])*ratio
	y := point1[1] + (point2[1]-point1[1])*ratio
	return orb.Point{x, y}
}

// NearestPointOnLine returns the nearest point on the line to the point, and the distance to that point.
func NearestPointOnLine(line *orb.LineString, point orb.Point) (orb.Point, float64) {
	var (
		nearestPoint orb.Point
		minDistance  = math.MaxFloat64
	)

	for i := 0; i < len(*line)-1; i++ {
		start := (*line)[i]
		end := (*line)[i+1]

		np, distance := nearestPointOnSegment(start, end, point)

		if distance < minDistance {
			minDistance = distance
			nearestPoint = np
		}
	}

	return nearestPoint, minDistance
}

// nearestPointOnSegment returns the nearest point on the line segment defined by the start and end points,
// and the distance the point is from the start of the line segment.
func nearestPointOnSegment(segmentStartPoint, segmentEndPoint, targetPoint orb.Point) (orb.Point, float64) {
	// Calculate the differences in x and y co-ordinates for the start and end points of the line segment.
	deltaX := segmentEndPoint[0] - segmentStartPoint[0]
	deltaY := segmentEndPoint[1] - segmentStartPoint[1]

	// Calculate the differences in x and y co-ordinates for the target point and segment start point.
	deltaXFromStart := targetPoint[0] - segmentStartPoint[0]
	deltaYFromStart := targetPoint[1] - segmentStartPoint[1]

	// Calculate the square of the length of the line segment.
	lenSquared := deltaX*deltaX + deltaY*deltaY

	// If the length is zero (start and end are the same), return the start point and the distance to the target point.
	if lenSquared == 0.0 {
		return segmentStartPoint, math.Sqrt(deltaXFromStart*deltaXFromStart + deltaYFromStart*deltaYFromStart)
	}

	// Calculate the projection of the point onto the line segment.
	projectionRatio := (deltaXFromStart*deltaX + deltaYFromStart*deltaY) / lenSquared

	// If the projection falls before the start of the line segment, return the start point and the distance to the target point.
	if projectionRatio < 0 {
		return segmentStartPoint, math.Sqrt(deltaXFromStart*deltaXFromStart + deltaYFromStart*deltaYFromStart)
	} else if projectionRatio > 1 {
		// If the projection falls after the end of the line segment, return the end point and the distance to the target point.
		deltaX = targetPoint[0] - segmentEndPoint[0]
		deltaY = targetPoint[1] - segmentEndPoint[1]
		return segmentEndPoint, math.Sqrt(deltaX*deltaX + deltaY*deltaY)
	}

	// If the projection falls on the line segment, calculate the co-ordinates of the projection point.
	projection := orb.Point{segmentStartPoint[0] + projectionRatio*deltaX, segmentStartPoint[1] + projectionRatio*deltaY}

	// Calculate the differences in x and y co-ordinates for the target point and the projection point.
	deltaX = targetPoint[0] - projection[0]
	deltaY = targetPoint[1] - projection[1]

	// Return the projection point and the distance to the target point.
	return projection, math.Sqrt(deltaX*deltaX + deltaY*deltaY)
}

// DistanceAlongLine returns the distance along the line from the start to the point.
func DistanceAlongLine(line *orb.LineString, point orb.Point) float64 {
	totalDistance := 0.0

	for i := 0; i < len(*line)-1; i++ {
		start, end := (*line)[i], (*line)[i+1]

		if pointOnLineSegment(start, end, point) {
			totalDistance += planar.Distance(start, point)
			break
		} else {
			totalDistance += planar.Distance(start, end)
		}
	}

	return totalDistance
}

// pointOnLineSegment returns true if the point is on the line segment defined by the start and end points.
func pointOnLineSegment(startPoint, endPoint, targetPoint orb.Point) bool {
	const Epsilon = 1e-9
	segmentLength := planar.Distance(startPoint, endPoint)
	distanceToSegmentStart := planar.Distance(startPoint, targetPoint)
	distanceToSegmentEnd := planar.Distance(endPoint, targetPoint)
	return math.Abs(segmentLength-(distanceToSegmentStart+distanceToSegmentEnd)) < Epsilon
}

// pointAtDistanceAlongLine returns the point (interpolated if necessary) at the given distance (metres) along a linestring.
// The orb library does not have this function, so this is a rework of the library's spherical equivalent.
func pointAtDistanceAlongLine(distance float64, line orb.LineString) orb.Point {
	numPoints := len(line)

	if distance < 0 || numPoints == 1 {
		return line[0]
	}

	var (
		travelled          = 0.0
		fromPoint, toPoint orb.Point
	)

	for i := 1; i < numPoints; i++ {
		fromPoint, toPoint = line[i-1], line[i]
		actualSegmentDistance := planar.Distance(fromPoint, toPoint)
		expectedSegmentDistance := distance - travelled

		if expectedSegmentDistance < actualSegmentDistance {
			return interpolatePoint(fromPoint, toPoint, expectedSegmentDistance/actualSegmentDistance)
		}

		travelled += actualSegmentDistance
	}

	return toPoint
}
