//
//
//

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthResponse struct {
	Expires string `json:"expires"`
	Token   string `json:"token"`
}

func getAuthToken(client *http.Client, url string) (string, error) {

	payload, err := httpGet(client, url)
	if err != nil {
		return "", err
	}

	resp := AuthResponse{}
	err = json.Unmarshal(payload, &resp)
	if err != nil {
		fmt.Printf("ERROR: json unmarshal of AuthResponse (%s)\n", err.Error())
		return "", err
	}

	return resp.Token, nil
}

//
// end of file
//
