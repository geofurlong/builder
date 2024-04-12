package geocode

import (
	"testing"

	"github.com/paulmach/orb"
)

type testValue struct {
	easting  int
	northing int
	osgr     string
	location string
}

func getTestCases() []testValue {
	return []testValue{
		{485_771, 92_079, "SZ8577192079", "Selsey Bill"},
		{480_479, 96_317, "SZ8047996317", "Bracklesham Bay"},
		{529_740, 180_773, "TQ2974080773", "Wardour St, London"},
		{496_685, 177_398, "SU9668577398", "Eton Guns, Eton"},
		{502_925, 160_554, "TQ0292560554", "Sheerwater, Woking"},
		{528_132, 184_415, "TQ2813284415", "Chalk Farm Tube Station"},
		{147_092, 30_263, "SW4709230263", "Penzance"},
		{224_726, 382_330, "SH2472682330", "Holyhead"},
		{150_552, 755_297, "NM5055255297", "Tobermory"},
		{175_405, 225_295, "SM7540525295", "St Davids"},
		{176_346, 827_263, "NG7634627263", "Kyle of Lochalsh"},
		{182_648, 44_827, "SW8264844827", "Truro"},
		{218_570, 597_931, "NX1857097931", "Girvan"},
		{266_617, 845_493, "NH6661745493", "Inverness"},
		{268_577, 551_148, "NX6857751148", "Kirkcudbright"},
		{293_994, 758_112, "NN9399458112", "Pitlochry"},
		{303_277, 827_846, "NJ0327727846", "Grantown-on-Spey"},
		{311_687, 968_450, "ND1168768450", "Thurso"},
		{317_814, 745_296, "NO1781445296", "Blairgowrie and Rattray"},
		{330_872, 436_578, "SD3087236578", "Blackpool"},
		{350_897, 239_976, "SO5089739976", "Hereford"},
		{390_129, 426_184, "SD9012926184", "Portsmouth"},
		{394_426, 806_556, "NJ9442606556", "Aberdeen"},
		{402_281, 202_031, "SP0228102031", "Cirencester"},
		{419_799, 585_952, "NZ1979985952", "Morpeth"},
		{424_572, 403_285, "SE2457203285", "Penistone"},
		{435_749, 316_723, "SK3574916723", "Ashby-de-la-Zouch"},
		{447_467, 1_141_434, "HU4746741434", "Lerwick"},
		{460_218, 452_158, "SE6021852158", "York"},
		{519_204, 298_638, "TL1920498638", "Peterborough"},
		{567_377, 194_515, "TQ6737794515", "Billericay"},
		{579_941, 183_634, "TQ7994183634", "Canvey Island"},
		{622_977, 308_549, "TG2297708549", "Norwich"},
		{631_587, 141_745, "TR3158741745", "Dover"},
		{654_779, 293_216, "TM5477993216", "Lowestoft"},
	}
}

func TestPointToOSGR(t *testing.T) {
	testCases := getTestCases()

	for _, tc := range testCases {
		res := PointToOSGR(orb.Point{float64(tc.easting), float64(tc.northing)})

		if res != tc.osgr {
			t.Errorf("PointToOSGR(%d, %d) = %s; want %s", tc.easting, tc.northing, res, tc.osgr)
		}
	}

}

func TestOSGRToPoint(t *testing.T) {
	testCases := getTestCases()

	for _, tc := range testCases {
		pt, err := OSGRToPoint(tc.osgr)
		if err != nil {
			t.Errorf("OSGRToPoint(%q) returned error: %v", tc.osgr, err)
		} else if int(pt.X()) != tc.easting || int(pt.Y()) != tc.northing {
			t.Errorf("OSGRToPoint(%q) = %d, %d; want %d, %d", tc.osgr, int(pt.X()), int(pt.Y()), tc.easting, tc.northing)
		}
	}
}
