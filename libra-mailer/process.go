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

// field name indicating email already sent
var emailSentFieldName = "email-sent"

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

	// easystore access
	es, err := newEasystore(cfg)
	if err != nil {
		fmt.Printf("ERROR: creating easystore (%s)\n", err.Error())
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

	// have we already sent the email
	if len(fields[emailSentFieldName]) != 0 {
		fmt.Printf("INFO: email already sent, ignoring\n")
		return nil
	}

	// mail attributes
	var mailSubject string
	var mailBody string

	// check the event type
	switch ev.EventName {
	case uvalibrabus.EventObjectCreate:
		// we send notifications for libraetd events only
		switch obj.Namespace() {
		case libraEtdNamespace:
			mail := ETD_OPTIONAL_INVITATION
			if fields["source"] == "sis" {
				mail = ETD_SIS_INVITATION
			}
			mailSubject, mailBody, err = emailSubjectAndBody(cfg, mail, obj)

		case libraOpenNamespace:
			fmt.Printf("INFO: uninteresting namespace for event, ignoring\n")
			return nil

		default:
			err = fmt.Errorf("unsupported namespace")
		}

	case uvalibrabus.EventWorkPublish:
		switch obj.Namespace() {
		case libraEtdNamespace:
			// FIXME: support advisor email too

			mailSubject, mailBody, err = emailSubjectAndBody(cfg, ETD_SUBMITTED_AUTHOR, obj)

		case libraOpenNamespace:

			mailSubject, mailBody, err = emailSubjectAndBody(cfg, OPEN_SUBMITTED_AUTHOR, obj)

		default:

			err = fmt.Errorf("unsupported namespace for work publish event")
		}

	default:
		fmt.Printf("INFO: uninteresting event, ignoring\n")
		return nil
	}

	// bail out in the event of an error
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return err
	}

	// send the mail
	mailRecipient := fmt.Sprintf("%s@virginia.edu", fields["depositor"])
	err = sendEmail(cfg, mailSubject, mailRecipient, []string{}, mailBody)
	if err != nil {
		return err
	}

	// a special case, we also need to email the registrar
	if ev.EventName == uvalibrabus.EventWorkPublish && obj.Namespace() == libraEtdNamespace {
		mailSubject, mailBody, err = emailSubjectAndBody(cfg, ETD_SUBMITTED_ADVISOR, obj)
		//mailRecipient = //FIXME
		//err = sendEmail(cfg, mailSubject, mailRecipient, []string{}, mailBody)
		//if err != nil {
		//	return err
		//}
	}

	// update the field to note that we have sent the email(s)
	fields[emailSentFieldName] = time.DateTime
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
