// Conversion functions between Easting / Northing points and Ordnance Survey Grid References (OSGR).

package geocode

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/paulmach/orb"
)

const (
	// Size of the single and double letter top-level OSGR tiles.
	PrimaryTileSize   = 500_000
	SecondaryTileSize = 100_000

	// Origins of the OS grid.
	GridOriginSWEasting  = -2 * PrimaryTileSize
	GridOriginSWNorthing = -PrimaryTileSize

	// Order of the lettered OSGR tiles.
	TileLetterOrder = "VWXYZQRSTULMNOPFGHJKABCDE"

	// Number of tiles per row in the lettered OSGR grid.
	TilesPerRow = 5
)

// letterOrigin represents an OSGR tile letter and associated tile origin co-ordinates.
type letterOrigin struct {
	letter string
	origin orb.Point
}

// tileAndOrigin returns the OSGR tile letter and associated tile origin co-ordinates for an Easting / Northing point.
func tileAndOrigin(pt orb.Point, tileSize int) letterOrigin {
	tilePosX := math.Floor(pt.X() / float64(tileSize))
	tilePosY := math.Floor(pt.Y() / float64(tileSize))
	letter := string(TileLetterOrder[int(tilePosX)+TilesPerRow*int(tilePosY)])
	originX := tilePosX * float64(tileSize)
	originY := tilePosY * float64(tileSize)
	return letterOrigin{letter, orb.Point{originX, originY}}
}

// tile100kmAndOrigin returns the OSGR tile letters and origin co-ordinates for an Easting / Northing point.
func tile100kmAndOrigin(pt orb.Point) letterOrigin {
	lo1 := tileAndOrigin(orb.Point{pt.X() - GridOriginSWEasting, pt.Y() - GridOriginSWNorthing}, PrimaryTileSize)
	lo2 := tileAndOrigin(orb.Point{pt.X() - lo1.origin.X() - GridOriginSWEasting, pt.Y() - lo1.origin.Y() - GridOriginSWNorthing}, SecondaryTileSize)
	letters := lo1.letter + lo2.letter
	originX := lo1.origin.X() + lo2.origin.X() + GridOriginSWEasting
	originY := lo1.origin.Y() + lo2.origin.Y() + GridOriginSWNorthing
	return letterOrigin{letters, orb.Point{originX, originY}}
}

// PointToOSGR returns the OSGR string (to 1 metre resolution) for the given Easting / Northing point.
func PointToOSGR(pt orb.Point) string {
	lo := tile100kmAndOrigin(pt)
	coordX := int(math.Floor(pt.X() - lo.origin.X()))
	coordY := int(math.Floor(pt.Y() - lo.origin.Y()))
	return fmt.Sprintf("%s%05d%05d", lo.letter, coordX, coordY)
}

// OSGRToPoint returns the Easting / Northing point for the given OSGR string.
func OSGRToPoint(osgr string) (orb.Point, error) {
	prefix := osgr[:2]

	numbers := osgr[2:]
	numbersMidPoint := len(numbers) / 2
	n1, err := strconv.Atoi(numbers[:numbersMidPoint])
	if err != nil {
		return orb.Point{0, 0}, err
	}

	n2, err := strconv.Atoi(numbers[numbersMidPoint:])
	if err != nil {
		return orb.Point{0, 0}, err
	}

	ixLetter1 := strings.Index(TileLetterOrder, string(prefix[0]))
	originX := PrimaryTileSize*(ixLetter1%TilesPerRow) + GridOriginSWEasting
	originY := PrimaryTileSize*(ixLetter1/TilesPerRow) + GridOriginSWNorthing

	ixLetter2 := strings.Index(TileLetterOrder, string(prefix[1]))
	originX += SecondaryTileSize * (ixLetter2 % TilesPerRow)
	originY += SecondaryTileSize * (ixLetter2 / TilesPerRow)

	depth := (len(osgr) - 2) / 2
	cellSize := SecondaryTileSize / int(math.Pow10(depth))
	originX += n1 * cellSize
	originY += n2 * cellSize

	return orb.Point{float64(originX), float64(originY)}, nil
}
