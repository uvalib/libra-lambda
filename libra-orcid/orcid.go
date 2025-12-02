//
//
//

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/uvalib/easystore/uvaeasystore"
	librametadata "github.com/uvalib/libra-metadata"
)

type OrcidActivityUpdateResponse struct {
	Status     int    `json:"status"`
	Message    string `json:"message"`
	UpdateCode string `json:"update_code"`
}

type OrcidActivityUpdate struct {
	UpdateCode string     `json:"update_code,omitempty"`
	Work       WorkSchema `json:"work,omitempty"`
}

type WorkSchema struct {
	Title           string   `json:"title,omitempty"`
	Abstract        string   `json:"abstract,omitempty"`
	PublicationDate string   `json:"publication_date,omitempty"`
	URL             string   `json:"url,omitempty"`
	Authors         []Person `json:"authors,omitempty"`
	ResourceType    string   `json:"resource_type,omitempty"`
}

type Person struct {
	Index     int    `json:"index"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

var ErrIncompleteData = fmt.Errorf("incomplete data")

func updateAuthorOrcidActivity(config *Config, eso uvaeasystore.EasyStoreObject, authorId string, updateCode string, auth string, client *http.Client) (string, error) {

	// create the update schema
	schema, err := createUpdateSchema(eso)
	if err != nil {
		return "", err
	}

	// ensure we have enough to do the activity update, otherwise it will be rejected
	if len(schema.Title) == 0 || len(schema.ResourceType) == 0 || len(schema.URL) == 0 {
		return "", ErrIncompleteData
	}

	// substitute values into url
	url := strings.Replace(config.OrcidSetActivityUrl, "{:id}", authorId, 1)
	url = strings.Replace(url, "{:auth}", auth, 1)

	// create the request payload
	req := OrcidActivityUpdate{}
	req.UpdateCode = updateCode
	req.Work = *schema

	pl, err := json.Marshal(req)
	if err != nil {
		fmt.Printf("ERROR: json marshal of OrcidActivityUpdate (%s)\n", err.Error())
		return "", err
	}

	buf, err := httpPut(client, url, pl, "application/json")
	if err != nil {
		fmt.Printf("ERROR: failed payload [%s]\n", string(pl))
		if buf != nil {
			// let's try and unmarshal the response anyway
			resp := OrcidActivityUpdateResponse{}
			err = json.Unmarshal(buf, &resp)
			if err == nil {
				// this is a special case
				if resp.Status == http.StatusConflict {
					fmt.Printf("INFO: reports update already applied [%s]\n", resp.Message)
					// assume all is well and we can ignore this
					return resp.UpdateCode, nil
				}
			}
			fmt.Printf("ERROR: failed response [%s]\n", string(buf))
		}
		return "", err
	}

	//fmt.Printf("PUT RESPONSE [%s]\n", string(buf))

	resp := OrcidActivityUpdateResponse{}
	err = json.Unmarshal(buf, &resp)
	if err != nil {
		fmt.Printf("ERROR: json unmarshal of OrcidActivityUpdateResponse (%s)\n", err.Error())
		return "", err
	}

	// all good apparently
	return resp.UpdateCode, nil
}

func createUpdateSchema(eso uvaeasystore.EasyStoreObject) (*WorkSchema, error) {

	// check we have metadata
	md := eso.Metadata()
	if md == nil {
		return nil, ErrNoMetadata
	}
	pl, err := md.Payload()
	if err != nil {
		return nil, err
	}

	// this is our update schema
	schema := WorkSchema{}
	fields := eso.Fields()
	schema.URL = fields["doi"]

	meta, err := librametadata.ETDWorkFromBytes(pl)
	if err != nil {
		return nil, err
	}

	schema.Authors = getEtdPersons(meta)
	schema.Abstract = truncateString(meta.Abstract, 5000)
	schema.PublicationDate = extractYYMMDD(fields["publish-date"])
	schema.Title = meta.Title

	schema.ResourceType = "supervised-student-publication"
	if strings.Contains(meta.Degree, "Doctor") == true {
		schema.ResourceType = "dissertation-thesis"
	}

	return &schema, nil
}

func getWorkAuthor(eso uvaeasystore.EasyStoreObject) (string, error) {

	// check we have metadata
	md := eso.Metadata()
	if md == nil {
		return "", ErrNoMetadata
	}
	pl, err := md.Payload()
	if err != nil {
		return "", err
	}

	meta, err := librametadata.ETDWorkFromBytes(pl)
	if err != nil {
		return "", err
	}
	if len(meta.Author.ComputeID) != 0 {
		return meta.Author.ComputeID, nil
	}

	// no error and no author
	return "", nil
}

func getEtdPersons(meta *librametadata.ETDWork) []Person {
	person := Person{Index: 0, FirstName: meta.Author.FirstName, LastName: meta.Author.LastName}
	return []Person{person}
}

// truncateString limits a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// attempt to extract a 4 digit year from the date string (crap, I know)
func extractYYMMDD(date string) string {
	if len(date) == 0 {
		return ""
	}

	re := regexp.MustCompile("\\d{4}")
	if re.MatchString(date) == true {
		return re.FindAllString(date, 1)[0]
	}
	return ""
}

//
// end of file
//
