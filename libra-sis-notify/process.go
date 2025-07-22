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

	fmt.Printf("INFO: EVENT %s from %s -> %s\n", messageId, messageSrc, ev.String())

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

	obj, err := getEasystoreObjectByKey(es, ev.Namespace, ev.Identifier, uvaeasystore.Fields)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	fields := obj.Fields()

	// we have already notified SIS, bail out unless this is a command event
	if len(fields[sisNotifiedFieldName]) != 0 && ev.EventName != uvalibrabus.EventCommandSisNotify {
		fmt.Printf("INFO: SIS already notified, ignoring\n")
		return nil
	}

	if strings.HasPrefix(fields["source-id"], "sis:") == true {

		// get a new http client and get an auth token
		httpClient := newHttpClient(1, 30)
		// important, cleanup properly
		defer httpClient.CloseIdleConnections()

		token, err := getAuthToken(httpClient, cfg.MintAuthUrl)
		if err != nil {
			return err
		}

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
		src := "<empty>"
		if len(fields["source-id"]) != 0 {
			src = fields["source-id"]
		}
		fmt.Printf("INFO: not a SIS work (source %s), ignoring\n", src)
		return nil
	}

	// audit this change
	who := "libra-sis-notify"
	bus, _ := NewEventBus(cfg.BusName, who)
	_ = pubAuditEvent(bus, obj, who, sisNotifiedFieldName, "", fields[sisNotifiedFieldName])

	// log the happy news
	fmt.Printf("INFO: EVENT %s from %s processed OK\n", messageId, messageSrc)
	return nil
}

//
// end of file
//
