PRAGMA cache_size = 35000;

BEGIN TRANSACTION;

CREATE TABLE gazetteer (
	elr VARCHAR NOT NULL,
	total_yards INTEGER NOT NULL,
	mileage VARCHAR NOT NULL,
	easting VARCHAR NOT NULL,
	northing VARCHAR NOT NULL,
	longitude VARCHAR NOT NULL,
	latitude VARCHAR NOT NULL,
	osgr VARCHAR NOT NULL,
	accuracy INTEGER NOT NULL,
	nr_region VARCHAR NULL,
	place_name VARCHAR NOT NULL,
	district VARCHAR NULL,
	county VARCHAR NULL,
	distance_m NUMBER NOT NULL,
	country VARCHAR NOT NULL,
	admin_area VARCHAR NOT NULL
);

-- Import the CSV file created by the gazetteer script into the gazetteer table.
-- NOTE: the SQlite `.import` command does not honour the `NOT NULL` column constraints.
.mode csv
.header OFF

-- "geofurlong_gazetteer_0022y_raw.csv" is replaced with its full path by the calling program.
.import -skip 1 geofurlong_gazetteer_0022y_raw.csv gazetteer

-- ORDNANCE SURVEY COUNTY / DISTRICT.
-- county (unitary) OR district can be blank, so combine into a single field to enhance reporting.
ALTER TABLE gazetteer ADD COLUMN county_district VARCHAR;
UPDATE gazetteer SET county_district = county WHERE district = '';
UPDATE gazetteer SET county_district = district WHERE county = '';
UPDATE gazetteer SET county_district = county || ' - ' || district WHERE county <> '' AND district <> '';

-- county and district columns now superseded, as combined above.
ALTER TABLE gazetteer DROP COLUMN county;
ALTER TABLE gazetteer DROP COLUMN district;

-- Index created to improve performance on SQL updates in subsequent sections.
-- From testing, no performance benefit gained by creating temporary indexes for subsequent updates.
CREATE INDEX ix_elr ON gazetteer (elr, total_yards);

-- ORDNANCE SURVEY POPULATED PLACE NAME.
-- Updates in section below generally based on manual review and local knowledge.
UPDATE gazetteer SET place_name = 'Prestonpans' WHERE elr='ECM8' and place_name='Preston' and total_yards < 11*1760;
UPDATE gazetteer SET place_name = 'Strathcarron' WHERE elr='KYL' and place_name='Srath Carrann';
UPDATE gazetteer SET place_name = 'Duirinish' WHERE elr='KYL' and place_name='Dìurinis';
UPDATE gazetteer SET place_name = 'Duncraig' WHERE elr='KYL' and place_name='Craig' and total_yards > 55*1760;
UPDATE gazetteer SET place_name = 'Stromeferry' WHERE elr='KYL' and place_name='Strome Ferry';
UPDATE gazetteer SET place_name = 'Leven' WHERE elr IN ('MTL1', 'MTL2') and place_name='Innerleven';
UPDATE gazetteer SET place_name = 'Alness' WHERE elr='WCK' and place_name='Shillinghill';

