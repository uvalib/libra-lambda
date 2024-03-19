package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config defines all of the service configuration parameters
type Config struct {
	// mailer configuration
	SMTPHost string // SMTP hostname
	SMTPPort int    // SMTP port number
	SMTPUser string // SMTP username
	SMTPPass string // SMTP password

	EmailSender string // the email sender
	SendEmail   bool   // do we send or just log

	// easystore configuration
	EsDbHost     string // database host
	EsDbPort     int    // database port
	EsDbName     string // database name
	EsDbUser     string // database user
	EsDbPassword string // database password
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

	cfg.EmailSender, err = ensureSetAndNonEmpty("NOTIFY_EMAIL_SENDER")
	if err != nil {
		return nil, err
	}
	cfg.SendEmail, err = envToBool("NOTIFY_EMAIL_SEND")
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

	fmt.Printf("[CONFIG] SMTPHost          = [%s]\n", cfg.SMTPHost)
	fmt.Printf("[CONFIG] SMTPPort          = [%d]\n", cfg.SMTPPort)
	fmt.Printf("[CONFIG] SMTPUser          = [%s]\n", cfg.SMTPUser)
	fmt.Printf("[CONFIG] SMTPPass          = [%s]\n", cfg.SMTPPass)

	fmt.Printf("[CONFIG] EmailSender       = [%s]\n", cfg.EmailSender)
	fmt.Printf("[CONFIG] SendEmail         = [%t]\n", cfg.SendEmail)

	fmt.Printf("[CONFIG] EsDbHost          = [%s]\n", cfg.EsDbHost)
	fmt.Printf("[CONFIG] EsDbPort          = [%d]\n", cfg.EsDbPort)
	fmt.Printf("[CONFIG] EsDbName          = [%s]\n", cfg.EsDbName)
	fmt.Printf("[CONFIG] EsDbUser          = [%s]\n", cfg.EsDbUser)
	fmt.Printf("[CONFIG] EsDbPassword      = [REDACTED]\n")

	return &cfg, nil
}
