//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/uvalib/easystore/uvaeasystore"
	librametadata "github.com/uvalib/libra-metadata"
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

	cfg.httpClient = newHttpClient(1, 30)
	// important, cleanup properly
	defer cfg.httpClient.CloseIdleConnections()

	cfg.AuthToken, err = getAuthToken(cfg.httpClient, cfg.MintAuthURL)
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

	eso, err := getEasystoreObjectByKey(es, ev.Namespace, ev.Identifier, uvaeasystore.Fields+uvaeasystore.Metadata)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	fields := eso.Fields()

	if len(fields["doi"]) > 0 && strings.HasPrefix(fields["doi"], Cfg().DOIBaseURL) == false {
		// doi exists but for the wrong environment.
		// if a production DOI is sent to test
		fmt.Printf("WARNING: DOI %s has the wrong Datacite hostname for this environment (%s). ", fields["doi"], Cfg().DOIBaseURL)
		return nil
	}

	if eso.Metadata() == nil {
		fmt.Printf("ERROR: unable to get metadata payload for ns/oid [%s/%s]\n", ev.Namespace, ev.Identifier)
		return ErrNoMetadata
	}

	mdBytes, err := eso.Metadata().Payload()
	if err != nil {
		fmt.Printf("ERROR: unable to get metadata payload from response: %s\n", err.Error())
		return err
	}

	spew.Dump(fields)

	var eventType string
	switch ev.EventName {
	// https://support.datacite.org/docs/doi-states
	case uvalibrabus.EventWorkPublish:
		if fields["draft"] == "false"  {
			// Publish Event for published work
			fmt.Printf("INFO: Publish Event for [%s/%s]\n", ev.Namespace, ev.Identifier)
			if len(fields["doi"]) > 0 {
				fmt.Printf("INFO: Publishing DOI %s \n", fields["doi"])
			} else {
				fmt.Printf("INFO: Publishing new DOI\n")
			}
			eventType = "publish"

		} else{
			fmt.Printf("ERROR: Can't publish a draft %s \n", ev.Identifier)
			return nil
		}

	case uvalibrabus.EventWorkUnpublish:
		fmt.Printf("INFO: Unpublish Event for [%s/%s]\n", ev.Namespace, ev.Identifier)
		// "registered" is a reserved DOI but not findable
		eventType = "register"

	case uvalibrabus.EventObjectCreate:
	  if len(fields["doi"]) == 0 && fields["draft"] == "true" {
			fmt.Printf("INFO: Registering new DOI for draft work [%s/%s].\n", ev.Namespace, ev.Identifier)
			eventType = "register"
		}

	case uvalibrabus.EventMetadataUpdate, uvalibrabus.EventCommandDoiSync :
		// No Event change for edits or resyncs
		if len(fields["doi"]) > 0 {
			fmt.Printf("INFO: Update Event for [%s/%s] with DOI %s\n", ev.Namespace, ev.Identifier, fields["doi"])
		} else {
			fmt.Printf("INFO: Update Event for [%s/%s] without DOI. One will be created.\n", ev.Namespace, ev.Identifier)
		}
	}


	work, err := librametadata.ETDWorkFromBytes(mdBytes)
	if err != nil {
		fmt.Printf("ERROR: unable to process ETD Work %s\n", err.Error())
		return err
	}

	payload := createETDPayload(work, fields)
	payload.Data.Attributes.Event = eventType
	payload.Data.Attributes.URL =
		fmt.Sprintf("%s/%s/%s", cfg.PublicURLBase, cfg.ETDPublicShoulder, ev.Identifier)

	// send to Datacite
	doi, err := sendToDatacite(&payload)
	if err != nil {
		fmt.Printf("ERROR: sending to Datacite (%s)\n", err.Error())
		return err
	}

	// Save the new DOI
	if !strings.HasSuffix(fields["doi"], doi) {
		fmt.Printf("INFO: New DOI for [%s/%s] is %s\n", ev.Namespace, ev.Identifier, doi)

		// Refresh easystore object
		eso, err = getEasystoreObjectByKey(es, ev.Namespace, ev.Identifier, uvaeasystore.Fields+uvaeasystore.Metadata)
		if err != nil {
			fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
			fmt.Printf("ERROR: DOI created but not saved for [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, doi)
			return err
		}
		fields = eso.Fields()

		fields["doi"] = fmt.Sprintf("%s/%s", cfg.DOIBaseURL, doi)
		eso.SetFields(fields)
		err = putEasystoreObject(es, eso, uvaeasystore.Fields)
		if err != nil {
			fmt.Printf("ERROR: unable to update object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
			fmt.Printf("ERROR: DOI created but not saved for [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, doi)
			return err
		}

		// audit this change
		who := "libra-doi"
		bus, _ := NewEventBus(cfg.BusName, who)
		_ = pubAuditEvent(bus, eso, who, "doi", "", fields["doi"])
	}

	// log the happy news
	fmt.Printf("INFO: EVENT %s from %s processed OK\n", messageId, messageSrc)
	return nil
}

//
// end of file
//
