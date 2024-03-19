//
//
//

package main

import (
	"embed"
	"github.com/uvalib/easystore/uvaeasystore"
)

// templates holds our email templates
//
//go:embed templates/*
var templates embed.FS

// the kind of email to render
type emailType int

const (
	ETD_OPTIONAL_CAN_DEPOSIT emailType = iota
	ETD_SIS_CAN_DEPOSIT
	ETD_SUBMITTED_AUTHOR
	ETD_SUBMITTED_ADVISOR
	OPEN_SUBMITTED_AUTHOR
)

func emailSubjectAndBody(cfg *Config, theType emailType, work uvaeasystore.EasyStoreObject) (string, string, error) {

	var template string
	var subject string
	switch theType {
	case ETD_OPTIONAL_CAN_DEPOSIT:
		template = "templates/libraetd-optional-can-deposit.template"
		subject = ""

	case ETD_SIS_CAN_DEPOSIT:
		template = "templates/libraetd-sis-can-deposit.template"
		subject = ""

	case ETD_SUBMITTED_AUTHOR:
		template = "templates/libraetd-submitted-author.template"
		subject = ""

	case ETD_SUBMITTED_ADVISOR:
		template = "templates/libraetd-submitted-advisor.template"
		subject = ""

	case OPEN_SUBMITTED_AUTHOR:
		template = "templates/libraopen-submitted-author.template"
		subject = ""

	default:
	}

	type EmailAttributes struct {
		Recipient   string
		FailedCount int
		Details     string
	}

	//
	//	// parse the template
	//	tmpl, err := template.New("email").Parse(cfg.EmailTemplate)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//	// populate the attributes
	//	attribs := EmailAttributes{
	//		Recipient: cfg.EmailRecipient,
	//		FailedCount: len(messageList),
	//	}
	//
	//	// render the template
	//	var renderedBuffer bytes.Buffer
	//	err = tmpl.Execute(&renderedBuffer, attribs)
	//	if err != nil {
	//		return "", err
	//	}
	//
	//return subject, renderedBuffer.String(), nil
	// TEMP
	return subject, template, nil
}

//
// end of file
//
