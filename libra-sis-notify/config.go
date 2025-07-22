package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config defines all of the service configuration parameters
type Config struct {

	// service endpoint configuration
	MintAuthUrl  string // mint auth token endpoint
	SisNotifyUrl string // the sis notify service

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
	cfg.SisNotifyUrl, err = ensureSetAndNonEmpty("SIS_NOTIFY_URL")
	if err != nil {
		return nil, err
	}

	cfg.EsProxyUrl, err = ensureSetAndNonEmpty("ES_PROXY_URL")
	if err != nil {
		return nil, err
	}

	cfg.BusName = envWithDefault("MESSAGE_BUS", "")

	fmt.Printf("[conf] MintAuthUrl  = [%s]\n", cfg.MintAuthUrl)
	fmt.Printf("[conf] SisNotifyUrl = [%s]\n", cfg.SisNotifyUrl)
	fmt.Printf("[conf] EsProxyUrl   = [%s]\n", cfg.EsProxyUrl)
	fmt.Printf("[conf] BusName      = [%s]\n", cfg.BusName)

	return &cfg, nil
}
