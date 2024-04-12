# Convert the Network Rail (ELR, Milepost, and Region) and Ordnance Survey (Populated Place and Administrative Area) datasets to SQLite format.
# Also converts the manually-maintained GeoFurlong ELR Excel file to a CSV file.

import logging
import subprocess
import sqlite3
import pandas as pd
import geopandas as gpd
import file_ops
import config


def miles_yards_to_total_yards(miles_yards: float) -> int:
    """Convert a floating-point decimal mileage (of form mmm.yyyy) to a total yards integer."""
    YARDS_IN_MILE = 1_760
    miles = int(miles_yards)
    return int(round(YARDS_IN_MILE * miles + 10_000 * (miles_yards - miles), 0))


def km_to_total_yards(kilometres: float) -> int:
    """Converts a floating-point kilometreage to a nearest total yards integer."""
    KM_IN_YARD = 0.0009144
    return int(kilometres / KM_IN_YARD + 0.5)


def to_total_yards(linear_unit: str, value: float) -> int:
    """Returns an equivalent total yardage integer based on the input value and linear unit."""
    if linear_unit == "M":
        return miles_yards_to_total_yards(value)
    return km_to_total_yards(value)


def convert_centrelines() -> None:
    """Convert the NR ELR Centre-line Shapefile to a SQLite database."""
    logging.info("Import NR ELR Centre-lines")
    file_ops.delete_file(CONFIG["cl_db"])

    tmp_cl_db = CONFIG["cl_db"].replace(".sqlite", "_tmp.sqlite")
    file_ops.delete_file(tmp_cl_db)

    # `CastToXY(geometry)` converts 3D linestrings to 2D.
    sql = f'SELECT elr, l_system, l_m_from, l_m_to, shape_len AS shape_length_m, CastToXY(geometry) AS geometry FROM NetworkReferenceLines {CONFIG["skip_elr_sql"]}'

    cmd_ogr_convert = [
        "ogr2ogr",
        "-f",
        "SQLite",
        tmp_cl_db,
        CONFIG["cl_shp"],
        "-nln",
        "cl",
        "-dialect",
        "sqlite",
        "-sql",
        sql,
    ]

    if subprocess.run(cmd_ogr_convert, check=True, encoding="utf-8").returncode != 0:
        raise Exception("Failed to import ELR centre-lines")

    # Some historic NR ELR centre-line datasets contain 100% duplicates, so these need removed.
    # In the event of this reoccurring, uncomment the lines below.

    # cmd_dedupe = ["ogrinfo", "-q", "-sql",
    #               "DELETE FROM cl WHERE rowid NOT IN (SELECT MIN(rowid) FROM cl GROUP BY elr)",
    #               tmp_cl_db]
    # if subprocess.run(cmd_dedupe, check=True, encoding="utf-8").returncode != 0:
    #     raise Exception("Failed to deduplicate ELR centre-lines")

    centrelines = gpd.read_file(tmp_cl_db, engine="pyogrio")

    # The ELR centre-line Shapefile linestring is orientated in mileage direction in all cases.
    # Convert ELR centre-line extents from decimal mile/yards to total yards.
    # Extents of metric ELRs are already provided as miles/yards in the NR FoI source.
    centrelines["total_yards_from"] = centrelines["l_m_from"].apply(miles_yards_to_total_yards)
    centrelines["total_yards_to"] = centrelines["l_m_to"].apply(miles_yards_to_total_yards)

    # Remove obsolete columns.
    centrelines.drop(["l_m_from", "l_m_to"], axis=1, inplace=True)

    # Export to SQLite format.
    centrelines.to_file(CONFIG["cl_db"], driver="SQLite", layer="cl", engine="pyogrio")  # type: ignore

    file_ops.delete_file(tmp_cl_db)


def convert_mileposts() -> None:
    """Convert the NR Milepost Shapefile to a SQLite database."""
    logging.info("Import NR Mileposts")
    file_ops.delete_file(CONFIG["mp_db"])

    sql = f'SELECT elr, m_system, waymark_va, geometry FROM NetworkWaymarks {CONFIG["skip_elr_sql"]}'

    cmd_ogr_convert = [
        "ogr2ogr",
        "-f",
        "SQLite",
        CONFIG["mp_db"],
        CONFIG["mp_shp"],
        "-nln",
        "mp",
        "-dialect",
        "sqlite",
        "-sql",
        sql,
    ]

    if subprocess.run(cmd_ogr_convert, check=True, encoding="utf-8").returncode != 0:
        raise Exception("Failed to import Mileposts")

    mileposts = gpd.read_file(CONFIG["mp_db"], engine="pyogrio")

    # The Milepost reporting values are in miles/yards or kilometres (unlike the ELR reporting values which are always miles/yards).
    # Convert Milepost `mileages` from source value (miles or kilometres) to total yards.
    # Column name specified as used later for calibration between mileposts, i.e. "from" and "to".
    mileposts["total_yards_from"] = mileposts.apply(lambda x: to_total_yards(x.m_system, x.waymark_va), axis=1)

    # Remove obsolete columns.
    mileposts.drop(["m_system", "waymark_va"], axis=1, inplace=True)

    # Export to SQLite format.
    mileposts.to_file(CONFIG["mp_db"], driver="SQLite", layer="mp", engine="pyogrio")  # type: ignore

    # Create SQLite composite index on Milepost table to reduce calibration time.
    conn = sqlite3.connect(CONFIG["mp_db"])
    cur = conn.cursor()
    cur.execute("CREATE UNIQUE INDEX ix_elr_total_yards ON mp (elr, total_yards_from)")
    conn.commit()
    conn.close()


