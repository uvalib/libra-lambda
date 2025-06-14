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

	// init the parameter client
	ssm, err := newParameterClient()
	if err != nil {
		fmt.Printf("ERROR: creating ssm client (%s)\n", err.Error())
		return err
	}

	// get our state information
	if err != nil {
		return err
	}
	sisLastProcessed, err := getParameter(ssm, cfg.SisIngestStateName)
	if err != nil {
		return err
	}
	fmt.Printf("INFO: last SIS = [%s]\n", sisLastProcessed)

	// get a new http client and get an auth token
	httpClient := newHttpClient(1, 30)
	token, err := getAuthToken(httpClient, cfg.MintAuthUrl)
	if err != nil {
		return err
	}

	// get inbound SIS items
	sisList, err := inboundSis(cfg, sisLastProcessed, token, httpClient)
	if err != nil {
		return err
	}

	// bail out if nothing to do
	if len(sisList) == 0 {
		fmt.Printf("INFO: nothing to do, terminating early\n")
		return nil
	}

	// easystore access
	es, err := newEasystoreProxy(cfg)
	if err != nil {
		fmt.Printf("ERROR: creating easystore proxy (%s)\n", err.Error())
		return err
	}

	// important, cleanup properly
	defer es.Close()

	// process inbound SIS items
	err = processSis(cfg, sisList, es)
	if err != nil {
		return err
	}

	// get the last one processed and update the state if necessary
	if len(sisList) != 0 {
		sisLast := lastSisId(sisList)
		if sisLastProcessed != sisLast {
			fmt.Printf("INFO: last SIS = [%s]\n", sisLast)
			err = setParameter(ssm, cfg.SisIngestStateName, sisLast)
			if err != nil {
				fmt.Printf("ERROR: setting parameter (%s)\n", err.Error())
				return err
			}
		}
	}

	return nil
}

//
// end of file
//
