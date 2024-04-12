# Read the project configuration file.

import os
import yaml


def read() -> dict:
    """Read the project configuration file, returning as a dictionary."""
    env_var = "GEOFURLONG_ROOT"
    config_file = "geofurlong_config.yaml"

    root_dir = os.environ.get(env_var, "")
    if root_dir == "":
        raise ValueError(f"The required environment variable {env_var} is not set.")

    with open(f"{root_dir}/{config_file}", "r") as config_file:
        config = yaml.safe_load(config_file)

    # Replace root_dir in the settings with the environment variable value.
    for key, value in config["settings"].items():
        if isinstance(value, str):
            config["settings"][key] = value.replace("${root_dir}", root_dir)

    return config["settings"]
