//
//
//

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/uvalib/easystore/uvaeasystore"
	librametadata "github.com/uvalib/libra-metadata"
)

type InboundSisResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Details []InboundSisItem `json:"details"`
}

type InboundSisItem struct {
	InboundId   string `json:"inbound_id"`
	Id          string `json:"id"`
	ComputingId string `json:"computing_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Title       string `json:"title"`
	Department  string `json:"department"`
	Degree      string `json:"degree"`
}

func processSis(cfg *Config, objs []InboundSisItem, es uvaeasystore.EasyStore) error {
	fmt.Printf("INFO: processing %d SIS item(s)\n", len(objs))

	// audit infrastructure
	auditWho := "libra-ingest"
	messageBus, _ := NewEventBus(cfg.BusName, auditWho)

	var returnErr error
	for _, o := range objs {
		fmt.Printf("INFO: processing SIS # %s for %s\n", o.InboundId, o.ComputingId)

		sourceId := fmt.Sprintf("sis:%s", o.Id)
		fields := uvaeasystore.DefaultEasyStoreFields()
		fields["source-id"] = sourceId

		// try and find an existing object
		esrs, err := getEasystoreObjectsByFields(es, libraEtdNamespace, fields, uvaeasystore.Fields+uvaeasystore.Metadata)
		if err != nil {
			fmt.Printf("ERROR: finding easystore object, continuing (%s)\n", err.Error())
			returnErr = err
			continue
		}

		// did we find an existing object?
		if esrs.Count() == 1 {
			eso, err := esrs.Next()
			if err != nil {
				fmt.Printf("ERROR: finding easystore object, continuing (%s)\n", err.Error())
				returnErr = err
				continue
			}

			// ensure the work is in draft state
			if eso.Fields()["draft"] == "true" {

				// if we have a metadata payload
				if eso.Metadata() != nil {
					esomd := eso.Metadata()
					pl, err := esomd.Payload()
					if err != nil {
						fmt.Printf("ERROR: getting metadata from easystore object [%s/%s], continuing (%s)\n", eso.Namespace(), eso.Id(), err.Error())
						returnErr = err
						continue
					}

					md, err := librametadata.ETDWorkFromBytes(pl)
					if err != nil {
						fmt.Printf("ERROR: unmarshaling metadata from easystore object [%s/%s], continuing (%s)\n", eso.Namespace(), eso.Id(), err.Error())
						returnErr = err
						continue
					}

					// is this a title update
					if md.Title != o.Title {
						fmt.Printf("INFO: title update for unpublished work [%s/%s]\n", eso.Namespace(), eso.Id())
						previous := md.Title
						md.Title = o.Title

						eso.SetMetadata(md)
						err = putEasystoreObject(es, eso, uvaeasystore.Metadata)
						if err != nil {
							fmt.Printf("ERROR: updating easystore object [%s/%s], continuing (%s)\n", eso.Namespace(), eso.Id(), err.Error())
							returnErr = err
							continue
						}
						// audit this change
						_ = pubAuditEvent(messageBus, eso, auditWho, "title", previous, o.Title)
					}
				} else {
					fmt.Printf("ERROR: sis update but work has missing metadata [%s/%s], ignoring\n", eso.Namespace(), eso.Id())
				}
			} else {
				fmt.Printf("WARNING: sis update for published work [%s/%s], ignoring\n", eso.Namespace(), eso.Id())
			}
		} else {
			// we did not find an existing one, create a new easystore object
			eso := uvaeasystore.NewEasyStoreObject(libraEtdNamespace, "")

			// add some fields
			fields["author"] = o.ComputingId
			fields["depositor"] = o.ComputingId
			fields["create-date"] = time.Now().UTC().Format(time.RFC3339)
			fields["source-id"] = sourceId
			fields["source"] = "sis"
			fields["draft"] = "true"
			eso.SetFields(fields)

			meta := librametadata.ETDWork{}
			meta.Program = o.Department
			meta.Degree = o.Degree
			meta.Title = o.Title
			meta.Author = librametadata.ContributorData{
				ComputeID:   o.ComputingId,
				FirstName:   o.FirstName,
				LastName:    o.LastName,
				Department:  o.Department,
				Institution: "University of Virginia",
			}

			// An ETDWork does not serialize the same way as an EasyStoreMetadata object
			// does when being managed by json.Marshal/json.Unmarshal so we wrap it in an object that
			// behaves appropriately
			pl, err := meta.Payload()
			if err != nil {
				log.Printf("ERROR: serializing ETDWork: %s, continuing", err.Error())
				returnErr = err
				continue
			}
			eso.SetMetadata(uvaeasystore.NewEasyStoreMetadata(meta.MimeType(), pl))

			// create the new object
			err = createEasystoreObject(es, eso)
			if err != nil {
				fmt.Printf("ERROR: creating easystore object, continuing (%s)\n", err.Error())
				returnErr = err
				continue
			}

			// audit this set of changes
			_ = pubAuditEvent(messageBus, eso, auditWho, "create-date", "", fields["create-date"])
			_ = pubAuditEvent(messageBus, eso, auditWho, "program", "", meta.Program)
			_ = pubAuditEvent(messageBus, eso, auditWho, "degree", "", meta.Degree)
			_ = pubAuditEvent(messageBus, eso, auditWho, "title", "", meta.Title)
			_ = pubAuditEvent(messageBus, eso, auditWho, "author.cid", "", meta.Author.ComputeID)
			_ = pubAuditEvent(messageBus, eso, auditWho, "author.firstname", "", meta.Author.FirstName)
			_ = pubAuditEvent(messageBus, eso, auditWho, "author.lastname", "", meta.Author.LastName)
			_ = pubAuditEvent(messageBus, eso, auditWho, "author.department", "", meta.Author.Department)
			_ = pubAuditEvent(messageBus, eso, auditWho, "author.institution", "", meta.Author.Institution)
		}
	}

	return returnErr
}

func inboundSis(config *Config, last string, auth string, client *http.Client) ([]InboundSisItem, error) {

	// substitute values into url
	url := strings.Replace(config.SisIngestUrl, "{:last}", last, 1)
	url = strings.Replace(url, "{:auth}", auth, 1)

	payload, err := httpGet(client, url)
	if err != nil {
		// special case of no items
		if strings.Contains(err.Error(), "HTTP 404") == true {
			return make([]InboundSisItem, 0), nil
		}
		return nil, err
	}

	resp := InboundSisResponse{}
	err = json.Unmarshal(payload, &resp)
	if err != nil {
		fmt.Printf("ERROR: json unmarshal of InboundSisResponse (%s)\n", err.Error())
		return nil, err
	}

	fmt.Printf("INFO: received %d SIS item(s)\n", len(resp.Details))
	return resp.Details, nil
}

func lastSisId(objs []InboundSisItem) string {
	last := "0"
	for _, in := range objs {
		l, _ := strconv.Atoi(last)
		c, _ := strconv.Atoi(in.InboundId)
		if c > l {
			last = in.InboundId
		}
	}
	return last
}

//
// end of file
//
