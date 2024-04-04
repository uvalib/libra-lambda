//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"

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

	fmt.Printf("Fields: %+v\n", fields)

	var payload DataciteData
	if ev.Namespace == cfg.ETDNamespace.Name {
		work, err := librametadata.ETDWorkFromBytes(mdBytes)
		if err != nil {
			fmt.Printf("ERROR: unable to process ETD Work %s\n", err.Error())
			return err
		}

		fmt.Printf("Metadata: %+v\n", work)
		payload = createETDPayload(work, cfg, fields)

	} else if ev.Namespace == cfg.OpenNamespace.Name {
		work, err := librametadata.OAWorkFromBytes(mdBytes)
		if err != nil {
			fmt.Printf("ERROR: unable to process OA Work  %s\n", err.Error())
			return err
		}

		fmt.Printf("Metadata: %+v\n", work)
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
	cfg.httpClient = *newHttpClient(1, 30)
	fmt.Printf("Payload: %+v\n", payload)

	//      if work.description.present?
	//        attributes['descriptions'] = [{
	//          description: work.description,
	//          descriptionType: 'Abstract'
	//        }]
	//      end
	//
	//      attributes[:creators] = authors_construct( work )
	//      attributes[:contributors] = contributors_construct( work )
	//      attributes[:subjects] = work.keyword.map{|k| {subject: k}} if work.keyword.present?
	//      attributes[:rightsList] = [{rights: work.rights.first}] if work.rights.present? && work.rights.first.present?
	//      attributes[:fundingReferences] = work.sponsoring_agency.map{|f| {funderName: f}} if work.sponsoring_agency.present?
	//      attributes[:types] = {resourceTypeGeneral: DC_GENERAL_TYPE_TEXT, resourceType: RESOURCE_TYPE_DISSERTATION}
	//
	//      yyyymmdd = extract_yyyymmdd_from_datestring( work.date_published )
	//      yyyymmdd = extract_yyyymmdd_from_datestring( work.date_created ) if yyyymmdd.blank?
	//      attributes[:dates] = [{date: yyyymmdd, dateType: 'Issued'}] if yyyymmdd.present?
	//      attributes[:publicationYear] = yyyymmdd.first(4) if yyyymmdd.present?
	//
	//      attributes[:url] = fully_qualified_work_url( work.id ) # 'http://google.com'
	//      attributes[:publisher] = work.publisher if work.publisher.present?
	//
	//      #puts "==> #{h.to_json}"
	//      payload = {
	//        data: {
	//          type: 'dois',
	//          attributes: attributes
	//        }
	//      }

	// Check DOI
	if len(fields["doi"]) == 0 {
		// No DOI present. Create one.
		fmt.Printf("INFO: DOI blank\n")
		doi, err := createDOI(cfg, eso)
		if err != nil {
			panic(err)
		}

		fmt.Println("New DOI: " + doi)

	} else {
		// Update DOI
		fmt.Printf("INFO: DOI for %s = %s\n", ev.Identifier, fields["doi"])
	}

	return nil
}

//
// end of file
//
