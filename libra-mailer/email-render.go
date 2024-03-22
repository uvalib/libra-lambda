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
	ETD_OPTIONAL_INVITATION emailType = iota
	ETD_SIS_INVITATION
	ETD_SUBMITTED_AUTHOR
	ETD_SUBMITTED_ADVISOR
	OPEN_SUBMITTED_AUTHOR
)

func emailSubjectAndBody(cfg *Config, theType emailType, work uvaeasystore.EasyStoreObject) (string, string, error) {

	var templateFile string
	var subject string
	switch theType {
	case ETD_OPTIONAL_INVITATION:
		templateFile = "templates/libraetd-optional-invitation.template"
		subject = "Access to upload your approved thesis to Libra"

	case ETD_SIS_INVITATION:
		templateFile = "templates/libraetd-sis-invitation.template"
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

	type Work struct {
		Degree string // degree name
		Title  string // work title
	}

	type Attributes struct {
		Doc Work

		Advisee                  string // for mail sent to registrar
		Availability             string // FIXME
		BaseUrl                  string // libra base URL
		Doi                      string // work DOI
		EmbargoReleaseDate       string // embargo release date
		EmbargoReleaseVisibility string // FIXME
		IsSis                    bool   // is this a SIS thesis
		License                  string // work license
		Recipient                string // mail recipient
		Sender                   string // mail sender
		Visibility               string // work visibility
	}

	// populate the work
	doc := Work{
		Degree: "placeholder degree", // FIXME
		Title:  "placeholder title",  // FIXME
	}

	//	populate the attributes
	fields := work.Fields()
	attribs := Attributes{
		Doc: doc,

		Advisee:            fields["depositor"],
		BaseUrl:            "https://bla.library.virginia.edu",
		Doi:                fields["doi"],
		EmbargoReleaseDate: fields["embargo-release"],
		IsSis:              fields["source"] == "sis",
		Recipient:          fields["depositor"],
		Sender:             cfg.EmailSender,
		Visibility:         fields["visibility"],
	}

	// render the template
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
