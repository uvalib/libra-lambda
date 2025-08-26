package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/uvalib/easystore/uvaeasystore"
	librametadata "github.com/uvalib/libra-metadata"
)

// https://support.datacite.org/reference/introduction

type AffiliationData struct {
	Name                        string `json:"name"`
	SchemeURI                   string `json:"schemeUri,omitempty"`
	AffiliationIdentifier       string `json:"affiliationIdentifier,omitempty"`
	AffiliationIdentifierScheme string `json:"affiliationIdentifierScheme,omitempty"`
}

type TitleData struct {
	Title string `json:"title"`
	Lang  string `json:"lang,omitempty"`
}
type DescriptionData struct {
	Description     string `json:"description"`
	DescriptionType string `json:"descriptionType"`
}
type PersonData struct {
	Name            string               `json:"name,omitempty"`
	GivenName       string               `json:"givenName"`
	FamilyName      string               `json:"familyName"`
	NameType        string               `json:"nameType,omitempty"`
	ContributorType string               `json:"contributorType,omitempty"`
	Affiliation     []AffiliationData    `json:"affiliation,omitempty"`
	NameIdentifiers []NameIdentifierData `json:"nameIdentifiers,omitempty"`
}
type NameIdentifierData struct {
	SchemeURI            string `json:"schemeUri"`
	NameIdentifier       string `json:"nameIdentifier"`
	NameIdentifierScheme string `json:"nameIdentifierScheme"`
}
type SubjectData struct {
	Subject string `json:"subject,omitempty"`
}
type RightsData struct {
	Rights string `json:"rights,omitempty"`
}
type FundingData struct {
	FunderName string `json:"funderName,omitempty"`
}
type TypeData struct {
	ResourceType        string `json:"resourceType,omitempty"`
	ResourceTypeGeneral string `json:"resourceTypeGeneral,omitempty"`
}
type DateData struct {
	Date     string `json:"date"`
	DateType string `json:"dateType"`
}

type PublisherData struct {
	Name                      string `json:"name"`
	SchemeURI                 string `json:"schemeUri,omitempty"`
	PublisherIdentifier       string `json:"publisherIdentifier,omitempty"`
	PublisherIdentifierScheme string `json:"publisherIdentifierScheme,omitempty"`
}

type AttributesData struct {
	Event             string            `json:"event,omitempty"` // eg: publish
	DOI               string            `json:"doi,omitempty"`   // Datacite generates a DOI when empty
	Prefix            string            `json:"prefix"`
	URL               string            `json:"url"`
	Titles            []TitleData       `json:"titles"`
	Descriptions      []DescriptionData `json:"descriptions,omitempty"`
	Creators          []PersonData      `json:"creators,omitempty"`
	Contributors      []PersonData      `json:"contributors,omitempty"`
	Subjects          []SubjectData     `json:"subjects,omitempty"`
	RightsList        []RightsData      `json:"rightsList,omitempty"`
	FundingReferences []FundingData     `json:"fundingReferences,omitempty"`
	Types             TypeData          `json:"types,omitempty"`
	Dates             []DateData        `json:"dates,omitempty"`
	PublicationYear   string            `json:"publicationYear,omitempty"`
	Publisher         PublisherData     `json:"publisher"`
}

type DataciteData struct {
	Data struct {
		TypeName   string         `json:"type"`
		Attributes AttributesData `json:"attributes"`
	} `json:"data"`
}

// UVAAffiliation with ROR ID
func UVAAffiliation() AffiliationData {
	return AffiliationData{
		Name:                        "University of Virginia",
		SchemeURI:                   "https://ror.org",
		AffiliationIdentifier:       "https://ror.org/0153tk833",
		AffiliationIdentifierScheme: "ROR",
	}
}

// UVAPublisher with ROR ID
func UVAPublisher() PublisherData {
	return PublisherData{
		Name:                      "University of Virginia",
		PublisherIdentifier:       "https://ror.org/0153tk833",
		SchemeURI:                 "https://ror.org",
		PublisherIdentifierScheme: "ROR",
	}
}

func createETDPayload(work *librametadata.ETDWork, fields uvaeasystore.EasyStoreObjectFields) DataciteData {
	var payload = DataciteData{}
	payload.Data.TypeName = "dois"
	// remove http://doi... prefix
	lastPath := regexp.MustCompile("[^/]+$")
	suffix := lastPath.FindString(fields["doi"])
	doi := ""
	if len(suffix) > 0 {
		doi = Cfg().IDService.Shoulder + "/" + suffix
	}

	payload.Data.Attributes = AttributesData{
		DOI:               doi,
		Prefix:            Cfg().IDService.Shoulder,
		Titles:            []TitleData{{Title: work.Title}},
		Descriptions:      getDescription(work.Abstract),
		Creators:          []PersonData{getPerson(work.Author, "")},
		Contributors:      getPersonList(work.Advisors, "RelatedPerson"),
		Subjects:          getKeywords(work.Keywords),
		RightsList:        getRights(work.License),
		FundingReferences: getSponsors(work.Sponsors),

		Types: TypeData{
			ResourceTypeGeneral: "Text",
			ResourceType:        "Dissertation",
		},
		Publisher: UVAPublisher(),
	}
	addDates(&payload, fields["publish-date"])
	return payload
}

