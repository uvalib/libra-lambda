//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
	"strings"
	"time"
)

// we treat libraETD and libraOpen events differently
var libraEtdNamespace = "libraetd"
var libraOpenNamespace = "libraopen"

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

	// init the S3 client
	err = initS3()
	if err != nil {
		fmt.Printf("ERROR: creating S3 client (%s)\n", err.Error())
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

	obj, err := getEasystoreObject(es, ev.Namespace, ev.Identifier)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// render the document
	buf, err := docRender(cfg, obj)
	if err != nil {
		return err
	}

	// populate the key template
	year := fmt.Sprintf("%04d", time.Now().Year())
	bucketKey := strings.Replace(cfg.BucketKeyTemplate, "{:id}", obj.Id(), 1)
	bucketKey = strings.Replace(bucketKey, "{:year}", year, 1)

	// upload to S3
	err = putS3(cfg.BucketName, bucketKey, buf)
	if err != nil {
		fmt.Printf("ERROR: uploading (%s)\n", err.Error())
		return err
	}

	return nil
}

//
// end of file
//
