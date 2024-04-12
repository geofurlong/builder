import convert

# Note: All Metric ELRs have positive kilometreages.


def test_miles_yards_to_total_yards():
    assert convert.miles_yards_to_total_yards(-0.1) == -1_000
    assert convert.miles_yards_to_total_yards(-0.01) == -100
    assert convert.miles_yards_to_total_yards(-0.001) == -10
    assert convert.miles_yards_to_total_yards(-0.0001) == -1
    assert convert.miles_yards_to_total_yards(0) == 0
    assert convert.miles_yards_to_total_yards(0.0001) == 1
    assert convert.miles_yards_to_total_yards(0.0002) == 2
    assert convert.miles_yards_to_total_yards(0.001) == 10
    assert convert.miles_yards_to_total_yards(0.01) == 100
    assert convert.miles_yards_to_total_yards(0.088) == 880
    assert convert.miles_yards_to_total_yards(0.1) == 1_000
    assert convert.miles_yards_to_total_yards(0.132) == 1_320
    assert convert.miles_yards_to_total_yards(1) == 1_760
    assert convert.miles_yards_to_total_yards(10) == 17_600
    assert convert.miles_yards_to_total_yards(99.1759) == 99 * 1_760 + 1_759
    assert convert.miles_yards_to_total_yards(123.0456) == 123 * 1_760 + 456


def test_kilometres_to_total_yards():
    assert convert.km_to_total_yards(0) == 1_760 * 0
    assert convert.km_to_total_yards(1.609344) == 1_760 * 1
    assert convert.km_to_total_yards(8.04672) == 1_760 * 5
    assert convert.km_to_total_yards(16.09344) == 1_760 * 10
    assert convert.km_to_total_yards(160.9344) == 1_760 * 100
    assert convert.km_to_total_yards(198.753984) == 1_760 * 123.5


def test_to_total_yards():
    assert convert.to_total_yards("M", -0.1) == -1_000
    assert convert.to_total_yards("M", 0) == 0
    assert convert.to_total_yards("M", 0.0001) == 1
    assert convert.to_total_yards("M", 0.001) == 10
    assert convert.to_total_yards("M", 0.01) == 100
    assert convert.to_total_yards("M", 0.1) == 1_000
    assert convert.to_total_yards("M", 1) == 1_760
    assert convert.to_total_yards("M", 1.0001) == 1_761
    assert convert.to_total_yards("M", 10.1759) == 17_600 + 1759
    assert convert.to_total_yards("M", 100.1759) == 176_000 + 1759

    # Metric linear reporting units.
    assert convert.to_total_yards("K", 0) == 0
    assert convert.to_total_yards("K", 1.609344 / 2.0) == 880
    assert convert.to_total_yards("K", 1.609344) == 1_760
