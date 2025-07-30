package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config defines all of the service configuration parameters
type Config struct {

	// service endpoint configuration
	MintAuthUrl         string // mint auth token endpoint
	OrcidGetDetailsUrl  string // get orcid get details endpoint
	OrcidSetActivityUrl string // get orcid set activity endpoint

	// easystore proxy configuration
	EsProxyUrl string // the easystore proxy endpoint

	// message bus configuration
	BusName string // the message bus name
}

func envWithDefault(env string, defaultValue string) string {
	val, set := os.LookupEnv(env)

	if set == false {
		fmt.Printf("INFO: environment variable not set: [%s] using default value [%s]\n", env, defaultValue)
		return defaultValue
	}

	return val
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
		return "", err
	}

	return val, nil
}

func envToInt(env string) (int, error) {

	number, err := ensureSetAndNonEmpty(env)
	if err != nil {
		return -1, err
	}

	n, err := strconv.Atoi(number)
	if err != nil {
		return -1, err
	}
	return n, nil
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error
	cfg.MintAuthUrl, err = ensureSetAndNonEmpty("MINT_AUTH_URL")
	if err != nil {
		return nil, err
	}
	cfg.OrcidGetDetailsUrl, err = ensureSetAndNonEmpty("ORCID_GET_DETAILS_URL")
	if err != nil {
		return nil, err
	}
	cfg.OrcidSetActivityUrl, err = ensureSetAndNonEmpty("ORCID_SET_ACTIVITY_URL")
	if err != nil {
		return nil, err
	}

	cfg.EsProxyUrl, err = ensureSetAndNonEmpty("ES_PROXY_URL")
	if err != nil {
		return nil, err
	}

	//cfg.BusName = envWithDefault("MESSAGE_BUS", "")

	fmt.Printf("[conf] MintAuthUrl         = [%s]\n", cfg.MintAuthUrl)
	fmt.Printf("[conf] OrcidGetDetailsUrl  = [%s]\n", cfg.OrcidGetDetailsUrl)
	fmt.Printf("[conf] OrcidSetActivityUrl = [%s]\n", cfg.OrcidSetActivityUrl)

	fmt.Printf("[conf] EsProxyUrl          = [%s]\n", cfg.EsProxyUrl)
	//fmt.Printf("[conf] BusName             = [%s]\n", cfg.BusName)

	return &cfg, nil
}

//
// end of file
//
