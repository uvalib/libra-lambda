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
	"regexp"
	"strings"
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

func updateAuthorOrcidActivity(config *Config, eso uvaeasystore.EasyStoreObject, authorId string, updateCode string, auth string, client *http.Client) (string, error) {

	// create the update schema
	schema, err := createUpdateSchema(eso)
	if err != nil {
		return "", err
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

	//fmt.Printf("PUT PAYLOAD [%s]\n", string(pl))

	buf, err := httpPut(client, url, pl)
	if err != nil {
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

	if eso.Namespace() == libraEtdNamespace {
		meta, err := librametadata.ETDWorkFromBytes(pl)
		if err != nil {
			return nil, err
		}

		schema.Authors = getEtdPersons(meta)
		schema.Abstract = meta.Abstract
		schema.PublicationDate = extractYYMMDD(meta.PublicationDate)
		schema.Title = meta.Title

		schema.ResourceType = "supervised-student-publication"
		if strings.Contains(meta.Degree, "Doctor") == true {
			schema.ResourceType = "dissertation-thesis"
		}
	}

	if eso.Namespace() == libraOpenNamespace {
		meta, err := librametadata.OAWorkFromBytes(pl)
		if err != nil {
			return nil, err
		}

		schema.Authors = getOpenPersons(meta)
		schema.Abstract = meta.Abstract
		schema.PublicationDate = extractYYMMDD(meta.PublicationDate)
		schema.ResourceType = MapResourceType(meta.ResourceType)
		schema.Title = meta.Title
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

	if eso.Namespace() == libraEtdNamespace {
		meta, err := librametadata.ETDWorkFromBytes(pl)
		if err != nil {
			return "", err
		}
		if len(meta.Author.ComputeID) != 0 {
			return meta.Author.ComputeID, nil
		}
	}

	if eso.Namespace() == libraOpenNamespace {
		meta, err := librametadata.OAWorkFromBytes(pl)
		if err != nil {
			return "", err
		}
		if len(meta.Authors) != 0 {
			if len(meta.Authors[0].ComputeID) != 0 {
				return meta.Authors[0].ComputeID, nil
			}
		}
	}

	// no error and no author
	return "", nil
}

func getEtdPersons(meta *librametadata.ETDWork) []Person {
	person := Person{Index: 0, FirstName: meta.Author.FirstName, LastName: meta.Author.LastName}
	return []Person{person}
}

func getOpenPersons(meta *librametadata.OAWork) []Person {
	persons := make([]Person, 0)
	for ix, p := range meta.Authors {
		person := Person{Index: ix, FirstName: p.FirstName, LastName: p.LastName}
		persons = append(persons, person)
	}
	return persons
}

func MapResourceType(rt string) string {
	switch rt {
	case "Article":
		return "journal-article"
	case "Book":
		return "book"
	case "Conference Paper":
		return "conference-paper"
	case "Part of Book":
		return "book-chapter"
	case "Report":
		return "report"
	case "Journal":
		return "journal-issue"
	case "Poster":
		return "conference-poster"
	default:
		return "other"
	}
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
