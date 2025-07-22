//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
	//"strings"
	//"time"
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

	// easystore access
	esro, err := newEasystoreReadonlyProxy(cfg)
	if err != nil {
		fmt.Printf("ERROR: creating easystore proxy (%s)\n", err.Error())
		return err
	}

	// important, cleanup properly
	defer esro.Close()

	obj, err := getEasystoreObjectByKey(esro, ev.Namespace, ev.Identifier, uvaeasystore.Metadata+uvaeasystore.Fields)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// get a new http client
	httpClient := newHttpClient(1, 30)
	// important, cleanup properly
	defer httpClient.CloseIdleConnections()

	// and update the index
	err = updateIndex(cfg, obj, httpClient)
	if err == nil {
		// log the happy news
		fmt.Printf("INFO: EVENT %s from %s processed OK\n", messageId, messageSrc)
	}

	return err
}

//
// end of file
//
