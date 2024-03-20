package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config defines all of the service configuration parameters
type Config struct {
	// DOI service configuration
	IdService IdServiceConfig
	DOIBaseURL string // base url for DOIs

	// easystore configuration
	EsDbHost     string // database host
	EsDbPort     int    // database port
	EsDbName     string // database name
	EsDbUser     string // database user
	EsDbPassword string // database password
}

type IdServiceConfig struct {
	BaseURL string
	Shoulder string
	User string
	Password string
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

func envToBool(env string) (bool, error) {

	str, err := ensureSetAndNonEmpty(env)
	if err != nil {
		return false, err
	}

	b, err := strconv.ParseBool(str)
	if err != nil {
		return false, err
	}
	return b, nil
}

// loadConfiguration will load the service configuration from env/cmdline
// and return a pointer to it. Any failures are fatal.
func loadConfiguration() (*Config, error) {

	var cfg Config

	var err error

	cfg.IdService.BaseURL, err = ensureSet("ID_SERVICE_BASE")
	if err != nil {
		return nil, err
	}

	cfg.IdService.Shoulder, err = ensureSet("ID_SERVICE_SHOULDER")
	if err != nil {
		return nil, err
	}


	cfg.IdService.User, err = ensureSet("ID_SERVICE_USER")
	if err != nil {
		return nil, err
	}

	cfg.IdService.Password, err = ensureSet("ID_SERVICE_PASSWORD")
	if err != nil {
		return nil, err
	}

	cfg.DOIBaseURL, err = ensureSet("DOI_BASE_URL")
	if err != nil {
		return nil, err
	}

	cfg.EsDbHost, err = ensureSet("ES_DBHOST")
	if err != nil {
		return nil, err
	}
	cfg.EsDbPort, err = envToInt("ES_DBPORT")
	if err != nil {
		return nil, err
	}
	cfg.EsDbName, err = ensureSet("ES_DBNAME")
	if err != nil {
		return nil, err
	}
	cfg.EsDbUser, err = ensureSet("ES_DBUSER")
	if err != nil {
		return nil, err
	}
	cfg.EsDbPassword, err = ensureSet("ES_DBPASSWORD")
	if err != nil {
		return nil, err
	}


	fmt.Printf("[conf] EsDbHost       = [%s]\n", cfg.EsDbHost)
	fmt.Printf("[conf] EsDbPort       = [%d]\n", cfg.EsDbPort)
	fmt.Printf("[conf] EsDbName       = [%s]\n", cfg.EsDbName)
	fmt.Printf("[conf] EsDbUser       = [%s]\n", cfg.EsDbUser)
	fmt.Printf("[conf] EsDbPassword   = [REDACTED]\n")

	return &cfg, nil
}
