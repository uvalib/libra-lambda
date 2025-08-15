package main

import (
	"fmt"
	"os"
)

// Config defines all of the service configuration parameters
type Config struct {
	BusName    string // message bus name
	SourceName string // message source name
}

func ensureSet(env string) (string, error) {
	val, set := os.LookupEnv(env)

	if set == false {
		err := fmt.Errorf("environment variable not set: [%s]", env)
		fmt.Printf("ERROR: %s\n", err.Error())
		return "", err
	}

	return val, nil
}

func ensureSetAndNonEmpty(env string) (string, error) {
	val, err := ensureSet(env)
	if err != nil {
		return "", err
	}

	if val == "" {
		err := fmt.Errorf("environment variable is empty: [%s]", env)
		fmt.Printf("ERROR: %s\n", err.Error())
		return "", err
	}

	return val, nil
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	cfg.BusName, err = ensureSetAndNonEmpty("MESSAGE_BUS")
	if err != nil {
		return nil, err
	}

	cfg.SourceName, err = ensureSetAndNonEmpty("MESSAGE_SOURCE")
	if err != nil {
		return nil, err
	}

	fmt.Printf("[conf] BusName    = [%s]\n", cfg.BusName)
	fmt.Printf("[conf] SourceName = [%s]\n", cfg.SourceName)

	return &cfg, nil
}

//
// end of file
//
