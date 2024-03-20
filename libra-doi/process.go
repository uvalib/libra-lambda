//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"

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

	obj, err := getEasystoreObject(es, ev.Namespace, ev.Identifier)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// object fields contain useful state information
	fields := obj.Fields()
	//fmt.Printf("%+v\n", fields)

	currentDOI := fields["doi"]

	// Check DOI
	if len(currentDOI) == 0 {
		// No DOI present. Create one.
		fmt.Printf("INFO: DOI blank\n")
	} else {
		// Update DOI
		fmt.Printf("INFO: DOI for %s = %s\n", ev.Identifier, currentDOI)
	}


	return nil
}

//
// end of file
//
