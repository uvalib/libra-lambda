package main

import (
	"fmt"
	"os"
	"strconv"
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
