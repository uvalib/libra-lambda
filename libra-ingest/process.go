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

	// init the parameter client
	err = initParameter()
	if err != nil {
		return err
	}

	// get our state information
	sisLast, err := getParameter(cfg.SisIngestStateName)
	if err != nil {
		return err
	}
	optionalLast, err := getParameter(cfg.OptionalIngestStateName)
	if err != nil {
		return err
	}

	fmt.Printf("INFO: latest SIS = [%s]\n", sisLast)
	fmt.Printf("INFO: latest OPT = [%s]\n", optionalLast)

	// get a new http client and get an auth token
	//httpClient := newHttpClient(1, 30)
	//token, err := getAuthToken(cfg.MintAuthUrl, httpClient)
	//if err != nil {
	//	return err
	//}

	//err = inboundSis(cfg, sisLast, token, httpClient)
	//err = inboundOptional(cfg, optionalLast, token, httpClient)

	//err = setParameter(cfg.OptionalIngestStateName, next)
	//if err != nil {
	//	return err
	//}

	return nil
}

//
// end of file
//
