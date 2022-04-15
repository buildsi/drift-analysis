package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type config struct {
	General        general
	Authentication authentication
	Database       database
	Relay          relay
}

type general struct {
	Version string
	Port    int
}
type authentication struct {
	User     string
	Password string
}

type database struct {
	Location string
}

type relay struct {
	Location         string
	User             string
	Password         string
	TempFileLocation string
}

func defaultConf(path string) config {
	result := config{
		General: general{
			Version: "0.0.1",
			Port:    8080,
		},
		Authentication: authentication{
			User:     "",
			Password: "",
		},
		Database: database{
			Location: filepath.Join(path, "drift-server.db"),
		},
		Relay: relay{
			Location:         "",
			User:             "",
			Password:         "",
			TempFileLocation: filepath.Join(path, "tmp"),
		},
	}
	return result
}

var (
	// Global config with the database.
	Config config
)

// initialize the app config system. If a config doesn't exist, create one.
// If the config is out of date read the current config and rebuild with new fields.
func init() {

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	Config = defaultConf(pwd)
}