def convert_nr_regions() -> None:
    """Convert the NR Region Shapefile to a SQLite database."""
    logging.info("Import NR Regions")
    nr_regions = gpd.read_file(CONFIG["nr_region_shp"], encoding="utf8", engine="pyogrio")

    # Remove unwanted columns.
    nr_regions.drop(
        [
            "OBJECTID",
            "SHAPE_Leng",
            "SHAPE_Area",
        ],
        axis=1,
        inplace=True,
    )

    # Rename NR Region column to avoid a naming clash with OS Region column when spatially joining.
    nr_regions.rename(
        columns={
            "REGION_NAM": "nr_region",
        },
        inplace=True,
    )

    # Export to SQLite format.
    nr_regions.to_file(CONFIG["nr_region_db"], driver="SQLite", layer="nr_region", engine="pyogrio")  # type: ignore


def convert_os_places() -> None:
    """Convert the OS Populated Place Shapefile to a SQLite database."""
    logging.info("Import OS Populated Places")
    os_places = gpd.read_file(CONFIG["os_place_shp"], encoding="utf8", engine="pyogrio")

    # Remove unwanted columns. Country or Region columns are not wanted as these are established
    # from the point in polygon spatial join via the Administrative Area dataset.
    os_places.drop(
        ["LOCAL_TYPE", "POSTCODE_D", "REGION", "COUNTRY"],
        axis=1,
        inplace=True,
    )

    os_places.rename(
        columns={
            "NAME1": "place_name",
            "NAME1_LANG": "name_lang",
            "NAME2": "place_name_alt",
            "NAME2_LANG": "name_alt_lang",
            "DISTRICT_B": "district",
            "COUNTY_UNI": "county_unitary",
        },
        inplace=True,
    )

    # Set OS Populated Place name to ENG variation, i.e. where Welsh or Gaelic are
    # listed as the primary Populated Place name.
    os_places["place_name"] = os_places.apply(
        lambda x: (x["place_name_alt"] if x["name_alt_lang"] == "eng" else x["place_name"]),
        axis=1,
    )

    # Remove obsolete columns.
    os_places.drop(["name_lang", "place_name_alt", "name_alt_lang"], axis=1, inplace=True)

    # Export to SQLite format.
    os_places.to_file(CONFIG["os_place_db"], driver="SQLite", layer="os_place", engine="pyogrio")  # type: ignore


def convert_os_admin_areas() -> None:
    """Convert the OS Administrative Area Shapefile to a SQLite database."""
    logging.info("Import OS Administrative Areas")
    os_admin_areas = gpd.read_file(CONFIG["os_admin_area_shp"], encoding="utf8", engine="pyogrio")

    # Remove unwanted columns.
    os_admin_areas.drop(["ID_0", "ID_1", "ID_2"], axis=1, inplace=True)

    os_admin_areas.rename(
        columns={
            "NAME_1": "country",
            "NAME_2": "admin_area",
        },
        inplace=True,
    )

    # Export to SQLite format.
    os_admin_areas.to_file(
        CONFIG["os_admin_area_db"],
        driver="SQLite",
        layer="os_admin_area",
        engine="pyogrio",
    )  # type: ignore


def convert_elrs() -> None:
    """Convert the GeoFurlong manually-maintained ELR Excel file to a CSV file."""
    logging.info("Import ELRs (attributes)")
    elr_xlsx = pd.ExcelFile(CONFIG["elr_xlsx"])
    df_elr = pd.read_excel(elr_xlsx, "master")

    df_elr = df_elr[
        [
            "elr",
            "route",
            "section",
            "remarks",
            "quail_book",
            "grouping",
            "neighbours",
        ]
    ]

    # Columns `quail_book`, `grouping`, and `neighbours` are delimited with ";" in the source XLSX file.
    elr_db = CONFIG["elr_csv"]
    file_ops.delete_file(elr_db)

    # Export to CSV format.
    df_elr.to_csv(elr_db, index=False)


if __name__ == "__main__":
    logging.basicConfig(
        level=logging.INFO,
        handlers=[
            logging.FileHandler("conversion.log"),
            logging.StreamHandler(),
        ],
        format="%(levelname)s,%(asctime)s,%(message)s",
    )

    CONFIG = config.read()

    file_ops.check_files_exist(
        (
            CONFIG["cl_shp"],
            CONFIG["mp_shp"],
            CONFIG["nr_region_shp"],
            CONFIG["os_place_shp"],
            CONFIG["os_admin_area_shp"],
            CONFIG["elr_xlsx"],
        )
    )

    logging.info(f"Data conversion started (version {CONFIG['version']})")
    convert_centrelines()
    convert_mileposts()
    convert_elrs()
    convert_nr_regions()
    convert_os_places()
    convert_os_admin_areas()
    logging.info("Data conversion ended")
