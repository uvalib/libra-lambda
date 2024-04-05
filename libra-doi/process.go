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

func process(messageID string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to librabus event
	ev, err := uvalibrabus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra bus event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("EVENT %s from:%s -> %s\n", messageID, messageSrc, ev.String())

	// initial namespace validation
	if ev.Namespace != libraEtdNamespace && ev.Namespace != libraOpenNamespace {
		fmt.Printf("WARNING: unsupported namespace (%s), ignoring\n", ev.Namespace)
		return nil
	}

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

	// important, cleanup properly
	defer es.Close()

	eso, err := getEasystoreObjectByKey(es, ev.Namespace, ev.Identifier, uvaeasystore.Fields+uvaeasystore.Metadata)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	fields := eso.Fields()

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

	var payload DataciteData
	if ev.Namespace == cfg.ETDNamespace.Name {
		work, err := librametadata.ETDWorkFromBytes(mdBytes)
		if err != nil {
			fmt.Printf("ERROR: unable to process ETD Work %s\n", err.Error())
			return err
		}

		//spew.Dump(work)
		payload = createETDPayload(work, cfg, fields)

	} else if ev.Namespace == cfg.OpenNamespace.Name {
		work, err := librametadata.OAWorkFromBytes(mdBytes)
		if err != nil {
			fmt.Printf("ERROR: unable to process OA Work  %s\n", err.Error())
			return err
		}

		//spew.Dump(work)
		payload = createOAPayload(work, cfg, fields)
	}

	payload.Data.Attributes.URL =
		fmt.Sprintf("%s/public/%s/%s", cfg.PublicURLBase, cfg.OAPublicShoulder, ev.Identifier)

	if len(payload.Data.Attributes.DOI) == 0 &&
		fields["draft"] == "false" {
		// No DOI but the work is published
		// Maybe this should follow the bus event

		payload.Data.Attributes.Event = "publish"

	} // else Datacite creates a draft by default

	spew.Dump(payload)

	cfg.httpClient = *newHttpClient(1, 30)
	// send to Datacite
	doi, err := sendToDatacite(cfg, &payload)
	if err != nil {
		fmt.Printf("ERROR: sending to Datacite (%s)\n", err.Error())
		return err
	}

	// Save the new DOI
	if !strings.HasSuffix(fields["doi"], doi) {
		fmt.Printf("INFO: New DOI for %s is %s\n", ev.Identifier, doi)
		fields["doi"] = fmt.Sprintf("%s/%s", cfg.DOIBaseURL, doi)
		eso.SetFields(fields)
		return putEasystoreObject(es, eso, uvaeasystore.Fields)
	}

	// DOI update complete
	fmt.Printf("INFO: Successfully updated\n")

	return nil
}

//
// end of file
//
