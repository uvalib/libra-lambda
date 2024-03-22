//
//
//

package main

import (
	"bytes"
	"embed"
	"github.com/uvalib/easystore/uvaeasystore"
	"text/template"
)

// templates holds our templates
//
//go:embed templates/*
var templates embed.FS

func docRender(cfg *Config, work uvaeasystore.EasyStoreObject) ([]byte, error) {

	// read the template
	templateFile := "templates/solr-doc.template"
	templateStr, err := templates.ReadFile(templateFile)
	if err != nil {
		return nil, err
	}

	// parse the templateFile
	tmpl, err := template.New("doc").Parse(string(templateStr))
	if err != nil {
		return nil, err
	}

	type Work struct {
		Id string // work identifier
	}
	type Attributes struct {
		Doc Work
	}

	//	populate the work
	doc := Work{
		Id: work.Id(),
	}

	//	populate the attributes
	//fields := work.Fields()
	attribs := Attributes{
		Doc: doc,
	}

	// render the template
	var renderedBuffer bytes.Buffer
	err = tmpl.Execute(&renderedBuffer, attribs)
	if err != nil {
		return nil, err
	}

	return renderedBuffer.Bytes(), nil
}

//
// end of file
//
