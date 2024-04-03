package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var maxHttpRetries = 3
var retrySleepTime = 100 * time.Millisecond

func newHttpClient(maxConnections int, timeout int) *http.Client {

	return &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: maxConnections,
		},
		Timeout: time.Duration(timeout) * time.Second,
	}
}

func httpGet(client *http.Client, url string) ([]byte, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("ERROR: GET %s failed with error (%s)\n", url, err)
		return nil, err
	}

	var response *http.Response
	count := 0
	for {
		start := time.Now()
		response, err = client.Do(req)
		duration := time.Since(start)
		fmt.Printf("INFO: GET %s (elapsed %d ms)\n", url, duration.Milliseconds())

		count++
		if err != nil {
			if canRetry(err) == false {
				fmt.Printf("ERROR: GET %s failed with error (%s)\n", url, err)
				return nil, err
			}

			// break when tried too many times
			if count >= maxHttpRetries {
				return nil, err
			}

			fmt.Printf("ERROR: GET %s failed with error, retrying (%s)\n", url, err)

			// sleep for a bit before retrying
			time.Sleep(retrySleepTime)
		} else {

			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				logLevel := "ERROR"
				// log not found as informational instead of as an error
				if response.StatusCode == http.StatusNotFound {
					logLevel = "INFO"
				}
				fmt.Printf("%s: GET %s failed with status %d\n", logLevel, url, response.StatusCode)

				body, _ := ioutil.ReadAll(response.Body)

				return body, fmt.Errorf("request returns HTTP %d", response.StatusCode)
			} else {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return nil, err
				}

				//fmt.Printf( body )
				return body, nil
			}
		}
	}
}

func httpPut(client *http.Client, url string, payload []byte) ([]byte, error) {

	var reader *bytes.Reader
	if payload != nil {
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequest("PUT", url, reader)
	if err != nil {
		fmt.Printf("ERROR: PUT %s failed with error (%s)\n", url, err)
		return nil, err
	}

	var response *http.Response
	count := 0
	for {
		start := time.Now()
		response, err = client.Do(req)
		duration := time.Since(start)
		fmt.Printf("INFO: PUT %s (elapsed %d ms)\n", url, duration.Milliseconds())

		count++
		if err != nil {
			if canRetry(err) == false {
				fmt.Printf("ERROR: PUT %s failed with error (%s)\n", url, err)
				return nil, err
			}

			// break when tried too many times
			if count >= maxHttpRetries {
				return nil, err
			}

			fmt.Printf("ERROR: PUT %s failed with error, retrying (%s)\n", url, err)

			// sleep for a bit before retrying
			time.Sleep(retrySleepTime)
		} else {

			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				logLevel := "ERROR"
				// log not found as informational instead of as an error
				if response.StatusCode == http.StatusNotFound {
					logLevel = "INFO"
				}
				fmt.Printf("%s: PUT %s failed with status %d\n", logLevel, url, response.StatusCode)

				body, _ := ioutil.ReadAll(response.Body)

				return body, fmt.Errorf("request returns HTTP %d", response.StatusCode)
			} else {
				body, err := ioutil.ReadAll(response.Body)
				if err != nil {
					return nil, err
				}

				//fmt.Printf( body )
				return body, nil
			}
		}
	}
}

// examines the error and decides if it can be retried
func canRetry(err error) bool {

	if strings.Contains(err.Error(), "operation timed out") == true {
		return true
	}

	if strings.Contains(err.Error(), "Client.Timeout exceeded") == true {
		return true
	}

	if strings.Contains(err.Error(), "write: broken pipe") == true {
		return true
	}

	if strings.Contains(err.Error(), "no such host") == true {
		return true
	}

	if strings.Contains(err.Error(), "network is down") == true {
		return true
	}

	return false
}

//
// end of file
//
