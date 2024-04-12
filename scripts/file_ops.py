# Functions for file operations.

import os

FILE_NOT_FOUND = -1


def delete_file(fn: str) -> None:
    """Deletes an operating system file."""
    if os.path.isfile(fn):
        os.remove(fn)


def file_size(fn: str) -> int:
    """Returns the file size in bytes, or -1 if it doesn't exist."""
    if os.path.isfile(fn):
        return os.stat(fn).st_size
    else:
        return FILE_NOT_FOUND


def check_files_exist(reqd_files: tuple[str, ...]) -> None:
    """Checks for existence of required files, aborting if any are absent."""
    for reqd_file in reqd_files:
        if file_size(reqd_file) == FILE_NOT_FOUND:
            raise Exception(f"Required file not present: {reqd_file}")
