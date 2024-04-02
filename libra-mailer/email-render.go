//
//
//

package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	librametadata "github.com/uvalib/libra-metadata"
	"text/template"
	"time"
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

// values extracted from the work used by the template rendering
type Work struct {
	Degree string // degree name
	Title  string // work title
}

// used for time stuff
var location, _ = time.LoadLocation("America/New_York")

func emailSubjectAndBody(cfg *Config, theType emailType, obj uvaeasystore.EasyStoreObject) (string, string, error) {

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

	type Attributes struct {
		Work Work

		Advisee                  string // for mail sent to registrar
		Availability             string // display version of visibility
		BaseUrl                  string // libra base URL
		Doi                      string // work DOI
		EmbargoReleaseDate       string // embargo release date
		EmbargoReleaseVisibility string // embargo release visibility
		IsSis                    bool   // is this a SIS thesis
		License                  string // work license
		Recipient                string // mail recipient
		Sender                   string // mail sender
		Visibility               string // work visibility
	}

	// populate the work
	work, err := extractAtributes(obj)
	if err != nil {
		return "", "", err
	}

	fields := obj.Fields()

	// determine the availability string
	availability := determineAvailability(fields)

	// determine the base URL
	baseUrl := cfg.EtdBaseUrl
	if obj.Namespace() == libraOpenNamespace {
		baseUrl = cfg.OpenBaseUrl
	}

	//	populate the attributes
	attribs := Attributes{
		Work: *work,

		Advisee:                  fields["depositor"],
		Availability:             availability,
		BaseUrl:                  baseUrl,
		Doi:                      fields["doi"],
		EmbargoReleaseDate:       fields["embargo-release"],
		EmbargoReleaseVisibility: fields["embargo-release-visibility"],
		IsSis:                    fields["source"] == "sis",
		Recipient:                fields["depositor"],
		Sender:                   cfg.EmailSender,
		Visibility:               fields["visibility"],
	}

	// render the template
	var renderedBuffer bytes.Buffer
	err = tmpl.Execute(&renderedBuffer, attribs)
	if err != nil {
		return "", "", err
	}

	return subject, renderedBuffer.String(), nil
}

func extractAtributes(obj uvaeasystore.EasyStoreObject) (*Work, error) {

	switch obj.Namespace() {
	case libraEtdNamespace:
		return extractEtdAtributes(obj)
	case libraOpenNamespace:
		return extractOpenAtributes(obj)
	default:
		return nil, fmt.Errorf("unsupported namespace")
	}
}

func extractEtdAtributes(obj uvaeasystore.EasyStoreObject) (*Work, error) {

	// extract the metadata
	if obj.Metadata() == nil {
		fmt.Printf("ERROR: unable to get metadata payload for ns/oid [%s/%s]\n", obj.Namespace(), obj.Id())
		return nil, ErrNoMetadata
	}

	md := obj.Metadata()
	pl, err := md.Payload()
	if err != nil {
		return nil, err
	}
	meta, err := librametadata.ETDWorkFromBytes(pl)
	if err != nil {
		return nil, err
	}

	// populate the work
	work := Work{
		Degree: meta.Degree,
		Title:  meta.Title,
	}

	return &work, nil
}

func extractOpenAtributes(obj uvaeasystore.EasyStoreObject) (*Work, error) {

	// extract the metadata
	if obj.Metadata() == nil {
		fmt.Printf("ERROR: unable to get metadata payload for ns/oid [%s/%s]\n", obj.Namespace(), obj.Id())
		return nil, ErrNoMetadata
	}

	md := obj.Metadata()
	pl, err := md.Payload()
	if err != nil {
		return nil, err
	}
	meta, err := librametadata.OAWorkFromBytes(pl)
	if err != nil {
		return nil, err
	}

	// populate the work
	work := Work{
		Degree: "None", // no degree program for an open item
		Title:  meta.Title,
	}

	return &work, nil
}

func determineAvailability(fields uvaeasystore.EasyStoreObjectFields) string {

	ava := "public access immediately"

	// if we have an embargo release date
	if len(fields["embargo-release"]) != 0 {

		format := "2006-01-02T15:04:05+00:00" // yeah, crap right
		dt, err := time.ParseInLocation(format, fields["embargo-release"], location)
		if err != nil {
			return ava + " (cannot decode embargo release date)"
		}

		// are we still under embargo
		if dt.After(time.Now()) {
			ava = fmt.Sprintf("public access on %s", dt.Format("%B %-d, %Y"))
		}
	}

	return ava
}

//
// end of file
//
