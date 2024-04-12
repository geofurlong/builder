BEGIN;

CREATE TABLE gazetteer_grouping (group_id INTEGER, group_name VARCHAR NOT NULL, PRIMARY KEY (group_id));
INSERT INTO gazetteer_grouping (group_id, group_name) VALUES (1, 'nr_region');
INSERT INTO gazetteer_grouping (group_id, group_name) VALUES (2, 'country_admin_area'); -- Country and admin area name railway point is within.
INSERT INTO gazetteer_grouping (group_id, group_name) VALUES (3, 'district_place');  -- District and place name of nearest place to railway point.


CREATE TABLE gazetteer_aggregated (
    elr VARCHAR NOT NULL, 
    group_id INTEGER NOT NULL, 
    offset_from INT NOT NULL, 
    offset_to INT NOT NULL, 
    mileage_from VARCHAR NOT NULL, 
    mileage_to VARCHAR NOT NULL, 
    value_1 VARCHAR NOT NULL, 
    value_2 VARCHAR, 
    min_distance INT, 
    max_distance INT, 
    mean_distance INT
);

.mode csv
.header ON
.import "gazetteer_aggregated.csv" gazetteer_aggregated --skip 1


CREATE UNIQUE INDEX ix_aggregated ON gazetteer_aggregated (elr, group_id, offset_from, offset_to);

CREATE VIEW gazetteer_aggregated_summary AS
SELECT
    elr, group_name, offset_from, offset_to, mileage_from, mileage_to, value_1, value_2, min_distance, max_distance, mean_distance
FROM
    gazetteer_aggregated
JOIN
    gazetteer_grouping ON gazetteer_aggregated.group_id = gazetteer_grouping.group_id;


CREATE TABLE gazetteer_by_nr_region (
    elr VARCHAR NOT NULL, 
    offset_from INT NOT NULL, 
    offset_to INT NOT NULL, 
    mileage_from VARCHAR NOT NULL, 
    mileage_to VARCHAR NOT NULL, 
    nr_region VARCHAR NOT NULL, 
    PRIMARY KEY (elr, offset_from, offset_to)
);

INSERT INTO gazetteer_by_nr_region (elr, offset_from, offset_to, mileage_from, mileage_to, nr_region)
SELECT
    elr, offset_from, offset_to, mileage_from, mileage_to, value_1 
FROM
    gazetteer_aggregated 
WHERE
    group_id = (SELECT group_id FROM gazetteer_grouping WHERE group_name = 'nr_region')
ORDER BY
    elr, offset_from;


CREATE TABLE gazetteer_by_country_admin_area (
    elr VARCHAR NOT NULL, 
    offset_from INT NOT NULL, 
    offset_to INT NOT NULL, 
    mileage_from VARCHAR NOT NULL, 
    mileage_to VARCHAR NOT NULL, 
    country VARCHAR NOT NULL, 
    admin_area VARCHAR NOT NULL,
    PRIMARY KEY (elr, offset_from, offset_to)
);

INSERT INTO gazetteer_by_country_admin_area (elr, offset_from, offset_to, mileage_from, mileage_to, country, admin_area)
SELECT
    elr, offset_from, offset_to, mileage_from, mileage_to, value_1, value_2
FROM
    gazetteer_aggregated
WHERE
    group_id = (SELECT group_id FROM gazetteer_grouping WHERE group_name = 'country_admin_area')
ORDER BY
    elr, offset_from;


CREATE TABLE gazetteer_by_nearest_place (
    elr VARCHAR NOT NULL, 
    offset_from INT NOT NULL, 
    offset_to INT NOT NULL, 
    mileage_from VARCHAR NOT NULL, 
    mileage_to VARCHAR NOT NULL, 
    district VARCHAR NOT NULL, 
    place VARCHAR NOT NULL, 
    distance_min INT NOT NULL, 
    distance_max INT NOT NULL, 
    distance_mean INT NOT NULL, 
    PRIMARY KEY (elr, offset_from, offset_to)
);

INSERT INTO gazetteer_by_nearest_place (elr, offset_from, offset_to, mileage_from, mileage_to, district, place, distance_min, distance_max, distance_mean)
SELECT
    elr, offset_from, offset_to, mileage_from, mileage_to, value_1, value_2, min_distance, max_distance, mean_distance
FROM
    gazetteer_aggregated 
WHERE
    group_id = (SELECT group_id FROM gazetteer_grouping WHERE group_name = 'district_place')
ORDER BY
    elr, offset_from;


CREATE TABLE elr_by_country_admin_area AS
SELECT country, admin_area, GROUP_CONCAT(elr, ";") AS elrs
FROM (
	SELECT DISTINCT country, admin_area, elr
	FROM gazetteer_by_country_admin_area
)
GROUP BY country, admin_area;

CREATE INDEX ix_admin_area ON elr_by_country_admin_area(admin_area);


CREATE TABLE elr_by_nearest_place AS
SELECT district, place, GROUP_CONCAT(elr, ";") AS elrs
FROM (
	SELECT DISTINCT district, place, elr
	FROM gazetteer_by_nearest_place
)
GROUP BY district, place;

CREATE INDEX ix_place_name ON elr_by_nearest_place(place);


COMMIT;

ANALYZE;
VACUUM;