-- ORDNANCE SURVEY ADMINISTRATIVE AREA.
-- Updates in section below generally caused by locations being close to the coastline.
UPDATE gazetteer SET admin_area='South Ayrshire' WHERE elr='AYH1' AND admin_area='';  -- Ayr Harbour
UPDATE gazetteer SET admin_area='North Ayrshire' WHERE elr='AYR4' AND admin_area='';  -- ~Troon
UPDATE gazetteer SET admin_area='North Ayrshire' WHERE admin_area='North Ayshire';  -- typo within OS source data
UPDATE gazetteer SET admin_area='Southampton' WHERE elr='BML1' AND admin_area='';  -- Southampton
UPDATE gazetteer SET admin_area='South Gloucestershire' WHERE elr='BSW' AND admin_area=''; -- and county='South Gloucestershire';
UPDATE gazetteer SET admin_area='Cumbria' WHERE elr='CBC1' AND admin_area='';  -- Cumbrian Coast (south)
UPDATE gazetteer SET admin_area='Cumbria' WHERE elr='CBC2' AND admin_area='';  -- Cumbrian Coast (north)
UPDATE gazetteer SET admin_area='Wandsworth' WHERE elr='CKL' AND admin_area='Kensington and Chelsea';
UPDATE gazetteer SET admin_area='Flintshire'   WHERE elr='CNH3' AND admin_area='' AND total_yards BETWEEN 181*1760 AND 203*1760;
UPDATE gazetteer SET admin_area='Denbighshire' WHERE elr='CNH3' AND admin_area='' AND total_yards BETWEEN 203*1760 AND 210*1760;
UPDATE gazetteer SET admin_area='Conwy'        WHERE elr='CNH3' AND admin_area='' AND total_yards BETWEEN 210*1760 AND 232*1760;
UPDATE gazetteer SET admin_area='Gwynedd'      WHERE elr='CNH3' AND admin_area='' AND total_yards BETWEEN 232*1760 AND 241.25*1760;
UPDATE gazetteer SET admin_area='Anglesey'     WHERE elr='CNH3' AND admin_area='' AND total_yards BETWEEN 241.25*1760 AND 261*1760;
UPDATE gazetteer SET admin_area='Gwynedd'      WHERE elr='CNH3' AND admin_area='' AND total_yards > 261*1760;
UPDATE gazetteer SET admin_area='Plymouth' WHERE elr='CWR' AND admin_area='';  -- Cattewater Branch (Plymouth)
UPDATE gazetteer SET admin_area='Devon'    WHERE elr='DAC' AND admin_area='' AND total_yards <  224.75*1760;
UPDATE gazetteer SET admin_area='Plymouth' WHERE elr='DAC' AND admin_area='' AND total_yards >= 224.75*1760;
UPDATE gazetteer SET admin_area='Plymouth' WHERE elr='DAC' AND admin_area=''; -- AND county='City of Plymouth';
UPDATE gazetteer SET admin_area='Gwynedd' WHERE elr='DJP' AND admin_area='';  -- Barnmouth
UPDATE gazetteer SET admin_area='Northumberland' WHERE elr='ECM7' AND admin_area='Scottish Borders';  -- SCO/ENG border
UPDATE gazetteer SET admin_area='River Forth' WHERE elr='ECN2' AND admin_area='' AND total_yards < 20000;
UPDATE gazetteer SET admin_area='River Tay'   WHERE elr='ECN2' AND admin_area='' AND total_yards > 99000;
UPDATE gazetteer SET admin_area='Angus' WHERE elr='ECN4' AND admin_area='';  -- Montrose etc.
UPDATE gazetteer SET admin_area='Cornwall' WHERE elr='FAL3' AND admin_area='';  -- Falmouth Docks
UPDATE gazetteer SET admin_area='Kent' WHERE elr='FFH2' AND admin_area='';  -- Folkestone Harbour
UPDATE gazetteer SET admin_area='Pembrokeshire' WHERE elr='FSH' AND admin_area='';  -- Fishguard Harbour
UPDATE gazetteer SET admin_area='Renfrewshire' WHERE elr='GOU2' AND admin_area='';  -- Paisley to Gourock
UPDATE gazetteer SET admin_area='Poole' WHERE elr='HAG' AND admin_area='';  -- Hamworthy
UPDATE gazetteer SET admin_area='North Ayrshire' WHERE elr='HUN' AND admin_area='';  -- Hunterston
UPDATE gazetteer SET admin_area='Isle of Wight' WHERE elr='IOW' AND admin_area='';
UPDATE gazetteer SET admin_area='Fife' WHERE elr='KNE1' AND admin_area='';  -- Kincardine
UPDATE gazetteer SET admin_area='Highland' WHERE elr='KYL' AND admin_area='';
UPDATE gazetteer SET admin_area='Cornwall' WHERE elr='LOF' AND admin_area='';
UPDATE gazetteer SET admin_area='North Ayrshire' WHERE elr='LGS1' AND admin_area='';  -- Saltcoats
UPDATE gazetteer SET admin_area='North East Lincolnshire' WHERE elr='MAC3' AND admin_area=''; -- Grimsby - TODO check inconsistencies admin vs county
UPDATE gazetteer SET admin_area='Somerset' WHERE elr='MIN' AND admin_area='';  -- Minehead
-- TODO UPDATE gazetteer SET admin_area='Merseyside' WHERE elr='MIR2' AND admin_area='' and district='Liverpool';
UPDATE gazetteer SET admin_area='Halton' WHERE elr='MIR2' AND admin_area=''; -- and district='Wirral';  | -- Birkenhead / Liverpool
UPDATE gazetteer SET admin_area='Highland' WHERE elr='MLG2' AND admin_area='';
UPDATE gazetteer SET admin_area='Devon' WHERE elr='MLN1' AND admin_area='';
UPDATE gazetteer SET admin_area='Plymouth' WHERE elr='MLN2' AND admin_area='' and county_district='City of Plymouth';
UPDATE gazetteer SET admin_area='Cornwall' WHERE elr='MLN2' AND admin_area='' and county_district='Cornwall';
UPDATE gazetteer SET admin_area='North East Lincolnshire' WHERE elr='MWN'; -- AND admin_area='';  |  -- Cleethorpes
UPDATE gazetteer SET admin_area='Argyll and Bute' WHERE elr='NEM7' AND admin_area='';  -- Cardross
UPDATE gazetteer SET admin_area='Merseyside' WHERE elr='PJL' AND admin_area='Manchester' AND total_yards > 800;
UPDATE gazetteer SET admin_area='Cumbria' WHERE elr='RDK1' AND admin_area='';  -- Barrow-in-Furness
-- UPDATE gazetteer SET country ='England' WHERE elr='SBA1' AND admin_area='Shropshire' AND country='Wales';
UPDATE gazetteer SET admin_area='Gwynedd' WHERE elr='SBA2' AND admin_area='';
UPDATE gazetteer SET admin_area='Nottinghamshire' WHERE elr='SCB' AND admin_area='Derbyshire' AND total_yards < 246950;
UPDATE gazetteer SET admin_area='Kent' WHERE elr='SEJ2' AND admin_area='';   -- Sheerness
UPDATE gazetteer SET admin_area='Dundee' WHERE elr='SCM5' AND admin_area='';  -- Dundee
UPDATE gazetteer SET admin_area='Cornwall' WHERE elr='SIV' AND admin_area='';  -- ~St Ives
UPDATE gazetteer SET admin_area='Dumfries and Galloway' WHERE elr='STR4' AND admin_area='';  -- Stranraer
UPDATE gazetteer SET admin_area='Swansea'         WHERE elr='SWM2' AND admin_area='' AND total_yards < 391000;
UPDATE gazetteer SET admin_area='Carmarthenshire' WHERE elr='SWM2' AND admin_area='' AND total_yards > 417000;
-- UPDATE gazetteer SET country ='England' WHERE elr='SWM2' AND admin_area='Gloucestershire' AND country='Wales';
UPDATE gazetteer SET admin_area='Newham' WHERE elr='TAH3' AND admin_area='Redbridge';
UPDATE gazetteer SET admin_area='Newham' WHERE elr='TLL' AND admin_area='Redbridge';
UPDATE gazetteer SET admin_area='Highland' WHERE elr='WCK' AND admin_area='';
UPDATE gazetteer SET admin_area='Portsmouth' WHERE elr='WPH2' AND admin_area='';  -- Broadmarsh
UPDATE gazetteer SET admin_area='Cheshire' WHERE elr='WSJ2' AND total_yards >= 210*1760;  -- Chester
UPDATE gazetteer SET admin_area='Fife' WHERE elr='ZZG1' AND admin_area='';  -- Kincardine

