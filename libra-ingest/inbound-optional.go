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

type InboundOptionalResponse struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Details []InboundOptionalItem `json:"details"`
}

type InboundOptionalItem struct {
	Id         string `json:"id"`
	Requester  string `json:"requester"`
	For        string `json:"for"`
	Department string `json:"department"`
	Degree     string `json:"degree"`
}

func processOptional(cfg *Config, objs []InboundOptionalItem, es uvaeasystore.EasyStore) error {
	fmt.Printf("processing %d optional item(s)\n", len(objs))

	var returnErr error
	for _, o := range objs {
		fmt.Printf("INFO: processing optional #%s for %s\n", o.Id, o.For)

		// new easystore object
		eso := uvaeasystore.NewEasyStoreObject(libraEtdNamespace, "")

		// add some fields
		fields := uvaeasystore.DefaultEasyStoreFields()
		fields["author"] = o.For
		fields["depositor"] = fmt.Sprintf("%s@virginia.edu", o.For)
		fields["create-date"] = time.Now().Format(time.RFC3339)
		fields["source-id"] = fmt.Sprintf("optional:%s", o.Id)
		fields["source"] = "optional"
		eso.SetFields(fields)

		// add the metadata
		meta := librametadata.ETDWork{}
		meta.Department = o.Department
		meta.Degree = o.Degree
		meta.Author.ComputeID = o.For
		eso.SetMetadata(meta)

		// create the new object
		err := createEasystoreObject(es, eso)
		if err != nil {
			fmt.Printf("ERROR: creating easystore object, continuing (%s)\n", err.Error())
			returnErr = err
		}
	}

	return returnErr
}

func inboundOptional(config *Config, last string, auth string, client *http.Client) ([]InboundOptionalItem, error) {

	// substitute values into url
	url := strings.Replace(config.OptionalIngestUrl, "{:last}", last, 1)
	url = strings.Replace(url, "{:auth}", auth, 1)

	payload, err := httpGet(client, url)
	if err != nil {
		// special case of no items
		if strings.Contains(err.Error(), "HTTP 404") == true {
			return make([]InboundOptionalItem, 0), nil
		}
		return nil, err
	}

	resp := InboundOptionalResponse{}
	err = json.Unmarshal(payload, &resp)
	if err != nil {
		fmt.Printf("ERROR: json unmarshal of InboundOptionalResponse (%s)\n", err.Error())
		return nil, err
	}

	fmt.Printf("received %d optional item(s)\n", len(resp.Details))
	return resp.Details, nil
}

func lastOptionalId(objs []InboundOptionalItem) string {
	last := "0"
	for _, in := range objs {
		l, _ := strconv.Atoi(last)
		c, _ := strconv.Atoi(in.Id)
		if c > l {
			last = in.Id
		}
	}
	return last
}

//
// end of file
//
