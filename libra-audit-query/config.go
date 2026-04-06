package main

import (
	"fmt"
)

// Config defines all of the service configuration parameters
type Config struct {
	// database configuration
	DbHost     string // database host
	DbPort     int    // database port
	DbName     string // database name
	DbUser     string // database user
	DbPassword string // database password
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	cfg.DbHost, err = ensureSetAndNonEmpty("DB_HOST")
	if err != nil {
		return nil, err
	}
	cfg.DbPort, err = envToInt("DB_PORT")
	if err != nil {
		return nil, err
	}
	cfg.DbName, err = ensureSetAndNonEmpty("DB_NAME")
	if err != nil {
		return nil, err
	}
	cfg.DbUser, err = ensureSetAndNonEmpty("DB_USER")
	if err != nil {
		return nil, err
	}
	cfg.DbPassword, err = ensureSetAndNonEmpty("DB_PASSWORD")
	if err != nil {
		return nil, err
	}

	fmt.Printf("[conf] DbHost          = [%s]\n", cfg.DbHost)
	fmt.Printf("[conf] DbPort          = [%d]\n", cfg.DbPort)
	fmt.Printf("[conf] DbName          = [%s]\n", cfg.DbName)
	fmt.Printf("[conf] DbUser          = [%s]\n", cfg.DbUser)
	fmt.Printf("[conf] DbPassword      = [REDACTED]\n")

	return &cfg, nil
}

//
// end of file
//
