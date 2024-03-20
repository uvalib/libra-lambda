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

type InboundOptionalResponse struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Details []InboundSisItem `json:"details"`
}

type InboundOptionalItem struct {
	Id         string `json:"id"`
	Requester  string `json:"requester"`
	For        string `json:"for"`
	Department string `json:"department"`
	Degree     string `json:"degree"`
}

func inboundOptional(config *Config, last string, auth string, client *http.Client) error {

	// substitute values into url
	url := strings.Replace(config.OptionalIngestUrl, "{:last}", last, 1)
	url = strings.Replace(url, "{:auth}", auth, 1)

	payload, err := httpGet(url, client)
	if err != nil {
		return err
	}

	resp := InboundOptionalResponse{}
	err = json.Unmarshal(payload, &resp)
	if err != nil {
		fmt.Printf("ERROR: json unmarshal of InboundOptionalResponse (%s)\n", err.Error())
		return err
	}

	fmt.Printf("received %d item(s)\n", len(resp.Details))
	return nil
}

//
// end of file
//