-- COUNTY / DISTRICT (from Ordnance Survey).
-- Updates in section below set Welsh county / district names to ENG version.
UPDATE gazetteer SET county_district='Blaenau Gwent ' WHERE county_district='Blaenau Gwent - Blaenau Gwent';
UPDATE gazetteer SET county_district='Bridgend' WHERE county_district='Pen-y-bont ar Ogwr - Bridgend';
UPDATE gazetteer SET county_district='Caerphilly' WHERE county_district='Caerffili - Caerphilly';
UPDATE gazetteer SET county_district='Cardiff' WHERE county_district='Caerdydd - Cardiff';
UPDATE gazetteer SET county_district='Carmarthenshire' WHERE county_district='Sir Gaerfyrddin - Carmarthenshire';
UPDATE gazetteer SET county_district='Ceredigion' WHERE county_district='Sir Ceredigion - Ceredigion';
UPDATE gazetteer SET county_district='Conwy ' WHERE county_district='Conwy - Conwy';
UPDATE gazetteer SET county_district='Denbighshire' WHERE county_district='Sir Ddinbych - Denbighshire';
UPDATE gazetteer SET county_district='Flintshire' WHERE county_district='Sir y Fflint - Flintshire';
UPDATE gazetteer SET county_district='Gwynedd ' WHERE county_district='Gwynedd - Gwynedd';
UPDATE gazetteer SET county_district='Isle of Anglesey' WHERE county_district='Sir Ynys Mon - Isle of Anglesey';
UPDATE gazetteer SET county_district='Merthyr Tydfil' WHERE county_district='Merthyr Tudful - Merthyr Tydfil';
UPDATE gazetteer SET county_district='Monmouthshire' WHERE county_district='Sir Fynwy - Monmouthshire';
UPDATE gazetteer SET county_district='Neath Port Talbot' WHERE county_district='Castell-nedd Port Talbot - Neath Port Talbot';
UPDATE gazetteer SET county_district='Newport' WHERE county_district='Casnewydd - Newport';
UPDATE gazetteer SET county_district='Pembrokeshire' WHERE county_district='Sir Benfro - Pembrokeshire';
UPDATE gazetteer SET county_district='Powys ' WHERE county_district='Powys - Powys';
UPDATE gazetteer SET county_district='Rhondda Cynon Taf ' WHERE county_district='Rhondda Cynon Taf - Rhondda Cynon Taf';
UPDATE gazetteer SET county_district='Swansea' WHERE county_district='Abertawe - Swansea';
UPDATE gazetteer SET county_district='Torfaen' WHERE county_district='Tor-faen - Torfaen';
UPDATE gazetteer SET county_district='Vale of Glamorgan' WHERE county_district='Bro Morgannwg - the Vale of Glamorgan';
UPDATE gazetteer SET county_district='Wrexham' WHERE county_district='Wrecsam - Wrexham';

