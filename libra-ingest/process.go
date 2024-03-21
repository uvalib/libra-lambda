//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
)

var namespace = "libraetd"

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

	// init the parameter client
	err = initParameter()
	if err != nil {
		return err
	}

	// get our state information
	sisLastProcessed, err := getParameter(cfg.SisIngestStateName)
	if err != nil {
		return err
	}
	optionalLastProcessed, err := getParameter(cfg.OptionalIngestStateName)
	if err != nil {
		return err
	}

	fmt.Printf("INFO: last SIS = [%s]\n", sisLastProcessed)
	fmt.Printf("INFO: last OPT = [%s]\n", optionalLastProcessed)

	// get a new http client and get an auth token
	httpClient := newHttpClient(1, 30)
	token, err := getAuthToken(cfg.MintAuthUrl, httpClient)
	if err != nil {
		return err
	}

	// get inbound SIS items
	//sisList, err := inboundSis(cfg, sisLastProcessed, token, httpClient)

	// get inbound optional items
	optionalList, err := inboundOptional(cfg, optionalLastProcessed, token, httpClient)
	if err != nil {
		return err
	}

	// bail out if nothing to do
	if len(optionalList) == 0 {
		//if len(sisList) == 0 && len(optionalList) == 0 {
		fmt.Printf("INFO: nothing to do, terminating early\n")
		return nil
	}

	// easystore access
	es, err := newEasystore(cfg)
	if err != nil {
		fmt.Printf("ERROR: creating easystore (%s)\n", err.Error())
		return err
	}

	// process inbound SIS items
	//err = processSis(cfg, sisList, es)
	//if err != nil {
	//	return err
	//}

	// get the last one processed and update the state if necessary
	//sisLast := lastSisId(sisList)
	//if sisLastProcessed != sisLast {
	//	fmt.Printf("INFO: last SIS = [%s]\n", sisLast)
	//err = setParameter(cfg.SisIngestStateName, sisLast)
	//if err != nil {
	//	return err
	//}
	//}

	// process inbound optional items
	err = processOptional(cfg, optionalList, es)
	if err != nil {
		return err
	}

	// get the last one processed and update the state if necessary
	optionalLast := lastOptionalId(optionalList)
	if optionalLastProcessed != optionalLast {
		fmt.Printf("INFO: last OPT = [%s]\n", optionalLast)
		err = setParameter(cfg.OptionalIngestStateName, optionalLast)
		if err != nil {
			return err
		}
	}

	return nil
}

//
// end of file
//
