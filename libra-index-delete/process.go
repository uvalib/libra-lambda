//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/uvalib/librabus-sdk/uvalibrabus"
)

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

	// get a new http client
	httpClient := newHttpClient(1, 30)
	// important, cleanup properly
	defer httpClient.CloseIdleConnections()

	url := strings.Replace(cfg.IndexDeleteUrl, "{:id}", ev.Identifier, 1)

	_, err = httpDelete(httpClient, url)
	if err == nil {
		// log the happy news
		fmt.Printf("INFO: EVENT %s from %s processed OK\n", messageId, messageSrc)
	} else {
		// log the sad news
		fmt.Printf("ERROR: EVENT %s from %s FAILED (%s)\n", messageId, messageSrc, err.Error())
	}

	return err
}

//
// end of file
//
