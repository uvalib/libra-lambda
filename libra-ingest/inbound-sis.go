//
//
//

package main

import (
	"encoding/json"
	"fmt"
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
						md.Title = o.Title

						eso.SetMetadata(md)
						err = putEasystoreObject(es, eso, uvaeasystore.Metadata)
						if err != nil {
							fmt.Printf("ERROR: updating easystore object [%s/%s], continuing (%s)\n", eso.Namespace(), eso.Id(), err.Error())
							returnErr = err
							continue
						}
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
			fields["depositor"] = fmt.Sprintf("%s@virginia.edu", o.ComputingId)
			fields["create-date"] = time.Now().Format(time.RFC3339)
			fields["source-id"] = sourceId
			fields["source"] = "sis"
			eso.SetFields(fields)

			meta := librametadata.ETDWork{}
			meta.Department = o.Department
			meta.Degree = o.Degree
			meta.Title = o.Title
			meta.Author = librametadata.ContributorData{
				ComputeID:   o.ComputingId,
				FirstName:   o.FirstName,
				LastName:    o.LastName,
				Program:     o.Department,
				Institution: "University of Virginia",
			}
			eso.SetMetadata(meta)

			// create the new object
			err := createEasystoreObject(es, eso)
			if err != nil {
				fmt.Printf("ERROR: creating easystore object, continuing (%s)\n", err.Error())
				returnErr = err
				continue
			}
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
