# GeoFurlong

Geospatial resources for mainline (Network Rail) railways in Britain.

## Description

This repository contains the source-code for the builder process, which transforms multiple geospatial data sources into optimised datasets for client application use. These output datasets are stored on [Dropbox](https://www.dropbox.com/scl/fo/1tzbnx4zaz61nwslqobn9/AOcmAULa79IuAjLNIvZThRQ?rlkey=fx8e3gxwhtgmlq2v0akv827ty&dl=0) and summarised in the [Data Catalogue](#data-catalogue) section below.

The output datasets provide ready-made geocoded data for railway locations, formed of an [Engineer's Line Reference](https://en.wikipedia.org/wiki/Engineer%27s_Line_Reference) (ELR) and associated mileage (or kilometreage) on that line. These precomputed geographic positions can be reverse-geocoded to establish the railway ELR and mileage relative to other features using standard nearest neighbour or point-in-polygon spatial analysis using GIS tools or software libraries.

Sample applications will be shared via separate repositories which utilise the geocoded datasets, dynamic geocoding (at several thousand points per second), and reverse geocoding for mobile applications.

🌐 [geofurlong.com](https://www.geofurlong.com/) is built using the output of this project, with interactive mapping of routes and junction areas, tabulated output detailing railway attributes, and spatial relationship with populated places and government administrative boundaries.

## Preface

GeoFurlong uses a _total yards_ value to define the reported linear distance along the railway, commonly referred to as a _mileage_, even for kilometre-based ELRs. This is a signed whole number which avoids the pitfalls associated with attempting to store mileage in the many permutations of decimal miles, fractional miles, or text format. These pitfalls are amplified when dealing with negative mileages. The _total yards_ unit is unambiguous, and efficient for sorting, filtering, and storage within systems.

GeoFurlong is opinionated and consistent in its textual presentation of mileages. For example, a mileage of `86 miles`, `7 yards` is presented as `86M 0007y`.

Recording of geographic position is [precise](https://en.wikipedia.org/wiki/Accuracy_and_precision) to one decimal place for Ordnance Survey Easting / Northing (i.e. 100 mm) and six decimal places for Longitude / Latitude (approximately 110 mm in Britain).

Linear accuracy is defined as the geographic measured distance versus the reported distance, both in metres. For example, if the measured distance between neighbouring quarter mileposts along an ELR centre-line was `403.836 metres`, the accuracy would be calculated as `+1.5 metres` (as a quarter mile being 440 yards, or `402.336 metres`). This is an example of what is commonly referred to as a _long quarter mile_. The linear accuracy, computed to maximum available decimal places, is used to produce the linear calibration statistics per ELR; it is subsequently truncated to a whole number for presentation in other data sets.

The computed geographic position for a defined ELR and mileage may not be accurate in all instances. In a number of locations, the position may be incorrect by a significant linear distance, particularly on closed or partially-closed lines. The manually-maintained _ELR_ dataset (via the `Remarks` column) identifies ELRs which exhibit potentially poor accuracy.

The build process computes the estimated linear position for a given mileage on an ELR by calibrating against mileposts on that ELR. For each ELR, calibration in undertaken using the virtual centre-line geometry, reported start and finish mileages, combined with the milepost position and value. The computed geographic distance along the segment between mileposts are compared against the reported mileages for the mileposts and recorded in a detailed calibration statistics database. This calibration process allows an estimation of the linear accuracy to be provided when geocoding from ELR and Mileage to geographic position.

Noting the linear calibration process described above, inaccuracies in estimating geographic position of a mileage on an ELR can result as a consequence of individual or combined factors which are out with the control of this project, including:

- **ELR**: Incorrect geometry; incorrect start / finish reported mileage; remodelled track layout.
- **Milepost**: Incorrect position; incorrect identifier; not recorded / physically missing.

>⚠️ Given these credible risks to positional accuracy, the GeoFurlong datasets must not be used for safety-critical decisions, nor relied upon as the primary geographic data source in production environments.

## Published Data Files

### ELR Schema

| Column | Description | Unit / Type | Sample |
| ---- | ---- | ---- | --- |
|elr|ELR|text|WCM1|
|l_system|Linear Reporting Unit|text (M or K)|M|
|shape_length_m|Geographic Length|metres|135756.658175|
|total_yards_from|Mileage From|total yards (whole number)|-216|
|total_yards_to|Mileage To|total yards (whole number)|148224|
|route|Route|text|West Coast Main Line (WCML)|
|section|Section|text (optional)|Carlisle to Law Jn|
|remarks|Remarks|text (optional)||
|quail_book|TrackMap Book|text (`;` separated)|1;4;2|
|grouping|Grouping|text (`;` separated)|LEC1;LEC2;LEC3; ...|
|neighbours|Neighbours|text (`;` separated)|CGJ7;CSP;ECA1;ETC; ...|

### Precomputed Schema

| Column | Description | Unit / Type | Sample |
| ---- | ---- | ---- | --- |
|elr|ELR|text|ECM1|
|total_yards|Mileage|total yards (whole number)|5654|
|mileage|Mileage|text|3M 0374y|
|easting|OS Easting|metres (1 decimal place)|531412.3|
|northing|OS Northing|metres (1 decimal place)|187912.1|
|longitude|Longitude|degrees (6 decimal places)|-0.105067|
|latitude|Latitude|degrees (6 decimal places)|51.574767|
|osgr|OS Grid Reference|text|TQ3141287912|
|accuracy|Linear Accuracy|metres (whole number)|-2|

### Gazetteer Schema

| Column | Description | Unit / Type | Sample |
| ---- | ---- | ---- | --- |
|elr|ELR|text|ECM1|
|total_yards|Mileage|total yards (whole number)|5654|
|mileage|Mileage|text|3M 0374y|
|easting|OS Easting|metres (1 decimal place)|531412.3|
|northing|OS Northing|metres (1 decimal place)|187912.1|
|longitude|Longitude|degrees (6 decimal places)|-0.105067|
|latitude|Latitude|degrees (6 decimal places)|51.574767|
|osgr|OS Grid Reference|text|TQ3141287912|
|accuracy|Linear Accuracy|metres (whole number)|-2|
|nr_region|Network Rail Region|text|Eastern|
|place_name|Nearest Populated Place|text|Stroud Green|
|district|Nearest Populated Place's District|text|Haringey|
|county_unitary|Nearest Populated Place's County|text|Greater London|
|distance_m|Distance to nearest Populated Place|metres (whole number)|540|
|country|Country|text|England|
|admin_area|Administrative Area|text|Haringey|

### Aggregated Gazetteer

At the maximum resolution of 22 yards, the gazetteer table consists of over 850,000 entries. An alternative method of establishing the geographic context of the railway positions is made available by grouping the following attributes into a mileage range: Network Rail Region, Government Administrative Area, and nearest Populated Place (and its corresponding County / District) in the following tables, each which have a significantly reduced number of entries:

- `geofurlong_gazetteer_by_nr_region.csv`
- `geofurlong_gazetteer_by_country_admin_area.csv`
- `geoofurlong_gazetteer_by_nearest_place.csv`

### Data Catalogue

| Filename | Description | Record Count | File Size |
| --- | --- | ---: | --: |
| geofurlong_elr.csv | ELR master list | 1,589 | 187.2 KB |
| geofurlong_elr_metric.csv | ELR (metric) | 19 | 92.0 B |
| geofurlong_precomputed_0022y.csv | Geographic positions at 22 yard intervals | 884,780 | 62.4 MB |
| geofurlong_precomputed_0110y.csv | Geographic positions at 110 yard intervals | 179,632 | 12.7 MB |
| geofurlong_precomputed_0220y.csv | Geographic positions at 220 yard intervals | 91,536 | 6.4 MB |
| geofurlong_precomputed_0440y.csv | Geographic positions at 440 yard intervals | 47,515 | 3.3 MB |
| geofurlong_precomputed_1760y.csv | Geographic positions at 1760 yard interval | 14,469 | 1.0 MB |
| geofurlong_precomputed_8800y.csv | Geographic positions at 8800 yard intervals | 5,722 | 407.3 KB |
| geofurlong_gazetteer_0022y.csv | Gazetteer at 22 yard intervals | 884,780 | 73.3 MB |
| geofurlong_gazetteer_0110y.csv | Gazetteer at 110 yard intervals | 179,632 | 14.9 MB |
| geofurlong_gazetteer_0220y.csv | Gazetteer at 220 yard intervals | 91,536 | 7.6 MB |
| geofurlong_gazetteer_0440y.csv | Gazetteer at 440 yard intervals | 47,515 | 3.9 MB |
| geofurlong_gazetteer_1760y.csv | Gazetteer at 1760 yard intervals | 14,469 | 1.2 MB |
| geofurlong_gazetteer_8800y.csv | Gazetteer at 8800 yard intervals | 5,722 | 478.9 KB |
| geofurlong_gazetteer_by_nr_region.csv | Gazetteer by Network Rail region | 1,640 | 76.1 KB |
| geofurlong_gazetteer_by_country_admin_area.csv | Gazetteer by Country and Administrative Area | 2,399 | 129.3 KB |
| geofurlong_gazetteer_by_nearest_place.csv | Gazetteer by Nearest Populated Place (and County / District) | 14,507 | 1.1 MB |
| geofurlong_gazetteer_aggregated.csv | Gazetteer (aggregated) | 18,544 | 1.6 MB |
| geofurlong_elr_by_country_admin_area.csv | ELRs within each Country and Administrative Area | 166 | 13.3 KB |
| geofurlong_elr_by_nearest_place.csv | ELRs with Nearest Populated Place (and County / District) | 10,416 | 361.3 KB |
| geofurlong_calibration_simplified.csv | Linear calibration (simplified) | 44,257 | 2.3 MB |
| geofurlong_calibration_full.csv | Linear calibration (full) | 44,257 | 5.4 MB |
| geofurlong_calibration_statistics.csv | Linear calibration statistics | 1,589 | 379.3 KB |

### To Do

Supplementary files will be distributed to support sample client and third-party applications:

- Production database (SQLite format).
- Geospatial databases at varying yardage intervals (GeoPackage format) for `precomputed` and `gazetteer` tables.
- Improved linear positioning accuracy could be obtained by utilising more recent surveyed position of mileposts. Milepost positions are regularly surveyed as a matter of course during topographic survey on the network. The process of collating and overriding milepost positions to improve the calibration accuracy is currently not within the scope of this project.

## Builder Technical Details

### Software Stack

GeoFurlong is primarily developed in the [Go](https://go.dev/) programming language, delegating certain input and output geospatial file operations to [Python](https://www.python.org/) scripts, utilising well-proven libraries. Input data files are in ESRI [Shapefile](https://en.wikipedia.org/wiki/Shapefile) format, intermediate files as comma-separated value ([CSV](https://en.wikipedia.org/wiki/Comma-separated_values)) format, and output files predominantly as [SQLite](https://en.wikipedia.org/wiki/SQLite) databases (with geometry columns stored in well-known binary ([WKB](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry)) format).

### Process

- Manual validation / preparation (see below).
- Conversion of source geospatial to optimised SQLite format: ELRs, Mileposts, Network Rail Regions, Ordnance Survey Administrative Areas, and Ordnance Survey Populated Places.
- Calibrate mileposts along each ELR centre-line geometry to maximise linear positional accuracy.
- Build optimised production database of ELRs and associate linear calibration.
- Precompute geographic positions for all ELRs at multiple yardage intervals: 22, 110, 220, 440, 1760 (one mile), and 8800 (5 miles).
- Build a gazetteer of railway positions combining Network Railway Region, Ordnance Survey Administrative Area and Populated Place datasets at multiple yardage intervals: 22, 110, 220, 440, 1760 (one mile), and 8800 (5 miles).
- Build an aggregated gazetteer, based on 22 yard intervals.

### Source Data Preparation

The Scottish Region geometry from the Network Rail data source has been identified as being invalid due to it containing a self-intersecting ring. This has been manually corrected prior to the data import phase using [QGIS](https://qgis.org/en/site/).

Several populated places are present which share a common geographic position within the Ordnance Survey data source. These have been assessed manually, then removed prior to the data import phase, as noted in the table below.

| Area                          | Deleted Place     | Retained Place       |
|:------------------------------|:------------------|:---------------------|
| Abertawe - Swansea            | Mount Pleasant    | Clydach              |
| Cornwall                      | Toldish           | Indian Queens        |
| County Durham                 | Catchgate         | Annfield Plain       |
| Devon                         | Bishop's Clyst    | Clyst St Mary        |
| Dorset                        | Dudsbury          | West Parley          |
| Dorset                        | Pidney            | Hazelbury Bryan      |
| Dumfries and Galloway         | Minnigaff         | Newton Stewart       |
| Fife                          | Town Centre       | Glenrothes           |
| Gloucestershire               | South Woodchester | Woodchester          |
| Kent                          | Boughton Street   | Boughton Under Blean |
| Moray                         | Old Keith         | Keith                |
| North Lanarkshire             | Garnqueen         | Glenboig             |
| North Lanarkshire             | Wester Auchinloch | Auchinloch           |
| Pen-y-bont ar Ogwr - Bridgend | Evanstown         | Gilfach Goch         |
| Somerset                      | Highbury          | Coleford             |
| Swindon                       | North Wroughton   | Wroughton            |
| Yorkshire and the Humber      | Westfield         | Brampton             |

### Directory Structure

|Directory|Contents|
|----|----|
|`(root)`|Configuration file|
|`cmd/builder`|Master builder program|
|`data/cache`|Serialised cache file|
|`data/gazetteer`|Railway points combined with location data at multiple intervals|
|`data/import`|ELR attributes file|
|`data/import/foi_nr`|Network Rail import files|
|`data/import/foi_os`|Ordnance Survey import files|
|`data/precomputed`|Railway geographic locations precomputed at multiple intervals|
|`data/production`|Database containing ELR centre-lines and calibration data|
|`data/staging`|Intermediate data files used during build process|
|`pkg/geocode`|Go support library files|
|`scripts`|Python and SQL support scripts|

### Data Usage

For most applications, end-users will likely utilise the files contained in the `data/precomputed` or `data/gazetteer` directories, as these contain pre-computed geographic information for regular points (at multiple resolutions) along each ELR. These ready-made tabular files provide simple lookup access to geographic positions for ELRs and mileages without the need for any complex computation.

In addition to these files, developers may use the database in `data/production`, combined with the Go library files in `pkg/geocode` for custom applications to compute the geographic position of an ELR and mileage combination dynamically. This library exposes function to establish a `point` for a single mileage or `substring` for a mileage range. Client libraries for other programming languages are in progress to integrate with the database in `data/production`.

### Definitions

|Definition|Description|
|---|---|
|ELR|Engineer's Line Reference|
|NR|Network Rail|
|OS|Ordnance Survey|

### Credits

The project's software is built upon a framework of open-source software packages and libraries, utilising portions of geospatial datasets which have been released under permissive licences.

### Data

- [Network Rail](https://www.networkrail.co.uk/) data via a Freedom of Information Request by Peter Hicks - [download](https://files.whatdotheyknow.com/request/update_geospatial_data/FOI202400062/).
- [Ordnance Survey](https://www.ordnancesurvey.co.uk/) data packaged by Alasdair Rae of [Automatic Knowledge Ltd](https://automaticknowledge.co.uk/) - [download](https://automaticknowledge.co.uk/resources/).
- The master list of ELR properties is manually maintained by the author, with properties of each ELR, such as neighbours and grouping, partially automated.

### Software

- [Go](https://go.dev/) programming language.
- [GDAL](https://gdal.org/) vector translator library.
- [SQLite](https://sqlite.org/) database library and tools.
- [orb](https://github.com/paulmach/orb) 2D geometry library for Go.
- [go-sqlite3](https://github.com/mattn/go-sqlite3) SQLite database library for Go.
- [go-proj](https://github.com/twpayne/go-proj) co-ordinate transformation library for Go.
- [Python](https://www.python.org/) programming language.
- [pandas](https://github.com/pandas-dev/pandas) data analysis library for Python.
- [GeoPandas](https://github.com/geopandas/) and [Shapely](https://github.com/shapely/shapely) geospatial data manipulation libraries for Python.
- Reference has been made to a Python [library](https://gitlab.com/jbrobertson/os-grid-reference/) for manipulation of OS Grid References.
- [QGIS](https://qgis.org/en/site/) Geographical Information System.
- [VisiData](https://www.visidata.org/) interactive tabular data analysis tool.

### Disclaimer

The output data is provided as is, with no warranty of any kind, express or implied.

In no event shall the GeoFurlong author be liable for any claim, damages or other liability, whether in an action of contract, tort or otherwise, arising from, out of or in connection with this repository or data output.

### Licence

All data outputs generated by the GeoFurlong tools are released under a permissive [CC BY Creative Commons Licence](https://creativecommons.org/licenses/by/4.0/). This licence allows reusers to distribute, remix, adapt, and build upon the material in any medium or format, so long as attribution is given to the creator. The licence allows for commercial use. however it should be noted that the GeoFurlong data is aggregated from Network Rail and Ordnance Survey data, both released under [Open Government Licences](https://www.nationalarchives.gov.uk/doc/open-government-licence/version/3/), and these contributing licences must be respected.

The project's source code is released under the permissive [MIT Licence](https://opensource.org/licenses/MIT) with a view to benefitting those working in the railway environment and foster further innovation.

### Author

Alan Morrison _CEng MICE Eur Ing FPWI_