-- COUNTRY (from Ordnance Survey).
-- Updates in section below generally located on or close to coastline.
UPDATE gazetteer SET country='Scotland' WHERE elr='AYH1' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='AYR4' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='BML1' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='BSW' AND country='';  -- TODO check if Wales 13.0440 - 13.0990
UPDATE gazetteer SET country='England' WHERE elr='CBC1' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='CBC2' AND country='';
UPDATE gazetteer SET country='Wales' WHERE elr='CNH3' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='CWR' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='DAC' AND country='';
UPDATE gazetteer SET country='Wales' WHERE elr='DJP' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='ECM7' AND country='Scotland';
UPDATE gazetteer SET country='Scotland' WHERE elr='ECN2' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='ECN4' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='FAL3' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='FFH2' AND country='';
UPDATE gazetteer SET country='Wales' WHERE elr='FSH' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='GOU2' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='HAG' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='HUN' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='IOW' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='KYL' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='KNE1' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='LGS1' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='LOF' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='MAC3' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='MIN' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='MIR2' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='MLG2' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='MLN1' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='MLN2' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='MWN' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='NEM7' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='RDK1' AND country='';
UPDATE gazetteer SET country='Wales' WHERE elr='SBA2' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='SEJ2' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='SCM5' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='SIV' AND country='';
UPDATE gazetteer SET country='Scotland' WHERE elr='STR4' AND country='';
UPDATE gazetteer SET country='Wales' WHERE elr='SWM2' AND country='' AND admin_area IN ('Swansea', 'Carmarthenshire');
UPDATE gazetteer SET country='Scotland' WHERE elr='WCK' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='WPH2' AND country='';
UPDATE gazetteer SET country='England' WHERE elr='WSJ2' AND total_yards >= 210*1760;  -- Chester
UPDATE gazetteer SET country='Scotland' WHERE elr='ZZG1' AND country='';

