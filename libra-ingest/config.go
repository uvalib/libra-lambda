package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config defines all of the service configuration parameters
type Config struct {
	MintAuthUrl string // mint auth token endpoint

	SisIngestUrl            string // the sis ingest service
	OptionalIngestUrl       string // the optional ingest service
	SisIngestStateName      string // the sis ingest ssm state name
	OptionalIngestStateName string // the optional ingest ssm state name

	// easystore configuration
	EsDbHost     string // database host
	EsDbPort     int    // database port
	EsDbName     string // database name
	EsDbUser     string // database user
	EsDbPassword string // database password
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

	cfg.SisIngestUrl, err = ensureSetAndNonEmpty("SIS_INGEST_URL")
	if err != nil {
		return nil, err
	}
	cfg.OptionalIngestUrl, err = ensureSetAndNonEmpty("OPTIONAL_INGEST_URL")
	if err != nil {
		return nil, err
	}
	cfg.SisIngestStateName, err = ensureSetAndNonEmpty("SIS_INGEST_STATE_NAME")
	if err != nil {
		return nil, err
	}
	cfg.OptionalIngestStateName, err = ensureSetAndNonEmpty("OPTIONAL_INGEST_STATE_NAME")
	if err != nil {
		return nil, err
	}

	cfg.EsDbHost, err = ensureSetAndNonEmpty("ES_DBHOST")
	if err != nil {
		return nil, err
	}
	cfg.EsDbPort, err = envToInt("ES_DBPORT")
	if err != nil {
		return nil, err
	}
	cfg.EsDbName, err = ensureSetAndNonEmpty("ES_DBNAME")
	if err != nil {
		return nil, err
	}
	cfg.EsDbUser, err = ensureSetAndNonEmpty("ES_DBUSER")
	if err != nil {
		return nil, err
	}
	cfg.EsDbPassword, err = ensureSetAndNonEmpty("ES_DBPASSWORD")
	if err != nil {
		return nil, err
	}

	fmt.Printf("[conf] MintAuthUrl             = [%s]\n", cfg.MintAuthUrl)

	fmt.Printf("[conf] SisIngestUrl            = [%s]\n", cfg.SisIngestUrl)
	fmt.Printf("[conf] OptionalIngestUrl       = [%d]\n", cfg.OptionalIngestUrl)
	fmt.Printf("[conf] SisIngestStateName      = [%s]\n", cfg.SisIngestStateName)
	fmt.Printf("[conf] OptionalIngestStateName = [%s]\n", cfg.OptionalIngestStateName)

	fmt.Printf("[conf] EsDbHost                = [%s]\n", cfg.EsDbHost)
	fmt.Printf("[conf] EsDbPort                = [%d]\n", cfg.EsDbPort)
	fmt.Printf("[conf] EsDbName                = [%s]\n", cfg.EsDbName)
	fmt.Printf("[conf] EsDbUser                = [%s]\n", cfg.EsDbUser)
	fmt.Printf("[conf] EsDbPassword            = [REDACTED]\n")

	return &cfg, nil
}
