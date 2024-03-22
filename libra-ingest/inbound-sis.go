//
//
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"github.com/uvalib/libra-metadata"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	fmt.Printf("processing %d SIS item(s)\n", len(objs))

	var returnErr error
	for _, o := range objs {
		fmt.Printf("INFO: processing SIS #%s for %s/%s (%s)\n", o.InboundId, o.FirstName, o.LastName, o.ComputingId)

		// FIXME, match against existing SIS entry

		// new easystore object
		eso := uvaeasystore.NewEasyStoreObject(namespace, "")

		// add some fields
		fields := uvaeasystore.DefaultEasyStoreFields()
		fields["author"] = o.ComputingId
		fields["depositor"] = o.ComputingId
		fields["create-date"] = time.Now().Format(time.RFC3339)
		fields["source-id"] = fmt.Sprintf("sis:%s", o.Id)
		fields["source"] = "sis"
		eso.SetFields(fields)

		meta := librametadata.ETDWork{}
		meta.Degree = o.Degree
		meta.Title = o.Title
		meta.Author = librametadata.StudentData{
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

	return returnErr
}

func inboundSis(config *Config, last string, auth string, client *http.Client) ([]InboundSisItem, error) {

	// substitute values into url
	url := strings.Replace(config.SisIngestUrl, "{:last}", last, 1)
	url = strings.Replace(url, "{:auth}", auth, 1)

	payload, err := httpGet(url, client)
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

	fmt.Printf("received %d SIS item(s)\n", len(resp.Details))
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
