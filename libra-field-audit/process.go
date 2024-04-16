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

	audit, err := uvalibrabus.MakeAuditEvent(ev.Detail)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra audit event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("AUDIT %s changed %s from:%s to:%s\n", audit.Who, audit.FieldName, audit.Before, audit.After)

	return nil
}

//
// end of file
//
