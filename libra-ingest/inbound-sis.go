//
//
//

package main

import (
	"encoding/json"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"net/http"
	"strconv"
	"strings"
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

	for _, o := range objs {
		fmt.Printf("INFO: processing SIS #%s for %s/%s (%s)\n", o.Id, o.FirstName, o.LastName, o.ComputingId)
	}

	return nil
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