-- Gazetteer: NR REGION.
-- Updates in section below at coastlines and interfaces with other NR Regions.
UPDATE gazetteer SET nr_region='Southern' WHERE elr='AGW' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='ATG' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Scotland' WHERE elr='AYH1' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='BGK' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Southern' WHERE elr='BLP' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr IN ('BOK1' ,'BOK2', 'BOK3', 'BOK4') AND nr_region IN ('North West & Central', 'Southern');
UPDATE gazetteer SET nr_region='Wales & Western' WHERE elr='BRB' AND nr_region='Southern';
-- TODO review BOK5 (Eastern / NW&C)
-- TODO and BOK6
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='BRI2' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='CAW' AND nr_region='North West & Central';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='CRF3' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='DWW2' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='ECM1' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='ECM7' AND nr_region='Scotland';
UPDATE gazetteer SET nr_region='Scotland' WHERE elr='ECM8' AND nr_region='Eastern';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='ELL5' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Wales & Western' WHERE elr='FAL3' AND nr_region='';
UPDATE gazetteer SET nr_region='Southern' WHERE elr='FFH2' AND nr_region='';
UPDATE gazetteer SET nr_region='Wales & Western' WHERE elr='FSH' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='GLT' AND nr_region='';
UPDATE gazetteer SET nr_region='Southern' WHERE elr='HAG' AND nr_region='';
UPDATE gazetteer SET nr_region='Southern' WHERE elr='IOW' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='IPD' AND nr_region='';
UPDATE gazetteer SET nr_region='Scotland' WHERE elr='KYL' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='LTN1' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='MCL' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='MEB1' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Wales & Western' WHERE elr='MLN1' AND nr_region IN ('Eastern', 'Southern');
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='NAY' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='NKE1' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Wales & Western' WHERE elr='OOC1' AND nr_region='Eastern';
UPDATE gazetteer SET nr_region='North West & Central' WHERE elr='OOS' AND nr_region='Eastern';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='RBY' AND nr_region='';
UPDATE gazetteer SET nr_region='North West & Central' WHERE elr='RDK1' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='SAR2' AND nr_region='Southern';
UPDATE gazetteer SET nr_region='Wales & Western' WHERE elr='SBK2' AND nr_region='';
UPDATE gazetteer SET nr_region='Southern' WHERE elr='SOY' AND nr_region='';
UPDATE gazetteer SET nr_region='Scotland' WHERE elr='STR4' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='THN' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='TIR' AND nr_region='';
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='TLL' AND nr_region='Southern';

-- WMB is a short ELR which straddles three NR Regions.
UPDATE gazetteer SET nr_region='Southern' WHERE elr='WMB' AND total_yards <= -44;
UPDATE gazetteer SET nr_region='North West & Central' WHERE elr='WMB' AND total_yards BETWEEN -43 AND 197;
UPDATE gazetteer SET nr_region='Eastern' WHERE elr='WMB' AND total_yards >= 198;

UPDATE gazetteer SET nr_region='Southern' WHERE elr='WLL9' AND nr_region='Wales & Western';
UPDATE gazetteer SET nr_region='North West & Central' WHERE elr='WMD1' AND nr_region='Eastern';
UPDATE gazetteer SET nr_region='Southern' WHERE elr='WPH2' AND nr_region='';
UPDATE gazetteer SET nr_region='Scotland' WHERE elr='WYS' AND nr_region='';
UPDATE gazetteer SET nr_region='Southern' WHERE elr IN ('RDO1', 'TRL1', 'TRL2', 'TRL3') AND nr_region='Eastern'; -- CTRL (HS1)

COMMIT;


-- Database normalisation section.
BEGIN TRANSACTION;
ALTER TABLE gazetteer ADD COLUMN nr_region_id INTEGER;
CREATE TABLE nr_region (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR NOT NULL);
CREATE INDEX ix_tmp_gaz ON gazetteer(nr_region);
INSERT INTO nr_region SELECT null, nr_region FROM gazetteer GROUP BY nr_region;
CREATE INDEX ix_tmp_lookup ON nr_region(name);
UPDATE gazetteer SET nr_region_id = (SELECT id FROM nr_region WHERE gazetteer.nr_region = nr_region.name);
DROP INDEX ix_tmp_gaz;
ALTER TABLE gazetteer DROP COLUMN nr_region;
DROP INDEX ix_tmp_lookup;
COMMIT;

