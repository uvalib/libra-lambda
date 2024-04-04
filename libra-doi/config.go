package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Config defines all of the service configuration parameters
type Config struct {
	IDService  IDServiceConfig
	DOIBaseURL string // base url for DOIs

	ETDNamespace  Namespace
	OpenNamespace Namespace

	ResourceTypes []ResourceType

	// easystore configuration
	EsDbHost     string // database host
	EsDbPort     int    // database port
	EsDbName     string // database name
	EsDbUser     string // database user
	EsDbPassword string // database password

	BusName    string // name of the bus
	SourceName string // name of the source

	httpClient http.Client // shared http client

}

// Namespace such as libraopen or libraetd
type Namespace struct {
	Name string
	Path string
}

type ResourceType struct {
	Value    string
	Label    string
	Category string
	Oa       bool
	Etd      bool
}

// IDServiceConfig for DOI service
type IDServiceConfig struct {
	BaseURL  string
	Shoulder string
	User     string
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

	cfg.IDService.BaseURL, err = ensureSet("ID_SERVICE_BASE")
	if err != nil {
		return nil, err
	}

	cfg.IDService.Shoulder, err = ensureSet("ID_SERVICE_SHOULDER")
	if err != nil {
		return nil, err
	}

	cfg.IDService.User, err = ensureSet("ID_SERVICE_USER")
	if err != nil {
		return nil, err
	}

	cfg.IDService.Password, err = ensureSet("ID_SERVICE_PASSWORD")
	if err != nil {
		return nil, err
	}

	cfg.DOIBaseURL, err = ensureSet("DOI_BASE_URL")
	if err != nil {
		return nil, err
	}

	cfg.OpenNamespace = Namespace{
		Name: "libraopen",
		Path: "oa",
	}
	cfg.ETDNamespace = Namespace{
		Name: "libraetd",
		Path: "etd",
	}

	log.Printf("INFO: load resource types")
	bytes, err := os.ReadFile("libra-doi/data/resourceTypes.json")
	if err != nil {
		log.Printf("ERROR: unable to load resourceTypes: %s", err.Error())
	} else {
		err = json.Unmarshal(bytes, &cfg.ResourceTypes)
		if err != nil {
			log.Printf("ERROR: unable to parse resourceTypes.json: %s", err.Error())
		}
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

	cfg.BusName = envWithDefault("MESSAGE_BUS", "")
	cfg.SourceName = envWithDefault("MESSAGE_SOURCE", "")

	fmt.Printf("[conf] EsDbHost       = [%s]\n", cfg.EsDbHost)
	fmt.Printf("[conf] EsDbPort       = [%d]\n", cfg.EsDbPort)
	fmt.Printf("[conf] EsDbName       = [%s]\n", cfg.EsDbName)
	fmt.Printf("[conf] EsDbUser       = [%s]\n", cfg.EsDbUser)
	fmt.Printf("[conf] EsDbPassword   = [REDACTED]\n")
	fmt.Printf("[conf] BusName        = [%s]\n", cfg.BusName)
	fmt.Printf("[conf] SourceName     = [%s]\n", cfg.SourceName)

	fmt.Printf("[conf] DOIBaseURL        = [%s]\n", cfg.DOIBaseURL)
	fmt.Printf("[conf] IDServiceBase     = [%s]\n", cfg.IDService.BaseURL)
	fmt.Printf("[conf] IDServiceShoulder = [%s]\n", cfg.IDService.Shoulder)
	fmt.Printf("[conf] IDServiceUser     = [%s]\n", cfg.IDService.User)
	fmt.Printf("[conf] IDServicePassword = [REDACTED]\n")

	return &cfg, nil
}
