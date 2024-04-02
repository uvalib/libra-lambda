//
//
//

package main

import (
	"bytes"
	"embed"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"github.com/uvalib/libra-metadata"
	"text/template"
)

// templates holds our templates
//
//go:embed templates/*
var templates embed.FS

func docRender(cfg *Config, work uvaeasystore.EasyStoreObject) ([]byte, error) {

	var templateFile string
	switch work.Namespace() {
	case libraEtdNamespace, libraOpenNamespace:
		templateFile = fmt.Sprintf("templates/%s-solr-doc.template", work.Namespace())
	default:
		return nil, fmt.Errorf("unsupported namespace (%s)", work.Namespace())
	}

	// read the template
	templateStr, err := templates.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	// parse the templateFile
	tmpl, err := template.New("doc").Parse(string(templateStr))
	if err != nil {
		return nil, err
	}

	switch work.Namespace() {

	case libraOpenNamespace:
		return renderOpen(cfg, tmpl, work)
	case libraEtdNamespace:
		return renderEtd(cfg, tmpl, work)
	}

	return nil, fmt.Errorf("unsupported namespace (%s)", work.Namespace())
}

func renderEtd(cfg *Config, tmpl *template.Template, work uvaeasystore.EasyStoreObject) ([]byte, error) {
	type Attributes struct {
		Doi       string // work DOI
		Id        string // work identifier
		TitleSort string // field used by SOLR for sorting/grouping
		Title2Key string // field used by SOLR for sorting/grouping
		Title3Key string // field used by SOLR for sorting/grouping

		Work librametadata.ETDWork
	}

	// extract the metadata
	if work.Metadata() == nil {
		fmt.Printf("ERROR: unable to get metadata payload for ns/oid [%s/%s]\n", work.Namespace(), work.Id())
		return nil, ErrNoMetadata
	}

	md := work.Metadata()
	pl, err := md.Payload()
	if err != nil {
		return nil, err
	}
	meta, err := librametadata.ETDWorkFromBytes(pl)
	if err != nil {
		return nil, err
	}

	//	populate the attributes
	fields := work.Fields()
	languages := []string{meta.Language}
	titleForSort := titleSort(meta.Title, languages)
	title2Key := titleForSort + titleSuffix(meta.Author.FirstName, meta.Author.LastName)
	attribs := Attributes{
		Work:      *meta,
		Doi:       fields["doi"],
		Id:        work.Id(),
		TitleSort: titleForSort,
		Title2Key: title2Key,
		Title3Key: title2Key, // same as above
	}

	// render the template
	var renderedBuffer bytes.Buffer
	err = tmpl.Execute(&renderedBuffer, attribs)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", renderedBuffer.String())
	return renderedBuffer.Bytes(), nil
}

func renderOpen(cfg *Config, tmpl *template.Template, work uvaeasystore.EasyStoreObject) ([]byte, error) {

	type Attributes struct {
		Doi       string // work DOI
		Id        string // work identifier
		PubYear   string // publication year
		TitleSort string // field used by SOLR for sorting/grouping
		Title2Key string // field used by SOLR for sorting/grouping
		Title3Key string // field used by SOLR for sorting/grouping

		Work librametadata.OAWork
	}

	// extract the metadata
	if work.Metadata() == nil {
		fmt.Printf("ERROR: unable to get metadata payload for ns/oid [%s/%s]\n", work.Namespace(), work.Id())
		return nil, ErrNoMetadata
	}

	md := work.Metadata()
	pl, err := md.Payload()
	if err != nil {
		return nil, err
	}
	meta, err := librametadata.OAWorkFromBytes(pl)
	if err != nil {
		return nil, err
	}

	//	populate the attributes
	fields := work.Fields()
	titleForSort := titleSort(meta.Title, meta.Languages)
	title2Key := titleForSort + titleSuffix(meta.Authors[0].FirstName, meta.Authors[0].LastName)
	attribs := Attributes{
		Work:      *meta,
		Doi:       fields["doi"],
		Id:        work.Id(),
		PubYear:   extractYYMMDD(meta.PublicationDate),
		TitleSort: titleForSort,
		Title2Key: title2Key,
		Title3Key: title2Key, // same as above
	}

	// render the template
	var renderedBuffer bytes.Buffer
	err = tmpl.Execute(&renderedBuffer, attribs)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", renderedBuffer.String())
	return renderedBuffer.Bytes(), nil
}

//
// end of file
//
