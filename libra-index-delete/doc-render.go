//
//
//

package main

import (
	"bytes"
	"embed"
	"fmt"
	"text/template"
)

// templates holds our templates
//
//go:embed templates/*
var templates embed.FS

func docRender(namespace string, oid string) ([]byte, error) {

	var templateFile = fmt.Sprintf("templates/%s-del-doc.template", namespace)

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

	type Attributes struct {
		Id string // work identifier
	}
	attribs := Attributes{
		Id: oid,
	}

	// render the template
	var renderedBuffer bytes.Buffer
	err = tmpl.Execute(&renderedBuffer, attribs)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("%s\n", renderedBuffer.String())
	return renderedBuffer.Bytes(), nil
}

//
// end of file
//
