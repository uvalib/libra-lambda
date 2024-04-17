//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
	"strings"
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

	// initial namespace validation
	if ev.Namespace != libraEtdNamespace && ev.Namespace != libraOpenNamespace {
		fmt.Printf("WARNING: unsupported namespace (%s), ignoring\n", ev.Namespace)
		return nil
	}

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

	fields := obj.Fields()
	if strings.HasPrefix(fields["source-id"], "sis:") == true {
		// notify SIS of the activity
		err = notifySis(cfg, fields, token, httpClient)
		if err != nil {
			fmt.Printf("ERROR: notifying SIS (%s)\n", err.Error())
			return err
		}

		// update the field to note that we have notified SIS
		fields[sisNotifiedFieldName] = time.DateTime
		obj.SetFields(fields)
		err = putEasystoreObject(es, obj, uvaeasystore.Fields)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err.Error())
			return err
		}
	} else {
		fmt.Printf("INFO: not a SIS work (source %s), ignoring\n", fields["source-id"])
		return nil
	}

	return nil
}

//
// end of file
//
