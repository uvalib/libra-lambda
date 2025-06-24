//
// main message processing
//

package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
)

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to librabus event
	ev, err := uvalibrabus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra bus event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("EVENT %s from:%s -> %s\n", messageId, messageSrc, ev.String())

	// initial namespace validation
	if ev.Namespace != libraEtdNamespace {
		fmt.Printf("WARNING: unsupported namespace (%s), ignoring\n", ev.Namespace)
		return nil
	}

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	// easystore access
	es, err := newEasystoreProxy(cfg)
	if err != nil {
		fmt.Printf("ERROR: creating easystore proxy (%s)\n", err.Error())
		return err
	}

	// important, cleanup properly
	defer es.Close()

	eso, err := getEasystoreObjectByKey(es, ev.Namespace, ev.Identifier, uvaeasystore.Fields+uvaeasystore.Metadata)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// get author details
	authorId, err := getWorkAuthor(eso)
	if err != nil {
		fmt.Printf("ERROR: cannot locate author ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	if len(authorId) == 0 {
		fmt.Printf("WARNING: cannot locate author ns/oid [%s/%s]\n", ev.Namespace, ev.Identifier)
		// noting to do
		return nil
	}

	// get a new http client and get an auth token
	httpClient := newHttpClient(1, 30)
	// important, cleanup properly
	defer httpClient.CloseIdleConnections()

	token, err := getAuthToken(httpClient, cfg.MintAuthUrl)
	if err != nil {
		return err
	}

	// attempt to get ORCID for the author
	orcid, err := getOrcidDetails(cfg.OrcidGetDetailsUrl, authorId, token, httpClient)
	if err != nil {
		fmt.Printf("ERROR: getting %s ORCID details ns/oid [%s/%s] (%s)\n", authorId, ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// if author does not have ORCID details, it's all over
	if len(orcid) == 0 {
		fmt.Printf("INFO: no ORCID details for %s ns/oid [%s/%s]\n", authorId, ev.Namespace, ev.Identifier)
		return nil
	}

	// use an existing ORCID update code (if one exists)
	fields := eso.Fields()
	updateCode := fields["orcid-update-code"]

	newCode, err := updateAuthorOrcidActivity(cfg, eso, authorId, updateCode, token, httpClient)
	if err != nil {
		if errors.Is(err, ErrIncompleteData) == true {
			fmt.Printf("WARNING: incomplete data for ORCID activity ns/oid [%s/%s]\n", ev.Namespace, ev.Identifier)
			return nil
		} else {
			fmt.Printf("ERROR: updating %s ORCID activity ns/oid [%s/%s] (%s)\n", authorId, ev.Namespace, ev.Identifier, err.Error())
			return err
		}
	}

	// do we have a new update code
	if len(newCode) != 0 && newCode != updateCode {
		fields["orcid-update-code"] = newCode
		eso.SetFields(fields)
		err = putEasystoreObject(es, eso, uvaeasystore.Fields)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return err
		}
	}

	return nil
}

//
// end of file
//
