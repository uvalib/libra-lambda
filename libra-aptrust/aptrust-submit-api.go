//
//
//

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

//
// Service API
//

type SubmitRegisterRequest struct {
	ClientIdentifier string `json:"cid"`        // the client identifier
	Collection       string `json:"collection"` // the collection name for the submission (optional)
	Storage          string `json:"storage"`    // the APT storage to use for this submission (optional)
}

type SubmitRegisterResponse struct {
	SubmissionIdentifier string `json:"sid"`
	DepositBucket        string `json:"bucket"`
	DepositPath          string `json:"path"`
}

type SubmitInitiateRequest struct {
	ClientIdentifier     string   `json:"cid"`         // the client identifier
	SubmissionIdentifier string   `json:"sid"`         // the submission identifier
	BagFolders           []string `json:"bag_folders"` // the bags to be included in this submission
}

type SubmitInitiateResponse struct {
	Submission string    `json:"submission"`
	Status     string    `json:"status"`
	Updated    time.Time `json:"updated"`
	// other stuff
}

func registerSubmission(cfg *Config, httpClient *http.Client) (*SubmitRegisterResponse, error) {

	req := SubmitRegisterRequest{}
	req.ClientIdentifier = cfg.APTServiceClient
	pl, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("ERROR: json marshal of SubmitRegisterRequest (%s)\n", err.Error())
		return nil, err
	}

	// post the request
	pl, err = httpPost(httpClient, cfg.APTServiceRegister, pl, "application/json")
	if err != nil {
		return nil, err
	}

	// and process the response
	resp := SubmitRegisterResponse{}
	err = json.Unmarshal(pl, &resp)
	if err != nil {
		fmt.Printf("ERROR: json unmarshal of SubmitRegisterResponse (%s)\n", err.Error())
		return nil, err
	}

	fmt.Printf("INFO: submit register complete (%s)\n", resp.SubmissionIdentifier)
	return &resp, nil
}

func initiateSubmission(cfg *Config, httpClient *http.Client, sid string, bagName string) error {

	req := SubmitInitiateRequest{}
	req.ClientIdentifier = cfg.APTServiceClient
	req.SubmissionIdentifier = sid
	req.BagFolders = []string{bagName}

	pl, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("ERROR: json marshal of SubmitInitiateRequest (%s)\n", err.Error())
		return err
	}

	// post the request
	pl, err = httpPost(httpClient, cfg.APTServiceSubmit, pl, "application/json")
	if err != nil {
		return err
	}

	// and process the response
	resp := SubmitInitiateResponse{}
	err = json.Unmarshal(pl, &resp)
	if err != nil {
		fmt.Printf("ERROR: json unmarshal of SubmitInitiateResponse (%s)\n", err.Error())
		return err
	}

	fmt.Printf("INFO: submit initiate complete")
	return nil
}

//
// end of file
//
