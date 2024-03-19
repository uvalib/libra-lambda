package main

import (
	"crypto/tls"
	"fmt"
	"gopkg.in/gomail.v2"
	"log"
	"time"
)

func sendEmail(cfg *Config, subject string, recipient string, cc []string, body string) error {

	mail := gomail.NewMessage()
	mail.SetHeader("MIME-version", "1.0")
	mail.SetHeader("Content-Type", "text/plain; charset=\"UTF-8\"")
	mail.SetHeader("Subject", subject)
	mail.SetHeader("To", recipient)
	mail.SetHeader("From", cfg.EmailSender)

	if len(cc) != 0 {
		mail.SetHeader("Cc", cc...)
	}

	mail.SetBody("text/plain", body)

	if cfg.SendEmail == false {
		fmt.Printf("INFO: Email is in debug mode. Logging message instead of sending\n")
		fmt.Printf("INFO: ==========================================================\n")
		_, _ = mail.WriteTo(log.Writer())
		fmt.Printf("\nINFO: ==========================================================\n")
		return nil
	}

	var dialer gomail.Dialer
	fmt.Printf("INFO: sending '%s' email to '%s'\n", subject, recipient)
	if cfg.SMTPPass != "" {
		fmt.Printf("INFO: sending email with auth\n")
		dialer = gomail.Dialer{Host: cfg.SMTPHost, Port: cfg.SMTPPort, Username: cfg.SMTPUser, Password: cfg.SMTPPass}
	} else {
		fmt.Printf("INFO: sending email with no auth\n")
		dialer = gomail.Dialer{Host: cfg.SMTPHost, Port: cfg.SMTPPort}
	}
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return dialAndSend(dialer, mail)
}

func dialAndSend(dialer gomail.Dialer, mail *gomail.Message) error {

	retryCount := 3
	retrySleepTime := 1 * time.Second
	currentCount := 0

	for {
		err := dialer.DialAndSend(mail)
		if err == nil {
			return nil
		}
		currentCount++

		// break when tried too many times
		if currentCount >= retryCount {
			err = fmt.Errorf("email send failed with error (%s), giving up", err)
			return err
		}

		fmt.Printf("WARNING: email send failed with error (%s), retrying...\n", err)

		// sleep for a bit before retrying
		time.Sleep(retrySleepTime)
	}
}

//
// end of file
//
