package geocode

import (
	"testing"
)

func TestValidELR(t *testing.T) {
	reg := RegexELR()

	cases := []struct {
		elr           string
		expectedValid bool
	}{
		{
			elr:           "NBK",
			expectedValid: true},
		{
			elr:           "ECM8",
			expectedValid: true},
		{
			elr:           "",
			expectedValid: false},
		{
			elr:           "e",
			expectedValid: false},
		{
			elr:           "E",
			expectedValid: false},
		{
			elr:           "EC",
			expectedValid: false},
		{
			elr:           "NBk",
			expectedValid: false},
		{
			elr:           "nbk",
			expectedValid: false},
		{
			elr:           "ecm8",
			expectedValid: false},
		{
			elr:           "1234",
			expectedValid: false},
		{
			elr:           "EC8M",
			expectedValid: false},
		{
			elr:           "ECMx",
			expectedValid: false},
		{
			elr:           "ECM1234",
			expectedValid: false},
		{
			elr:           "%*3_",
			expectedValid: false},
	}

	for _, c := range cases {
		res := reg.MatchString(c.elr)

		if res != c.expectedValid {
			t.Errorf("Valid ELR (%s) = %t; want %t", c.elr, res, c.expectedValid)
		}
	}

}

func TestExplodeTotalYards(t *testing.T) {
	tests := []struct {
		totalYards    int
		expectedMiles int
		expectedYards int
	}{
		{ // Zero total yards.
			totalYards:    0,
			expectedMiles: 0,
			expectedYards: 0,
		},
		{ // Positive total yards less than 1 mile.
			totalYards:    1,
			expectedMiles: 0,
			expectedYards: 1,
		},
		{ // Positive total yards, nearly 1 mile.
			totalYards:    1_759,
			expectedMiles: 0,
			expectedYards: 1_759,
		},
		{
			// Positive total yards, exactly 1 mile.
			totalYards:    1_760,
			expectedMiles: 1,
			expectedYards: 0,
		},
		{
			// Positive total yards, just over 1 mile.
			totalYards:    1_761,
			expectedMiles: 1,
			expectedYards: 1,
		},
		{
			// Negative total yards 1.
			totalYards:    -1,
			expectedMiles: 0,
			expectedYards: -1,
		},
		{
			// Negative total yards 2.
			totalYards:    -1_759,
			expectedMiles: 0,
			expectedYards: -1_759,
		},
		{
			// Large negative total yards.
			totalYards:    -1_761,
			expectedMiles: 0,
			expectedYards: -1_761,
		},
		{
			// Largest negative total yards (ELR HCS: known to be erronous).
			totalYards:    -2_884,
			expectedMiles: 0,
			expectedYards: -2_884,
		},
		{
			// Large total yards.
			totalYards:    17_601,
			expectedMiles: 10,
			expectedYards: 1,
		},
		{
			// Very large total yards.
			totalYards:    176_001,
			expectedMiles: 100,
			expectedYards: 1,
		},
	}

	for _, tt := range tests {
		t.Run("test", func(t *testing.T) {
			miles, yards := ExplodeTotalYards(tt.totalYards)
			if miles != tt.expectedMiles || yards != tt.expectedYards {
				t.Errorf("ExplodeTotalYards(%v) = %v, %v, want %v, %v", tt.totalYards, miles, yards, tt.expectedMiles, tt.expectedYards)
			}
		})
	}
}

func TestBuildTotalYards(t *testing.T) {
	cases := []struct {
		miles              int
		yards              int
		expectedTotalYards int
	}{
		{
			miles:              0,
			yards:              -2_884,
			expectedTotalYards: -2_884}, // ELR HCS: known to be erronous.
		{
			miles:              0,
			yards:              -9,
			expectedTotalYards: -9},
		{
			miles:              0,
			yards:              0,
			expectedTotalYards: 0},
		{
			miles:              0,
			yards:              2,
			expectedTotalYards: 2},

		{miles: 0,
			yards:              1_759,
			expectedTotalYards: 1_759},
		{
			miles:              1,
			yards:              0,
			expectedTotalYards: 1_760},
		{
			miles:              1,
			yards:              3,
			expectedTotalYards: 1_763},
		{
			miles:              10,
			yards:              5,
			expectedTotalYards: 17_605},
		{
			miles:              100,
			yards:              1_759,
			expectedTotalYards: 177_759},
	}

	for _, c := range cases {
		res := BuildTotalYards(c.miles, c.yards)
		if res != c.expectedTotalYards {
			t.Errorf("BuildTotalYards(%v, %v) = %v, want %v", c.miles, c.yards, res, c.expectedTotalYards)
		}
	}
}

func TestMetresToMiles(t *testing.T) {
	cases := []struct {
		metres         float64
		expectedString string
	}{
		{
			metres:         1_609.344 * 0,
			expectedString: "0.000 miles"},
		{
			metres:         1_609.344 * 1,
			expectedString: "1.000 miles"},
		{
			metres:         1_609.344 * 2,
			expectedString: "2.000 miles"},
		{
			metres:         1_609.344 * 9.5,
			expectedString: "9.500 miles"},
		{
			metres:         1_609.344 * 53,
			expectedString: "53.000 miles"},
	}

	for _, c := range cases {
		res := MetresToMiles(c.metres)
		if res != c.expectedString {
			t.Errorf("MetresToMiles(%v) = %v, want %v", c.metres, res, c.expectedString)
		}
	}
}

