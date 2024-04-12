package geocode

import (
	"testing"
)

func TestValidELR(t *testing.T) {
	reg := RegexELR()

	cases := []struct {
		elr      string
		expected bool
	}{
		{"NBK", true},
		{"ECM8", true},

		{"", false},
		{"e", false},
		{"E", false},
		{"EC", false},
		{"NBk", false},
		{"nbk", false},
		{"ecm8", false},
		{"1234", false},
		{"EC8M", false},
		{"ECMx", false},
		{"ECM1234", false},
		{"%*3_", false},
	}

	for _, c := range cases {
		res := reg.MatchString(c.elr)

		if res != c.expected {
			t.Errorf("Valid ELR (%s) = %t; want %t", c.elr, res, c.expected)
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
		miles    int
		yards    int
		expected int
	}{
		{0, -2_884, -2_884}, // ELR HCS: known to be erronous.
		{0, -9, -9},
		{0, 0, 0},
		{0, 2, 2},
		{0, 1_759, 1_759},
		{1, 0, 1_760},
		{1, 3, 1_763},
		{10, 5, 17_605},
		{100, 1_759, 177_759},
	}

	for _, c := range cases {
		res := BuildTotalYards(c.miles, c.yards)
		if res != c.expected {
			t.Errorf("BuildTotalYards(%v, %v) = %v, want %v", c.miles, c.yards, res, c.expected)
		}
	}
}

func TestMetresToMiles(t *testing.T) {
	cases := []struct {
		metres   float64
		expected string
	}{
		{1_609.344 * 0, "0.000 miles"},
		{1_609.344 * 1, "1.000 miles"},
		{1_609.344 * 2, "2.000 miles"},
		{1_609.344 * 9.5, "9.500 miles"},
		{1_609.344 * 53, "53.000 miles"},
	}

	for _, c := range cases {
		res := MetresToMiles(c.metres)
		if res != c.expected {
			t.Errorf("MetresToMiles(%v) = %v, want %v", c.metres, res, c.expected)
		}
	}
}

func TestFmtTotalYards(t *testing.T) {
	cases := []struct {
		requestedTY int
		metric      bool
		expected    string
	}{
		{0, false, "0M 0000y"},
		{1, false, "0M 0001y"},
		{1_760, false, "1M 0000y"},
		{1_760 + 1_759, false, "1M 1759y"},
		{1_760*35 + 880, false, "35M 0880y"},

		{-1, false, "0M -001y"},
		{-2_884, false, "0M -2884y"}, // ELR HCS: known to be erronous.
		{-1_456, false, "0M -1456y"}, // ELR CJA2.
		{-965, false, "0M -965y"},    // ELR LFL.

		{0, true, "0.000km"},
		{1, true, "0.001km"},
		{1_000, true, "0.914km"},
		{1_760, true, "1.609km"},
	}

	for _, c := range cases {
		res := FmtTotalYards(c.requestedTY, c.metric)
		if res != c.expected {
			t.Errorf("FmtTotalYards(%v, %v) = %v, want %v", c.requestedTY, c.metric, res, c.expected)
		}
	}
}

func TestFmtMileages(t *testing.T) {
	cases := []struct {
		tyFrom   int
		tyTo     int
		metric   bool
		expected string
	}{
		{0, 0, false, "0M 0000y to 0M 0000y"},
		{1, 1, false, "0M 0001y to 0M 0001y"},
		{1_760, 1_760, false, "1M 0000y to 1M 0000y"},
		{1_760 + 1, 1_760 + 1_759, false, "1M 0001y to 1M 1759y"},
		{1_760 + 1_759, 1_760 + 1_759, false, "1M 1759y to 1M 1759y"},
		{1_760*22 + 440, 1_760*33 + 1_320, false, "22M 0440y to 33M 1320y"},
		{1_760*35 + 880, 1_760*35 + 880, false, "35M 0880y to 35M 0880y"},
		{0, 0, true, "0.000km to 0.000km"},
		{1, 1, true, "0.001km to 0.001km"},
		{1_000, 1_760, true, "0.914km to 1.609km"},
		{5_468, 6_562, true, "5.000km to 6.000km"},  // 5,000 / 0.9144 = 5,468   6,000 / 0.9144 = 6,562
		{8_202, 10_116, true, "7.500km to 9.250km"}, // 7,500 / 0.9144 = 8,202   9,250 / 0.9144= 10,116
	}

	for _, c := range cases {
		res := FmtMileages(c.tyFrom, c.tyTo, c.metric)
		if res != c.expected {
			t.Errorf("FmtMileages(%v, %v, %v) = %v, want %v", c.tyFrom, c.tyTo, c.metric, res, c.expected)
		}
	}
}

func TestFmtTotalYardsMetric(t *testing.T) {
	cases := []struct {
		ty       int
		expected string
	}{
		{0, "0.000km"},
		{1_760, "1.609km"},
		{10_936, "10.000km"}, // 10,000 / 0.9144 = 10,936.133
		{17_600, "16.093km"},
		{176_000, "160.934km"},
		{5_468, "5.000km"},
		{10_116, "9.250km"},
	}

	for _, c := range cases {
		res := FmtTotalYards(c.ty, true)
		if res != c.expected {
			t.Errorf("FmtTotalYardsMetric(%d) = %s; want %s", c.ty, res, c.expected)
		}

	}
}
