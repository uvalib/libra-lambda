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
	cfg.AuthToken, err = getAuthToken(cfg.httpClient, cfg.MintAuthURL)
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

	var payload DataciteData

	// If there is no DOI and the work is published, set the event to publish
	if len(fields["doi"]) == 0 &&
		fields["draft"] == "false" {
		// No DOI and the work is published
		payload.Data.Attributes.Event = "publish"

	} else if len(fields["doi"]) == 0 && fields["draft"] == "true" &&
		ev.Namespace != cfg.ETDNamespace.Name {
		// Quit if draft with no DOI and not an ETD
		fmt.Printf("INFO: Skipping draft work %s \n", ev.Identifier)
		return nil

	} else if len(fields["doi"]) > 0 && fields["draft"] == "true" {
		// If the work has a DOI but is a draft, set the event to hide
		payload.Data.Attributes.Event = "hide"
	} // A draft is created when event is blank

	work, err := librametadata.ETDWorkFromBytes(mdBytes)
	if err != nil {
		fmt.Printf("ERROR: unable to process ETD Work %s\n", err.Error())
		return err
	}

	payload = createETDPayload(work, fields)
	payload.Data.Attributes.URL =
		fmt.Sprintf("%s/public/%s/%s", cfg.PublicURLBase, cfg.ETDPublicShoulder, ev.Identifier)

	spew.Dump(payload)

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
		return nil
	}

	// DOI update complete
	fmt.Printf("INFO: Successfully updated\n")

	return nil
}

//
// end of file
//
