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

type OrcidDetailsResponse struct {
	Status  int            `json:"status"`
	Message string         `json:"message"`
	Details []OrcidDetails `json:"results"`
}

type OrcidDetails struct {
	//ID    string `json:"id,omitempty"`
	//Cid   string `json:"cid,omitempty"`
	Orcid string `json:"orcid,omitempty"`
	//URI   string `json:"uri,omitempty"`
}

func getOrcidDetails(url string, cid string, auth string, client *http.Client) (string, error) {

	// substitute values into url
	url = strings.Replace(url, "{:id}", cid, 1)
	url = strings.Replace(url, "{:auth}", auth, 1)

	payload, err := httpGet(client, url)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") == true {
			return "", nil
		}
		return "", err
	}

	resp := OrcidDetailsResponse{}
	err = json.Unmarshal(payload, &resp)
	if err != nil {
		fmt.Printf("ERROR: json unmarshal of OrcidDetailsResponse (%s)\n", err.Error())
		return "", err
	}

	// if we have details, return them
	if len(resp.Details) != 0 {
		fmt.Printf("INFO: located ORCID [%s] for cid [%s]\n", resp.Details[0].Orcid, cid)
		return resp.Details[0].Orcid, nil
	}

	// no error, nothing found
	return "", nil
}

//
// end of file
//
