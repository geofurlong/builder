// Railway conversion functions for mileages and kilometreages.

package geocode

import (
	"fmt"
	"regexp"
)

const (
	YardsInMile      int     = 1_760     // Number of yards in a mile.
	QuarterMileYards         = 440       // Number of yards in a quarter mile.
	MetresInMile     float64 = 1_609.344 // Metres to miles conversion factor.
	YardsToMetres    float64 = 0.9144    // Yards to metres conversion factor.
)

// RegexELR returns a compiled regular expression for validating ELR codes.
func RegexELR() *regexp.Regexp {
	elrRegex, err := regexp.Compile(`[A-Z]{3}\d?$`)
	Check(err)
	return elrRegex
}

// ExplodeTotalYards takes a total yards value and returns the component miles and yards values.
func ExplodeTotalYards(totalYards int) (int, int) {
	if totalYards >= YardsInMile {
		return totalYards / YardsInMile, totalYards % YardsInMile
	}

	// For total yards less than one mile, return zero miles and the total yards as-is.
	// For edge cases where ELRs have negative mileages in excess of one mile, no attempt is made to convert to miles.
	return 0, totalYards
}

// BuildTotalYards takes miles and yards components and returns a total yards value.
func BuildTotalYards(miles int, yards int) int {
	return miles*YardsInMile + yards
}

// MetresToMiles takes a distance in metres and returns the equivalent decimal miles as a string.
func MetresToMiles(metres float64) string {
	return fmt.Sprintf("%.3f miles", metres/MetresInMile)
}

// FmtTotalYards takes a total yard value and returns a formatted miles / yards or kilometre string.
func FmtTotalYards(totalYards int, metric bool) string {
	if metric {
		return fmt.Sprintf("%.3fkm", float64(totalYards)*YardsToMetres/1_000)
	}

	miles, yards := ExplodeTotalYards((totalYards))
	return fmt.Sprintf("%dM %04dy", miles, yards)
}

// FmtMileages takes two total yards values and returns a formatted miles / yards or kilometre string.
func FmtMileages(tyFrom int, tyTo int, metric bool) string {
	return FmtTotalYards(tyFrom, metric) + " to " + FmtTotalYards(tyTo, metric)
}
