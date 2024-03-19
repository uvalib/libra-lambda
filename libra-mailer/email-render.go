//
//
//

package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"text/template"
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

	var templateFile string
	var subject string
	switch theType {
	case ETD_OPTIONAL_CAN_DEPOSIT:
		templateFile = "templates/libraetd-optional-can-deposit.template"
		subject = "Access to upload your approved thesis to Libra"

	case ETD_SIS_CAN_DEPOSIT:
		templateFile = "templates/libraetd-sis-can-deposit.template"
		subject = "Access to upload your approved thesis or dissertation to Libra"

	case ETD_SUBMITTED_AUTHOR:
		templateFile = "templates/libraetd-submitted-author.template"
		subject = "Successful deposit of your thesis or dissertation"

	case ETD_SUBMITTED_ADVISOR:
		templateFile = "templates/libraetd-submitted-advisor.template"
		subject = "Successful deposit of your student's thesis"

	case OPEN_SUBMITTED_AUTHOR:
		templateFile = "templates/libraopen-submitted-author.template"
		subject = "Work successfully deposited to Libra"

	default:
		return "", "", fmt.Errorf("unsupported email type")
	}

	// read the template
	templateStr, err := templates.ReadFile(templateFile)
	if err != nil {
		return "", "", err
	}

	// parse the templateFile
	tmpl, err := template.New("email").Parse(string(templateStr))
	if err != nil {
		return "", "", err
	}

	type EmailAttributes struct {
		Recipient string
		Sender    string
	}

	//	// populate the attributes
	attribs := EmailAttributes{
		Recipient: work.Fields()["depositor"],
		Sender:    cfg.EmailSender,
	}

	// render the templateFile
	var renderedBuffer bytes.Buffer
	err = tmpl.Execute(&renderedBuffer, attribs)
	if err != nil {
		return "", "", err
	}

	return subject, renderedBuffer.String(), nil
}

//
// end of file
//
