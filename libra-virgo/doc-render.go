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
	"strings"
	"text/template"
	"time"
)

// templates holds our templates
//
//go:embed templates/*
var templates embed.FS

func docRender(cfg *Config, work uvaeasystore.EasyStoreObject) ([]byte, error) {

	var templateFile = fmt.Sprintf("templates/%s-solr-doc.template", work.Namespace())

	// read the template
	templateStr, err := templates.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	// create the function map
	fMap := template.FuncMap{
		"ToLower":   strings.ToLower,
		"ToUpper":   strings.ToUpper,
		"XmlEncode": XmlEncode,
	}

	// parse the templateFile
	tmpl, err := template.New("doc").Funcs(fMap).Parse(string(templateStr))
	if err != nil {
		return nil, err
	}

	return renderEtd(cfg, tmpl, work)
}

func renderEtd(cfg *Config, tmpl *template.Template, work uvaeasystore.EasyStoreObject) ([]byte, error) {
	type Attributes struct {
		Doi           string // work DOI
		Id            string // work identifier
		IndexDateTime string // current date/time
		PubDate       string // publication date
		PubYear       string // publication year
		ReceivedDate  string // date received
		TitleSort     string // field used by SOLR for sorting/grouping
		Title2Key     string // field used by SOLR for sorting/grouping
		Visibility    string // whether the work is visible

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
		Work:          *meta,
		Doi:           fields["doi"],
		Id:            work.Id(),
		IndexDateTime: time.Now().Format("20060102150405"),
		PubDate:       fields["publish-date"],
		PubYear:       extractYYYY(fields["publish-date"]),
		ReceivedDate:  extractYYYY(fields["create-date"]),
		TitleSort:     titleForSort,
		Title2Key:     title2Key,
		Visibility:    workVisibility(fields),
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