func getDescription(abstract string) []DescriptionData {
	if abstract == "" {
		return []DescriptionData{}
	}
	return []DescriptionData{{
		Description:     abstract,
		DescriptionType: "Abstract",
	}}
}
func getRights(rights string) []RightsData {
	if rights == "" {
		return []RightsData{}
	}
	return []RightsData{{Rights: rights}}
}

func addDates(payload *DataciteData, publishDate string) {
	if len(publishDate) == 0 {
		currentYear := time.Now().Year()
		fmt.Printf("INFO: setting empty PublicationYear to: %d\n", currentYear)
		payload.Data.Attributes.PublicationYear = fmt.Sprintf("%d", currentYear)
		return
	}

	parsedDate, err := time.Parse(time.RFC3339, publishDate)
	if err != nil {
		fmt.Printf("WARN: unable to parse publish date %s\n", err.Error())
		return
	}

	payload.Data.Attributes.Dates = []DateData{{
		Date:     parsedDate.Format("2006-01-02"),
		DateType: "Issued",
	}}
	payload.Data.Attributes.PublicationYear = parsedDate.Format("2006")
}

func getSponsors(s []string) []FundingData {
	fundingList := []FundingData{}
	for _, sponsor := range s {
		fundingList = append(fundingList, FundingData{FunderName: sponsor})
	}
	return fundingList
}

func getKeywords(keywords []string) []SubjectData {
	// Keywords are mapped to subjects
	subjects := []SubjectData{}
	for _, keyword := range keywords {
		subjects = append(subjects, SubjectData{Subject: keyword})
	}
	return subjects

}

func sendToDatacite(payload *DataciteData) (string, error) {
	var response []byte
	var httpMethod, path string

	if len(payload.Data.Attributes.DOI) == 0 {
		// no DOI
		httpMethod = "POST"
		path = "/dois"

	} else {
		// update existing
		httpMethod = "PUT"
		// payload.Data.Attributes.DOI format should be: 10.18130/xxxx
		path = fmt.Sprintf("/dois/%s", payload.Data.Attributes.DOI)
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	fmt.Printf("INFO: JSON Payload to Datacite:\n%s\n", jsonPayload)

	req, err := http.NewRequest(httpMethod, Cfg().IDService.BaseURL+path, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	req.Header.Add("content-type", "application/vnd.api+json")
	req.Header.Add("accept", "application/json")
	req.SetBasicAuth(Cfg().IDService.User, Cfg().IDService.Password)

	response, err = httpSend(Cfg().httpClient, req)
	if err != nil {
		spew.Dump(req)
		spew.Dump(response)
		return "", err
	}

	// after success, we only care about the DOI here
	type DataciteResponseData struct {
		Data struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	var responseData DataciteResponseData
	err = json.Unmarshal(response, &responseData)
	if err != nil {
		return "", err
	}
	return responseData.Data.ID, nil
}

func getPersonList(contributors []librametadata.ContributorData, typeName string) []PersonData {
	var personList []PersonData
	for _, person := range contributors {
		if len(person.FirstName) > 0 || len(person.LastName) > 0 {
			personList = append(personList, getPerson(person, typeName))
		}
	}
	return personList
}

func getPerson(contributor librametadata.ContributorData, contribType string) PersonData {
	var person PersonData
	person.GivenName = contributor.FirstName
	person.FamilyName = contributor.LastName
	if len(person.GivenName) > 0 && len(person.FamilyName) > 0 {
		person.Name = fmt.Sprintf("%s, %s", person.FamilyName, person.GivenName)
	}
	person.ContributorType = contribType
	person.NameType = "Personal"
	if len(contributor.ComputeID) > 0 {
		person.Affiliation = []AffiliationData{UVAAffiliation()}
	} else {
		person.Affiliation = []AffiliationData{{Name: contributor.Institution}}
	}

	// Check for ORCID Account
	if contributor.ComputeID != "" {
		orcid, err := getOrcidDetails(Cfg().OrcidGetDetailsURL, contributor.ComputeID, Cfg().AuthToken, Cfg().httpClient)
		if err != nil {
			fmt.Printf("WARNING: unable to get ORCID details for %s\n", contributor.ComputeID)
		}

		if len(orcid) > 0 {
			person.NameIdentifiers = []NameIdentifierData{{
				NameIdentifier:       orcid,
				SchemeURI:            "https://orcid.org",
				NameIdentifierScheme: "ORCID",
			}}
		}
	}
	return person
}

//
// end of file
//
