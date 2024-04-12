package geocode

import (
	"math"
	"testing"

	"github.com/paulmach/orb"
)

func TestReproject(t *testing.T) {
	cases := []struct {
		easting  float64
		northing float64
		expected orb.Point
		location string
	}{
		{485_771, 92_079, orb.Point{-0.78628983, 50.722016}, "Selsey Bill"},
		{480_479, 96_317, orb.Point{-0.86031823, 50.760875}, "Bracklesham Bay"},
		{529_740, 180_773, orb.Point{-0.13181379, 51.510997}, "Wardour St, London"},
		{496_685, 177_398, orb.Point{-0.60888972, 51.487277}, "Eton Guns, Eton"},
		{502_925, 160_554, orb.Point{-0.52393303, 51.334773}, "Sheerwater, Woking"},
		{528_132, 184_415, orb.Point{-0.15364227, 51.544097}, "Chalk Farm Tube Station"},
		{147_092, 30_263, orb.Point{-5.5392674101162, 50.1184449113554}, "Penzance"},
		{224_726, 382_330, orb.Point{-4.6321875, 53.308828}, "Holyhead"},
		{150_552, 755_298, orb.Point{-6.06809695014501, 56.622901648826}, "Tobermory"},
		{175_405, 225_296, orb.Point{-5.26464006525205, 51.8808734514259}, "St Davids"},
		{176_346, 827_264, orb.Point{-5.71229896759099, 57.2810797079495}, "Kyle of Lochalsh"},
		{182_648, 44_827, orb.Point{-5.05101265950598, 50.2633300534275}, "Truro"},
		{218_570, 597_931, orb.Point{-4.85531770044002, 55.2421123945225}, "Girvan"},
		{266_617, 845_493, orb.Point{-4.22622129478818, 57.4799977149705}, "Inverness"},
		{268_577, 551_148, orb.Point{-4.04786239360251, 54.8378524191877}, "Kirkcudbright"},
		{293_994, 758_112, orb.Point{-3.73289806042956, 56.7026446631663}, "Pitlochry"},
		{303_277, 827_846, orb.Point{-3.60817223914678, 57.3308520860594}, "Grantown-on-Spey"},
		{311_687, 968_450, orb.Point{-3.52110010140828, 58.5950124669737}, "Thurso"},
		{317_814, 745_297, orb.Point{-3.33991695114065, 56.5923527415861}, "Blairgowrie and Rattray"},
		{330_873, 436_578, orb.Point{-3.05154404779973, 53.8208602230914}, "Blackpool"},
		{350_897, 239_976, orb.Point{-2.71754978850477, 52.0559792724345}, "Hereford"},
		{390_129, 426_184, orb.Point{-2.15110286235426, 53.7319540311404}, "Portsmouth"},
		{394_426, 806_556, orb.Point{-2.09375976250924, 57.1498507970959}, "Aberdeen"},
		{402_281, 202_031, orb.Point{-1.96837565467488, 51.7170067510382}, "Cirencester"},
		{419_799, 585_952, orb.Point{-1.69074152555053, 55.1674893825777}, "Morpeth"},
		{424_572, 403_286, orb.Point{-1.63081305258923, 53.5256667596675}, "Penistone"},
		{435_749, 316_724, orb.Point{-1.47184910059996, 52.7469578946588}, "Ashby-de-la-Zouch"},
		{447_467, 1_141_434, orb.Point{-1.14681809802032, 60.1546418281584}, "Lerwick"},
		{460_218, 452_158, orb.Point{-1.0836748355092, 53.9620087173741}, "York"},
		{519_204, 298_638, orb.Point{-0.242497599724569, 52.5725093308881}, "Peterborough"},
		{567_377, 194_515, orb.Point{0.416622675810784, 51.6245616474664}, "Billericay"},
		{579_941, 183_634, orb.Point{0.592366890161814, 51.5229534756215}, "Canvey Island"},
		{622_977, 308_549, orb.Point{1.29330728438099, 52.6288686480942}, "Norwich"},
		{631_587, 141_745, orb.Point{1.30842698233744, 51.1281377122637}, "Dover"},
		{654_779, 293_216, orb.Point{1.75043311819206, 52.4773484258696}, "Lowestoft"},
	}

	const Epsilon = 1e-5

	en := []orb.Point{}
	for _, cm := range cases {
		en = append(en, orb.Point{cm.easting, cm.northing})
	}

	ll := ReprojectMulti(en)
	for i, cm := range cases {
		deltaX := ll[i].Point().X() - cm.expected.X()
		deltaY := ll[i].Point().Y() - cm.expected.Y()

		if math.Abs(deltaX) > Epsilon || math.Abs(deltaY) > Epsilon {
			t.Log("ReprojectMulti error, should be:", cm.expected, "but got:", ll[i])
			t.Fail()
		}
	}

	pj := OSGBtoLongLat()
	for _, c := range cases {
		ll := Reproject(orb.Point{c.easting, c.northing}, pj)
		deltaX := ll.X() - c.expected.X()
		deltaY := ll.Y() - c.expected.Y()
		t.Logf("%.6f %.6f", deltaX, deltaY)

		if math.Abs(deltaX) > Epsilon || math.Abs(deltaY) > Epsilon {
			t.Log("Reproject error, should be:", c.expected, "but got:", ll, "for location:", c.location)
			t.Fail()
		}
	}
}