BEGIN TRANSACTION;
ALTER TABLE gazetteer ADD COLUMN country_id INTEGER;
CREATE TABLE country (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR NOT NULL);
CREATE INDEX ix_tmp_gaz ON gazetteer(country);
INSERT INTO country SELECT null, country FROM gazetteer GROUP BY country;
CREATE INDEX ix_tmp_lookup ON country(name);
UPDATE gazetteer SET country_id = (SELECT id FROM country WHERE gazetteer.country = country.name);
DROP INDEX ix_tmp_gaz;
ALTER TABLE gazetteer DROP COLUMN country;
DROP INDEX ix_tmp_lookup;
COMMIT;

BEGIN TRANSACTION;
ALTER TABLE gazetteer ADD COLUMN county_district_id INTEGER;
CREATE TABLE county_district (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR NOT NULL);
CREATE INDEX ix_tmp_gaz ON gazetteer(county_district);
INSERT INTO county_district SELECT null, county_district FROM gazetteer GROUP BY county_district;
CREATE INDEX ix_tmp_lookup ON county_district(name);
UPDATE gazetteer SET county_district_id = (SELECT id FROM county_district WHERE gazetteer.county_district = county_district.name);
DROP INDEX ix_tmp_gaz;
ALTER TABLE gazetteer DROP COLUMN county_district;
DROP INDEX ix_tmp_lookup;
COMMIT;

BEGIN TRANSACTION;
ALTER TABLE gazetteer ADD COLUMN admin_area_id INTEGER;
CREATE TABLE admin_area (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR NOT NULL);
CREATE INDEX ix_tmp_gaz ON gazetteer(admin_area);
INSERT INTO admin_area SELECT null, admin_area FROM gazetteer GROUP BY admin_area;
CREATE INDEX ix_tmp_lookup ON admin_area(name);
UPDATE gazetteer SET admin_area_id = (SELECT id FROM admin_area WHERE gazetteer.admin_area = admin_area.name);
DROP INDEX ix_tmp_gaz;
ALTER TABLE gazetteer DROP COLUMN admin_area;
DROP INDEX ix_tmp_lookup;
COMMIT;

BEGIN TRANSACTION;
ALTER TABLE gazetteer ADD COLUMN place_name_id INTEGER;
CREATE TABLE place_name (id INTEGER PRIMARY KEY AUTOINCREMENT, name VARCHAR NOT NULL);
CREATE INDEX ix_tmp_gaz ON gazetteer(place_name);
INSERT INTO place_name SELECT null, place_name FROM gazetteer GROUP BY place_name;
CREATE INDEX ix_tmp_lookup ON place_name(name);
UPDATE gazetteer SET place_name_id = (SELECT id FROM place_name WHERE gazetteer.place_name = place_name.name);
DROP INDEX ix_tmp_gaz;
ALTER TABLE gazetteer DROP COLUMN place_name;
DROP INDEX ix_tmp_lookup;
COMMIT;


-- Create a view to query the gazetteer and return denormalised locations.
CREATE VIEW gazetteer_summary AS
SELECT
	g.elr, g.total_yards, g.mileage, r.name AS nr_region, c.name AS country, aa.name AS admin_area, cd.name AS county_district, pn.name AS place_name, CAST(g.distance_m AS INTEGER) AS distance_m
FROM
	gazetteer g
JOIN
	nr_region r ON g.nr_region_id = r.id
JOIN
	country c ON g.country_id = c.id
JOIN
	admin_area aa ON g.admin_area_id = aa.id
JOIN
	county_district cd ON g.county_district_id = cd.id
JOIN
	place_name pn ON g.place_name_id = pn.id;


-- Helper tables - delimiter of ";" used to avoid CSV data transfer ambiguity.
CREATE TABLE elr_by_admin_area AS
SELECT country, admin_area, GROUP_CONCAT(elr, ";") AS elrs
FROM (
	SELECT DISTINCT country, admin_area, elr
	FROM gazetteer_summary
)
GROUP BY country, admin_area;

CREATE INDEX ix_admin_area ON elr_by_admin_area(admin_area);


CREATE TABLE elr_by_county_district_place_name AS
SELECT county_district, place_name, GROUP_CONCAT(elr, ";") AS elrs
FROM (
	SELECT DISTINCT county_district, place_name, elr
	FROM gazetteer_summary
)
GROUP BY county_district, place_name;

CREATE INDEX ix_place_name ON elr_by_county_district_place_name(place_name);


VACUUM;
ANALYZE;
