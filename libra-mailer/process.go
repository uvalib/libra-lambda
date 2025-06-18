//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
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
	es, err := newEasystoreProxy(cfg)
	if err != nil {
		fmt.Printf("ERROR: creating easystore proxy (%s)\n", err.Error())
		return err
	}

	// important, cleanup properly
	defer es.Close()

	obj, err := getEasystoreObjectByKey(es, ev.Namespace, ev.Identifier, uvaeasystore.Fields+uvaeasystore.Metadata)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	// object fields contain useful state information
	fields := obj.Fields()

	// ensure we have the depositor field else we have no-one to email
	v, f := fields["depositor"]
	if f == false || len(v) == 0 {
		fmt.Printf("ERROR: missing depositor field for object ns/oid [%s/%s]\n", ev.Namespace, ev.Identifier)
		return uvaeasystore.ErrBadParameter
	}

	// field we add to ensure we do not mail more than once
	emailSentFieldName := "unknown"

	var mailType emailType

	// check the event type
	switch ev.EventName {
	case uvalibrabus.EventObjectCreate, uvalibrabus.EventCommandMailInvite:
		mailType = ETD_OPTIONAL_INVITATION
		if fields["source"] == "sis" {
			mailType = ETD_SIS_INVITATION
		}
		emailSentFieldName = "invitation-sent"

	case uvalibrabus.EventWorkPublish, uvalibrabus.EventCommandMailSuccess:
		mailType = ETD_SUBMITTED_AUTHOR
		emailSentFieldName = "submitted-sent"

	default:
		fmt.Printf("INFO: uninteresting event, ignoring\n")
		return nil
	}

	// check to make sure we do not resend an email unless commanded to do so
	if ev.EventName == uvalibrabus.EventObjectCreate || ev.EventName == uvalibrabus.EventWorkPublish {
		if len(fields[emailSentFieldName]) != 0 {
			fmt.Printf("INFO: email already sent, ignoring\n")
			return nil
		}
	}

	// get a new http client and get an auth token
	httpClient := newHttpClient(1, 30)
	token, err := getAuthToken(httpClient, cfg.MintAuthUrl)
	if err != nil {
		return err
	}

	// lookup the user
	depositor, err := getUserDetails(cfg.UserInfoUrl, fields["depositor"], token, httpClient)
	if err != nil {
		return err
	}

	// if we did not find the user...
	if depositor == nil {
		fmt.Printf("ERROR: cannot find user details for [%s]\n", fields["depositor"])
		return err
	}

	// if the user does not have an email
	if len(depositor.Email) == 0 {
		fmt.Printf("ERROR: cannot find email for [%s]\n", fields["depositor"])
		return err
	}

	// render the email body and bail out in the event of an error
	mailSubject, mailBody, err := renderEmailSubjectAndBody(cfg, mailType, depositor, obj)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}

	// send the mail
	err = sendEmail(cfg, mailSubject, depositor.Email, []string{}, mailBody)
	if err != nil {
		return err
	}

	// a special case, we also need to email the registrar
	if ev.EventName == uvalibrabus.EventWorkPublish {
		//mailSubject, mailBody, err = renderEmailSubjectAndBody(cfg, ETD_SUBMITTED_ADVISOR, obj)
		//mailRecipient = //FIXME
		//err = sendEmail(cfg, mailSubject, mailRecipient, []string{}, mailBody)
		//if err != nil {
		//	return err
		//}
	}

	// update the field to note that we have sent the email(s)
	fields[emailSentFieldName] = time.Now().UTC().Format(time.RFC3339)
	obj.SetFields(fields)
	err = putEasystoreObject(es, obj, uvaeasystore.Fields)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}
	return nil
}

//
// end of file
//
