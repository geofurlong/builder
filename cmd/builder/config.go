// Read the project configuration file.

package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type GeofurlongConfig map[string]string

type Config struct {
	Settings GeofurlongConfig `yaml:"settings"`
}

// Read the project configuration file, located at directory pointed to by the GEOFURLONG_ROOT environment variable.
func readConfig() (GeofurlongConfig, error) {
	envVar := "GEOFURLONG_ROOT"
	configFile := "geofurlong_config.yaml"

	rootDir, exists := os.LookupEnv(envVar)
	if !exists {
		return nil, fmt.Errorf("the required environment variable %s is not set", envVar)
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/%s", rootDir, configFile))
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	for key, value := range config.Settings {
		config.Settings[key] = strings.Replace(value, "${root_dir}", rootDir, -1)
	}

	return config.Settings, nil
}