func TestFmtTotalYards(t *testing.T) {
	cases := []struct {
		requestedTY    int
		metric         bool
		expectedString string
	}{
		// Non-metric, positive mileages.
		{
			requestedTY:    0,
			metric:         false,
			expectedString: "0M 0000y"},
		{
			requestedTY:    1,
			metric:         false,
			expectedString: "0M 0001y"},
		{
			requestedTY:    1_760,
			metric:         false,
			expectedString: "1M 0000y"},
		{
			requestedTY:    1_760 + 1_759,
			metric:         false,
			expectedString: "1M 1759y"},
		{
			requestedTY:    1_760*35 + 880,
			metric:         false,
			expectedString: "35M 0880y"},

		// Non-metric, negative mileages.
		{
			requestedTY:    -1,
			metric:         false,
			expectedString: "0M -001y"},
		{
			requestedTY:    -2_884,
			metric:         false,
			expectedString: "0M -2884y"}, // ELR HCS: known to be erronous.
		{
			requestedTY:    -1_456,
			metric:         false,
			expectedString: "0M -1456y"}, // ELR CJA2.
		{
			requestedTY:    -965,
			metric:         false,
			expectedString: "0M -965y"}, // ELR LFL.

		// Metric, positive kilometreages.
		{
			requestedTY:    0,
			metric:         true,
			expectedString: "0.000km"},
		{
			requestedTY:    1,
			metric:         true,
			expectedString: "0.001km"},
		{
			requestedTY:    1_000,
			metric:         true,
			expectedString: "0.914km"},
		{
			requestedTY:    1_760,
			metric:         true,
			expectedString: "1.609km"},
	}

	for _, c := range cases {
		res := FmtTotalYards(c.requestedTY, c.metric)
		if res != c.expectedString {
			t.Errorf("FmtTotalYards(%v, %v) = %v, want %v", c.requestedTY, c.metric, res, c.expectedString)
		}
	}
}

func TestFmtMileages(t *testing.T) {
	cases := []struct {
		tyFrom         int
		tyTo           int
		metric         bool
		expectedString string
	}{
		{
			tyFrom:         0,
			tyTo:           0,
			metric:         false,
			expectedString: "0M 0000y to 0M 0000y"},
		{
			tyFrom:         1,
			tyTo:           1,
			metric:         false,
			expectedString: "0M 0001y to 0M 0001y"},
		{
			tyFrom:         1_760,
			tyTo:           1_760,
			metric:         false,
			expectedString: "1M 0000y to 1M 0000y"},
		{
			tyFrom:         1_760 + 1,
			tyTo:           1_760 + 1_759,
			metric:         false,
			expectedString: "1M 0001y to 1M 1759y"},
		{
			tyFrom:         1_760 + 1_759,
			tyTo:           1_760 + 1_759,
			metric:         false,
			expectedString: "1M 1759y to 1M 1759y"},
		{
			tyFrom:         1_760*22 + 440,
			tyTo:           1_760*33 + 1_320,
			metric:         false,
			expectedString: "22M 0440y to 33M 1320y"},
		{
			tyFrom:         1_760*35 + 880,
			tyTo:           1_760*35 + 880,
			metric:         false,
			expectedString: "35M 0880y to 35M 0880y"},
		{
			tyFrom:         0,
			tyTo:           0,
			metric:         true,
			expectedString: "0.000km to 0.000km"},
		{
			tyFrom:         1,
			tyTo:           1,
			metric:         true,
			expectedString: "0.001km to 0.001km"},
		{
			tyFrom:         1_000,
			tyTo:           1_760,
			metric:         true,
			expectedString: "0.914km to 1.609km"},
		{
			tyFrom:         5_468,
			tyTo:           6_562,
			metric:         true,
			expectedString: "5.000km to 6.000km"}, // 5,000 / 0.9144 = 5,468   6,000 / 0.9144 = 6,562
		{
			tyFrom:         8_202,
			tyTo:           10_116,
			metric:         true,
			expectedString: "7.500km to 9.250km"}, // 7,500 / 0.9144 = 8,202   9,250 / 0.9144= 10,116
	}

	for _, c := range cases {
		res := FmtMileages(c.tyFrom, c.tyTo, c.metric)
		if res != c.expectedString {
			t.Errorf("FmtMileages(%v, %v, %v) = %v, want %v", c.tyFrom, c.tyTo, c.metric, res, c.expectedString)
		}
	}
}

func TestFmtTotalYardsMetric(t *testing.T) {
	cases := []struct {
		ty                   int
		expectedKilometreage string
	}{
		{
			ty:                   0,
			expectedKilometreage: "0.000km"},
		{
			ty:                   1_760,
			expectedKilometreage: "1.609km"},
		{
			ty:                   10_936,
			expectedKilometreage: "10.000km"}, // 10,000 / 0.9144 = 10,936.133
		{
			ty:                   17_600,
			expectedKilometreage: "16.093km"},
		{
			ty:                   176_000,
			expectedKilometreage: "160.934km"},
		{
			ty:                   5_468,
			expectedKilometreage: "5.000km"},
		{
			ty:                   10_116,
			expectedKilometreage: "9.250km"},
	}

	for _, c := range cases {
		res := FmtTotalYards(c.ty, true)
		if res != c.expectedKilometreage {
			t.Errorf("FmtTotalYardsMetric(%d) = %s; want %s", c.ty, res, c.expectedKilometreage)
		}

	}
}
