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

type UserDetailsResponse struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Details *UserDetails `json:"user"`
}

type UserDetails struct {
	UserID      string `json:"cid,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	Initials    string `json:"initials,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	Email       string `json:"email,omitempty"`
}

func getUserDetails(url string, cid string, auth string, client *http.Client) (*UserDetails, error) {

	// substitute values into url
	url = strings.Replace(url, "{:id}", cid, 1)
	url = strings.Replace(url, "{:auth}", auth, 1)

	payload, err := httpGet(client, url)
	if err != nil {
		if strings.Contains(err.Error(), "HTTP 404") == true {
			return nil, nil
		}
		return nil, err
	}

	resp := UserDetailsResponse{}
	err = json.Unmarshal(payload, &resp)
	if err != nil {
		fmt.Printf("ERROR: json unmarshal of UserDetailsResponse (%s)\n", err.Error())
		return nil, err
	}

	// if we have details, return them
	if resp.Details != nil {
		fmt.Printf("INFO: located details for cid [%s]\n", cid)
		return resp.Details, nil
	}

	// no error, nothing found
	return nil, nil
}

//
// end of file
//
