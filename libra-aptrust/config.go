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

	// other configuration
	BagNameTemplate   string // the bag name template
	ScratchFilesystem string // the scratch filesystem
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error

	// APTrust submission service configuration
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

	// other configuration
	cfg.BagNameTemplate, err = ensureSetAndNonEmpty("BAG_NAME_TEMPLATE")
	if err != nil {
		return nil, err
	}
	cfg.ScratchFilesystem, err = ensureSetAndNonEmpty("SCRATCH_FS")
	if err != nil {
		return nil, err
	}

	// APTrust submission service configuration
	fmt.Printf("[conf] APTServiceRegister = [%s]\n", cfg.APTServiceRegister)
	fmt.Printf("[conf] APTServiceSubmit   = [%s]\n", cfg.APTServiceSubmit)
	fmt.Printf("[conf] APTServiceClient   = [%s]\n", cfg.APTServiceClient)

	// easystore proxy configuration
	fmt.Printf("[conf] EsProxyUrl         = [%s]\n", cfg.EsProxyUrl)

	// other configuration
	fmt.Printf("[conf] BagNameTemplate    = [%s]\n", cfg.BagNameTemplate)
	fmt.Printf("[conf] ScratchFilesystem  = [%s]\n", cfg.ScratchFilesystem)

	return &cfg, nil
}

//
// end of file
//
