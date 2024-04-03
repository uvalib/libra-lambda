//
// main message processing
//

package main

import (
	"encoding/json"
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
	if ev.Namespace != libraEtdNamespace && ev.Namespace != libraOpenNamespace {
		fmt.Printf("WARNING: unsupported namespace (%s), ignoring\n", ev.Namespace)
		return nil
	}

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	// easystore access
	es, err := newEasystore(cfg)
	if err != nil {
		fmt.Printf("ERROR: creating easystore (%s)\n", err.Error())
		return err
	}

	// important, cleanup properly
	defer es.Close()

	eso, err := getEasystoreObjectByKey(es, ev.Namespace, ev.Identifier, uvaeasystore.Fields+uvaeasystore.Metadata)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// get a new http client and get an auth token
	httpClient := newHttpClient(1, 30)
	token, err := getAuthToken(httpClient, cfg.MintAuthUrl)
	if err != nil {
		return err
	}

	// determine if we have ORCID details for the author
	orcidDetails, err := getAuthorOrcidDetails(cfg, eso, token, httpClient)
	if err != nil {
		fmt.Printf("ERROR: getting author ORCID details ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// if author does not have ORCID details, it's all over
	if orcidDetails == nil {
		fmt.Printf("INFO: no ORCID details for author of ns/oid [%s/%s]\n", ev.Namespace, ev.Identifier)
		return nil
	}

	// use an existing ORCID update code (if one exists)
	fields := eso.Fields()
	updateCode := fields["orcid-update-code"]

	newCode, err := updateAuthorOrcidActivity(cfg, eso, orcidDetails.Cid, updateCode, token, httpClient)
	if err != nil {
		fmt.Printf("ERROR: updating author ORCID activity ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// do we have a new update code
	if newCode != updateCode {
		fields["orcid-update-code"] = newCode
		eso.SetFields(fields)
		return putEasystoreObject(es, eso, uvaeasystore.Fields)
	}

	return nil
}

//
// end of file
//
