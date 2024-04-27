# Builds a gazetteer, combining Network Rail and Ordnance Survey datasets for a defined yardage resolution.

from os.path import abspath
import sys
import subprocess
import logging
import pandas as pd
import geopandas as gpd
import file_ops
import config


def build_csv(rail_points_fn: str, out_csv_fn: str) -> None:
    """Build and save a gazetteer for a defined yardage resolution."""
    logging.info("Loading precomputed railway points")
    rail_pts = pd.read_csv(rail_points_fn, encoding="utf8")

    rail_pts = gpd.GeoDataFrame(
        rail_pts,
        geometry=gpd.points_from_xy(rail_pts.easting, rail_pts.northing),
        crs="EPSG:27700",
    )  # type: ignore

    logging.info("Loading NR Regions")
    nr_regions = gpd.read_file(CONFIG["nr_region_db"], encoding="utf8", engine="pyogrio")

    logging.info("NR Regions spatial join")
    gdf_rail_nr = rail_pts.sjoin(nr_regions, how="left", predicate="within")

    # Spatial join with NR Regions may create duplicate rows for a single railway point,
    # so deduplicate, noting that the deleted row(s) will be arbitrary and is a relatively slow operation.
    logging.info("Removing duplicate spatial join rows")
    gdf_rail_nr = gdf_rail_nr.groupby(level=0).first()
    gdf_rail_nr = gdf_rail_nr.drop(["index_right"], axis=1)

    logging.info("Loading OS Populated Places")
    places = gpd.read_file(CONFIG["os_place_db"], encoding="utf8", engine="pyogrio")

    logging.info("Nearest Populated Places spatial join")
    gdf_rail_nr.crs = "EPSG:27700"
    gdf_rail_nr_os = gdf_rail_nr.sjoin_nearest(places, distance_col="distance_m")
    # Round distance from railway point to nearest Populated Place down to nearest metre.
    gdf_rail_nr_os["distance_m"] = gdf_rail_nr_os["distance_m"].astype("i8")
    gdf_rail_nr_os = gdf_rail_nr_os.drop(columns=["index_right"])

    logging.info("Loading OS Administrative Area boundaries")
    admin_areas = gpd.read_file(CONFIG["os_admin_area_db"], encoding="utf8", engine="pyogrio")

    logging.info("Administrative Area spatial join")
    gdf_rail_nr_os = gdf_rail_nr_os.sjoin(admin_areas, how="left", predicate="within")
    master = pd.DataFrame(gdf_rail_nr_os.drop(columns=["geometry", "index_right"]))

    # Final output contains railway, NR Region, Populated Place, and Administrative Area.
    logging.info("Sorting gazetteer by ELR and mileage")
    master = master.sort_values(["elr", "total_yards"])

    # Export gazetteer to CSV.
    logging.info("Saving gazetteer as CSV")
    master.to_csv(out_csv_fn, index=False, encoding="utf8")


def generate(resolution: int) -> None:
    """Create the gazetteer CSV and database files for the yardage resolution."""
    logging.info(f"Creating Gazetteer at {resolution}y resolution")
    rail_points = f'{CONFIG["precompute_dir"]}/geofurlong_precomputed_{resolution:04}y.csv'
    gazetteer_csv = f'{CONFIG["gazetteer_dir"]}/geofurlong_gazetteer_{resolution:04}y_raw.csv'
    gazetteer_db = f'{CONFIG["gazetteer_dir"]}/geofurlong_gazetteer_{resolution:04}y.sqlite'

    file_ops.delete_file(gazetteer_csv)
    file_ops.delete_file(gazetteer_db)

    build_csv(rail_points, gazetteer_csv)

    # Create an optimised and normalised SQLite gazetteer database.
    logging.info("Optimising gazetteer tables")
    with open("gazetteer_create.sql", "rb") as sql_script:
        stdin = sql_script.read()

    # Replace the default resolution in the SQL script with the current resolution.
    stdin = stdin.decode().replace("gazetteer_0022y_raw.csv", "gazetteer_{:04}y_raw.csv".format(resolution)).encode()

    if (
        subprocess.run(
            ["sqlite3", "-bail", gazetteer_db],
            input=stdin,
            cwd=abspath(CONFIG["gazetteer_dir"]),
        ).returncode
        != 0
    ):
        raise Exception("Failed to create gazetteer database")


if __name__ == "__main__":
    logging.basicConfig(
        level=logging.INFO,
        handlers=[
            logging.FileHandler("gazetteer.log"),
            logging.StreamHandler(),
        ],
        format="%(levelname)s,%(asctime)s,%(message)s",
    )

    CONFIG = config.read()

    if len(sys.argv) != 2:
        print("Usage: python3 gazetteer.py <resolution>")
        sys.exit(1)

    resolution = int(sys.argv[1])
    generate(resolution)
