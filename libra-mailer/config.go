package main

import (
	"fmt"
	"os"
	"strconv"
)

// Config defines all of the service configuration parameters
type Config struct {

	// service endpoint configuration
	MintAuthUrl string // mint auth token endpoint
	UserInfoUrl string // the user information service

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
	cfg.MintAuthUrl, err = ensureSetAndNonEmpty("MINT_AUTH_URL")
	if err != nil {
		return nil, err
	}
	cfg.UserInfoUrl, err = ensureSetAndNonEmpty("USER_INFO_URL")
	if err != nil {
		return nil, err
	}

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

	cfg.EsProxyUrl, err = ensureSetAndNonEmpty("ES_PROXY_URL")
	if err != nil {
		return nil, err
	}

	cfg.BusName = envWithDefault("MESSAGE_BUS", "")

	fmt.Printf("[conf] MintAuthUrl    = [%s]\n", cfg.MintAuthUrl)
	fmt.Printf("[conf] UserInfoUrl    = [%s]\n", cfg.UserInfoUrl)

	fmt.Printf("[conf] EtdBaseUrl     = [%s]\n", cfg.EtdBaseUrl)
	fmt.Printf("[conf] OpenBaseUrl    = [%s]\n", cfg.OpenBaseUrl)

	fmt.Printf("[conf] SMTPHost       = [%s]\n", cfg.SMTPHost)
	fmt.Printf("[conf] SMTPPort       = [%d]\n", cfg.SMTPPort)
	fmt.Printf("[conf] SMTPUser       = [%s]\n", cfg.SMTPUser)
	fmt.Printf("[conf] SMTPPass       = [%s]\n", cfg.SMTPPass)

	fmt.Printf("[conf] EmailSender    = [%s]\n", cfg.EmailSender)
	fmt.Printf("[conf] SendEmail      = [%t]\n", cfg.SendEmail)
	fmt.Printf("[conf] DebugRecipient = [%s]\n", cfg.DebugRecipient)

	fmt.Printf("[conf] EsProxyUrl     = [%s]\n", cfg.EsProxyUrl)
	fmt.Printf("[conf] BusName        = [%s]\n", cfg.BusName)

	return &cfg, nil
}

//
// end of file
//
