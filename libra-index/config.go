package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config defines all of the service configuration parameters
type Config struct {
	// index endpoint configuration
	IndexUpdateUrl string // the index update URL

	// easystore proxy configuration
	EsProxyUrl string // the easystore proxy endpoint
	//EsDbHost     string // database host
	//EsDbPort     int    // database port
	//EsDbName     string // database name
	//EsDbUser     string // database user
	//EsDbPassword string // database password

	// message bus configuration
	//BusName    string // the message bus name
	//SourceName string // the event source name
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
	cfg.IndexUpdateUrl, err = ensureSetAndNonEmpty("INDEX_UPDATE_URL")
	if err != nil {
		return nil, err
	}

	cfg.EsProxyUrl, err = ensureSetAndNonEmpty("ES_PROXY_URL")
	if err != nil {
		return nil, err
	}
	//cfg.EsDbHost, err = ensureSetAndNonEmpty("ES_DBHOST")
	//if err != nil {
	//	return nil, err
	//}
	//cfg.EsDbPort, err = envToInt("ES_DBPORT")
	//if err != nil {
	//	return nil, err
	//}
	//cfg.EsDbName, err = ensureSetAndNonEmpty("ES_DBNAME")
	//if err != nil {
	//	return nil, err
	//}
	//cfg.EsDbUser, err = ensureSetAndNonEmpty("ES_DBUSER")
	//if err != nil {
	//	return nil, err
	//}
	//cfg.EsDbPassword, err = ensureSetAndNonEmpty("ES_DBPASSWORD")
	//if err != nil {
	//	return nil, err
	//}

	//cfg.BusName = envWithDefault("MESSAGE_BUS", "")
	//cfg.SourceName = envWithDefault("MESSAGE_SOURCE", "")

	fmt.Printf("[conf] IndexUpdateUrl    = [%s]\n", cfg.IndexUpdateUrl)
	fmt.Printf("[conf] EsProxyUrl        = [%s]\n", cfg.EsProxyUrl)
	//fmt.Printf("[conf] EsDbHost          = [%s]\n", cfg.EsDbHost)
	//fmt.Printf("[conf] EsDbPort          = [%d]\n", cfg.EsDbPort)
	//fmt.Printf("[conf] EsDbName          = [%s]\n", cfg.EsDbName)
	//fmt.Printf("[conf] EsDbUser          = [%s]\n", cfg.EsDbUser)
	//fmt.Printf("[conf] EsDbPassword      = [REDACTED]\n")

	//fmt.Printf("[conf] BusName           = [%s]\n", cfg.BusName)
	//fmt.Printf("[conf] SourceName        = [%s]\n", cfg.SourceName)

	return &cfg, nil
}
