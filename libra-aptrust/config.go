package main

import (
	"fmt"
)

// Config defines all of the service configuration parameters
type Config struct {

	// APTrust submission service configuration
	APTServiceRegister string // url for APTrust submission registration
	APTServiceSubmit   string // url for APTrust submit
	APTServiceClient   string // client identifier for APTrust submit

	// easystore proxy configuration
	EsProxyUrl string // the easystore proxy endpoint
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	cfg.APTServiceRegister, err = ensureSetAndNonEmpty("APT_REGISTER_URL")
	if err != nil {
		return nil, err
	}
	cfg.APTServiceSubmit, err = ensureSetAndNonEmpty("APT_SUBMIT_URL")
	if err != nil {
		return nil, err
	}
	cfg.APTServiceClient, err = ensureSetAndNonEmpty("APT_CLIENT_ID")
	if err != nil {
		return nil, err
	}

	// easystore proxy configuration
	cfg.EsProxyUrl, err = ensureSetAndNonEmpty("ES_PROXY_URL")
	if err != nil {
		return nil, err
	}

	fmt.Printf("[conf] APTServiceRegister = [%s]\n", cfg.APTServiceRegister)
	fmt.Printf("[conf] APTServiceSubmit   = [%s]\n", cfg.APTServiceSubmit)
	fmt.Printf("[conf] APTServiceClient   = [%s]\n", cfg.APTServiceClient)

	// easystore proxy configuration
	fmt.Printf("[conf] EsProxyUrl         = [%s]\n", cfg.EsProxyUrl)

	return &cfg, nil
}

//
// end of file
//
