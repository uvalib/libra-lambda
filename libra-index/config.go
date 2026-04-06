package main

import (
	"fmt"
)

// Config defines all of the service configuration parameters
type Config struct {
	// index endpoint configuration
	IndexUpdateUrl string // the index update URL

	// easystore proxy configuration
	EsProxyUrl string // the easystore proxy endpoint

	// message bus configuration
	BusName string // the message bus name
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	cfg.IndexUpdateUrl, err = ensureSetAndNonEmpty("INDEX_UPDATE_URL")
	if err != nil {
		return nil, err
	}

	cfg.EsProxyUrl, err = ensureSetAndNonEmpty("ES_PROXY_URL")
	if err != nil {
		return nil, err
	}

	cfg.BusName = envWithDefault("MESSAGE_BUS", "")

	fmt.Printf("[conf] IndexUpdateUrl    = [%s]\n", cfg.IndexUpdateUrl)
	fmt.Printf("[conf] EsProxyUrl        = [%s]\n", cfg.EsProxyUrl)
	fmt.Printf("[conf] BusName           = [%s]\n", cfg.BusName)

	return &cfg, nil
}

//
// end of file
//
