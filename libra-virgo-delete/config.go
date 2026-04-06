package main

import (
	"fmt"
)

// Config defines all of the service configuration parameters
type Config struct {
	// upload configuration
	BucketName        string // the bucket name
	BucketKeyTemplate string // the bucket key template
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	cfg.BucketName, err = ensureSetAndNonEmpty("BUCKET_NAME")
	if err != nil {
		return nil, err
	}
	cfg.BucketKeyTemplate, err = ensureSetAndNonEmpty("BUCKET_KEY_TEMPLATE")
	if err != nil {
		return nil, err
	}

	fmt.Printf("[conf] BucketName        = [%s]\n", cfg.BucketName)
	fmt.Printf("[conf] BucketKeyTemplate = [%s]\n", cfg.BucketKeyTemplate)

	return &cfg, nil
}

//
// end of file
//
