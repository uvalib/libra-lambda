package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config defines all of the service configuration parameters
type Config struct {
	// configuration needed for mail content
	EtdBaseUrl  string // etd application base URL
	OpenBaseUrl string // open application base URL

	// mailer configuration
	EmailSender    string // the email sender
	SendEmail      bool   // do we send or just log
	DebugRecipient string // the debug recipient

	// SMTP configuration
	SMTPHost string // SMTP hostname
	SMTPPort int    // SMTP port number
	SMTPUser string // SMTP username
	SMTPPass string // SMTP password

	// easystore configuration
	EsDbHost     string // database host
	EsDbPort     int    // database port
	EsDbName     string // database name
	EsDbUser     string // database user
	EsDbPassword string // database password

	// message bus configuration
	BusName    string // the message bus name
	SourceName string // the event source name
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
	cfg.EtdBaseUrl, err = ensureSetAndNonEmpty("ETD_BASE_URL")
	if err != nil {
		return nil, err
	}
	cfg.OpenBaseUrl, err = ensureSetAndNonEmpty("OPEN_BASE_URL")
	if err != nil {
		return nil, err
	}

	cfg.SMTPHost, err = ensureSetAndNonEmpty("SMTP_HOST")
	if err != nil {
		return nil, err
	}
	cfg.SMTPPort, err = envToInt("SMTP_PORT")
	if err != nil {
		return nil, err
	}
	cfg.SMTPUser = envWithDefault("SMTP_USER", "")
	cfg.SMTPPass = envWithDefault("SMTP_PASSWORD", "")

	cfg.EmailSender, err = ensureSetAndNonEmpty("EMAIL_SENDER")
	if err != nil {
		return nil, err
	}
	cfg.SendEmail, err = envToBool("EMAIL_SEND")
	if err != nil {
		return nil, err
	}

	cfg.DebugRecipient = envWithDefault("DEBUG_RECIPIENT", "")

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

	cfg.BusName = envWithDefault("MESSAGE_BUS", "")
	cfg.SourceName = envWithDefault("MESSAGE_SOURCE", "")

	fmt.Printf("[conf] SMTPHost       = [%s]\n", cfg.SMTPHost)
	fmt.Printf("[conf] SMTPPort       = [%d]\n", cfg.SMTPPort)
	fmt.Printf("[conf] SMTPUser       = [%s]\n", cfg.SMTPUser)
	fmt.Printf("[conf] SMTPPass       = [%s]\n", cfg.SMTPPass)

	fmt.Printf("[conf] EmailSender    = [%s]\n", cfg.EmailSender)
	fmt.Printf("[conf] SendEmail      = [%t]\n", cfg.SendEmail)
	fmt.Printf("[conf] DebugRecipient = [%s]\n", cfg.DebugRecipient)

	fmt.Printf("[conf] EsDbHost       = [%s]\n", cfg.EsDbHost)
	fmt.Printf("[conf] EsDbPort       = [%d]\n", cfg.EsDbPort)
	fmt.Printf("[conf] EsDbName       = [%s]\n", cfg.EsDbName)
	fmt.Printf("[conf] EsDbUser       = [%s]\n", cfg.EsDbUser)
	fmt.Printf("[conf] EsDbPassword   = [REDACTED]\n")

	fmt.Printf("[conf] BusName        = [%s]\n", cfg.BusName)
	fmt.Printf("[conf] SourceName     = [%s]\n", cfg.SourceName)

	return &cfg, nil
}
