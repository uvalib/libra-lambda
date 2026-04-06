package main

import (
	"fmt"
)

// Config defines all of the service configuration parameters
type Config struct {
	// index endpoint configuration
	IndexDeleteUrl string // the index update URL
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	cfg.IndexDeleteUrl, err = ensureSetAndNonEmpty("INDEX_DELETE_URL")
	if err != nil {
		return nil, err
	}

	fmt.Printf("[conf] IndexDeleteUrl = [%s]\n", cfg.IndexDeleteUrl)

	return &cfg, nil
}

//
// end of file
//
