package config

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

// Configuration for drift server
type Config struct {
	Debug        bool   `envconfig:"DEBUG" default:"false"`
	Port         string `envconfig:"PORT" default:"8080"`
	DatabaseDSN  string `envconfig:"DATABASE_DSN" default:""`
	DatabaseType string `envconfig:"DATABASE_TYPE" default:"sqlite"`
	Username     string `envconfig:"USER" default:""`
	Password     string `envconfig:"PASS" default:""`
}

// NewConfig inits a new config
func NewConfig() Config {
	s := Config{}
	if err := envconfig.Process("DRIFTSERVER", &s); err != nil {
		log.Fatal(err)
	}
	return s
}
