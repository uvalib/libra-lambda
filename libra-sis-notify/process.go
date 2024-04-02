//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
	"time"
)

// field name indicating sis notified
var sisNotifiedFieldName = "sis-sent"

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to librabus event
	ev, err := uvalibrabus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra bus event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("EVENT %s from:%s -> %s\n", messageId, messageSrc, ev.String())

	// make sure this is an object we are interested in
	if ev.Namespace != libraEtdNamespace {
		fmt.Printf("INFO: uninteresting event, ignoring\n")
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

	obj, err := getEasystoreObjectByKey(es, ev.Namespace, ev.Identifier, uvaeasystore.Fields)
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

	// notify SIS of the activity
	fields := obj.Fields()
	err = notifySis(cfg, fields, token, httpClient)
	if err != nil {
		fmt.Printf("ERROR: notifying SIS (%s)\n", err.Error())
		return err
	}

	// update the field to note that we have notified SIS
	fields[sisNotifiedFieldName] = time.DateTime
	obj.SetFields(fields)
	return putEasystoreObject(es, obj, uvaeasystore.Fields)
}

//
// end of file
//
