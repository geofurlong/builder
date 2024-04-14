// External commands executed by the builder functions.

package main

import (
	"bytes"
	"geofurlong/pkg/geocode"
	"log"
	"os"
	"os/exec"
)

// runPython executes the external Python script, logging its output.
func runPython(cfg GeofurlongConfig, scriptFn string, params string) {
	originalDir, err := os.Getwd()
	geocode.Check(err)

	err = os.Chdir(cfg["scripts_dir"])
	geocode.Check(err)

	// Ensure the working directory is reset after the function completes.
	defer func() {
		err = os.Chdir(originalDir)
		geocode.Check(err)
	}()

	var cmd *exec.Cmd

	if params == "" {
		cmd = exec.Command("python3", scriptFn)
	} else {
		cmd = exec.Command("python3", scriptFn, params)
	}

	output, err := cmd.CombinedOutput()
	log.Println(string(output))
	geocode.Check(err)
}

// deleteFile deletes the specified file if it exists.
func deleteFile(filename string) {
	if _, err := os.Stat(filename); err == nil {
		err = os.Remove(filename)
		geocode.Check(err)
	} else if os.IsNotExist(err) {
	} else {
		geocode.Check(err)
	}
}

// runSQLiteCommand executes an SQL script on an SQLite database.
func runSQLiteCommand(db string, inputScript string) {
	cmd := exec.Command("sqlite3", "-bail", db)
	cmd.Stdin = bytes.NewBuffer([]byte(inputScript))
	err := cmd.Run()
	geocode.Check(err)
}
