//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
	"log"
)

// we treat libraETD and libraOpen events differently
var libraEtdNamespace = "libraetd"
var libraOpenNamespace = "libraopen"

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	fmt.Printf("EVENT %s from:%s -> %s\n", messageId, messageSrc, string(rawMsg))

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	busCfg := uvalibrabus.UvaBusConfig{
		Source:  cfg.SourceName,
		BusName: cfg.BusName,
	}

	// create message bus client
	bus, err := uvalibrabus.NewUvaBus(busCfg)
	if err != nil {
		log.Printf("ERROR: creating event bus client (%s)\n", err.Error())
		return err
	}
	fmt.Printf("Using: %s@%s\n", cfg.SourceName, cfg.BusName)

	// create event
	ev := uvalibrabus.UvaBusEvent{}
	ev.EventName = uvalibrabus.EventScheduleEtdIngest
	ev.Identifier = "none"

	// publish ETD namespace event
	ev.Namespace = libraEtdNamespace
	err = bus.PublishEvent(ev)
	if err != nil {
		log.Printf("ERROR: publishing event (%s)\n", err.Error())
		return err
	}

	// publish Open namespace event
	ev.Namespace = libraOpenNamespace
	err = bus.PublishEvent(ev)
	if err != nil {
		log.Printf("ERROR: publishing event (%s)\n", err.Error())
		return err
	}

	return nil
}

//
// end of file
//
