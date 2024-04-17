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

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to librabus event
	ev, err := uvalibrabus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra bus event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("EVENT %s from:%s -> %s\n", messageId, messageSrc, ev.String())

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

	// init the S3 client
	s3, err := newS3Client()
	if err != nil {
		fmt.Printf("ERROR: creating S3 client (%s)\n", err.Error())
		return err
	}

	// render the document
	buf, err := docRender(ev.Namespace, ev.Identifier)
	if err != nil {
		fmt.Printf("ERROR: rendering template (%s)\n", err.Error())
		return err
	}

	// populate the key template
	year := fmt.Sprintf("%04d", time.Now().Year())
	bucketKey := strings.Replace(cfg.BucketKeyTemplate, "{:year}", year, 1)
	bucketKey = strings.Replace(bucketKey, "{:namespace}", ev.Namespace, 1)
	bucketKey = strings.Replace(bucketKey, "{:id}", ev.Identifier, 1)

	// upload to S3
	err = putS3(s3, cfg.BucketName, bucketKey, buf)
	if err != nil {
		fmt.Printf("ERROR: uploading (%s)\n", err.Error())
		return err
	}

	return nil
}

//
// end of file
//
