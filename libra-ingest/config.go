package main

import (
	"fmt"
)

// Config defines all of the service configuration parameters
type Config struct {

	// service endpoint configuration
	MintAuthUrl        string // mint auth token endpoint
	UserInfoUrl        string // the user information service
	SisIngestUrl       string // the sis ingest service
	SisIngestStateName string // the sis ingest ssm state name

	// easystore proxy configuration
	EsProxyUrl string // the easystore proxy endpoint

	// message bus configuration
	BusName string // the message bus name
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	cfg.MintAuthUrl, err = ensureSetAndNonEmpty("MINT_AUTH_URL")
	if err != nil {
		return nil, err
	}
	cfg.UserInfoUrl, err = ensureSetAndNonEmpty("USER_INFO_URL")
	if err != nil {
		return nil, err
	}

	cfg.SisIngestUrl, err = ensureSetAndNonEmpty("SIS_INGEST_URL")
	if err != nil {
		return nil, err
	}
	cfg.SisIngestStateName, err = ensureSetAndNonEmpty("SIS_INGEST_STATE_NAME")
	if err != nil {
		return nil, err
	}

	cfg.EsProxyUrl, err = ensureSetAndNonEmpty("ES_PROXY_URL")
	if err != nil {
		return nil, err
	}

	cfg.BusName = envWithDefault("MESSAGE_BUS", "")

	fmt.Printf("[conf] EsProxyUrl              = [%s]\n", cfg.EsProxyUrl)
	fmt.Printf("[conf] MintAuthUrl             = [%s]\n", cfg.MintAuthUrl)
	fmt.Printf("[conf] UserInfoUrl             = [%s]\n", cfg.UserInfoUrl)
	fmt.Printf("[conf] SisIngestUrl            = [%s]\n", cfg.SisIngestUrl)
	fmt.Printf("[conf] SisIngestStateName      = [%s]\n", cfg.SisIngestStateName)
	fmt.Printf("[conf] BusName                 = [%s]\n", cfg.BusName)

	return &cfg, nil
}

//
// end of file
//
