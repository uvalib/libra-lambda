//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"

	"github.com/uvalib/easystore/uvaeasystore"
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

	// easystore access
	esro, err := newEasystoreReadonlyProxy(cfg)
	if err != nil {
		fmt.Printf("ERROR: creating easystore proxy (%s)\n", err.Error())
		return err
	}

	// important, cleanup properly
	defer esro.Close()

	obj, err := getEasystoreObjectByKey(esro, ev.Namespace, ev.Identifier, uvaeasystore.AllComponents)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// get a new http client
	httpClient := newHttpClient(1, 30)
	// important, cleanup properly
	defer httpClient.CloseIdleConnections()

	// write the content to the local filesystem
	bagName, files, err := createBagContent(cfg, httpClient, obj)
	if err != nil {
		fmt.Printf("ERROR: creating bag content for ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// register the incoming submission
	resp, err := registerSubmission(cfg, httpClient)
	if err != nil {
		fmt.Printf("ERROR: registering APTrust submission (%s)\n", err.Error())
		return err
	}

	// init the S3 client
	s3, err := newS3Client()
	if err != nil {
		fmt.Printf("ERROR: creating S3 client (%s)\n", err.Error())
		return err
	}

	// upload to S3
	err = uploadContent(cfg, s3, resp.DepositBucket, resp.DepositPath, bagName, files)
	if err != nil {
		return err
	}

	// initiate the submission
	err = initiateSubmission(cfg, httpClient, resp.SubmissionIdentifier, bagName)
	if err != nil {
		fmt.Printf("ERROR: initiating APTrust submission (%s)\n", err.Error())
		return err
	}

	// log the happy news
	fmt.Printf("INFO: EVENT %s from %s processed OK\n", messageId, messageSrc)
	return nil
}

//
// end of file
//
