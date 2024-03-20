//
//
//

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func inboundSis(config *Config, last string, auth string, client *http.Client) error {

	// substitute values into url
	url := strings.Replace(config.SisIngestUrl, "{:last}", last, 1)
	url = strings.Replace(url, "{:auth}", auth, 1)

	payload, err := httpGet(url, client)
	if err != nil {
		return err
	}

	resp := InboundSisResponse{}
	err = json.Unmarshal(payload, &resp)
	if err != nil {
		fmt.Printf("ERROR: json unmarshal of InboundSisResponse (%s)\n", err.Error())
		return err
	}

	fmt.Printf("received %d item(s)\n", len(resp.Details))
	return nil
}

//
// end of file
//
