package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
}
type DescriptionData struct {
	Description     string `json:"description"`
	DescriptionType string `json:"descriptionType"`
}
type PersonData struct {
	GivenName       string               `json:"givenName"`
	FamilyName      string               `json:"familyName"`
	NameType        string               `json:"nameType"`
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
	Subject string `json:"subject"`
}
type RightsData struct {
	Rights string `json:"rights"`
}
type FundingData struct {
	FunderName string `json:"funderName"`
}
type TypeData struct {
	ResourceType        string `json:"resourceType,omitempty"`
	ResourceTypeGeneral string `json:"resourceTypeGeneral,omitempty"`
}
type DateData struct {
	Date     string `json:"date"`
	DateType string `json:"dateType"`
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
	Publisher         string            `json:"publisher,omitempty"`

	Affiliation AffiliationData `json:"affiliation"`
}

type DataciteData struct {
	Data struct {
		ID         string         `json:"id,omitempty"`
		TypeName   string         `json:"type"`
		Attributes AttributesData `json:"attributes"`
	} `json:"data"`
}

func UVAAffiliation() AffiliationData {
	return AffiliationData{
		Name:                        "University of Virginia",
		SchemeURI:                   "https://ror.org",
		AffiliationIdentifier:       "https://ror.org/0153tk833",
		AffiliationIdentifierScheme: "ROR",
	}
}

func getGeneralResourceType(cfg *Config, resourceTypeStr string) string {
	for _, rt := range cfg.ResourceTypes {
		if rt.Value == resourceTypeStr {
			return rt.Category
		}
	}
	return ""
}

func createETDPayload(work *librametadata.ETDWork, cfg *Config, fields uvaeasystore.EasyStoreObjectFields) DataciteData {
	var payload = DataciteData{}
	payload.Data.TypeName = "dois"
	payload.Data.Attributes = AttributesData{
		// remove doi: prefix
		DOI:    strings.Replace(fields["doi"], "doi:", "", 1),          // bare DOI
		Prefix: strings.Replace(cfg.IDService.Shoulder, "doi:", "", 1), // bare prefix numerals
		Titles: []TitleData{{Title: work.Title}},
		Descriptions: []DescriptionData{{
			Description:     work.Abstract,
			DescriptionType: "Abstract",
		}},
		Creators:          []PersonData{parseContributor(work.Author)},
		Contributors:      parseContributors(work.Advisors),
		Subjects:          parseKeywords(work.Keywords),
		RightsList:        []RightsData{{Rights: work.License}},
		FundingReferences: parseSponsors(work.Sponsors),

		Affiliation: UVAAffiliation(),
		Types: TypeData{
			ResourceTypeGeneral: "Text",
			ResourceType:        "Dissertation",
		},
		Publisher: "University of Virginia",
	}
	addDates(&payload, fields["publish-date"])
	return payload
}
func createOAPayload(work *librametadata.OAWork, cfg *Config, fields uvaeasystore.EasyStoreObjectFields) DataciteData {
	var payload = DataciteData{}
	payload.Data.TypeName = "dois"
	payload.Data.Attributes = AttributesData{
		// remove doi: prefix
		DOI:    strings.Replace(fields["doi"], "doi:", "", 1),          // bare DOI
		Prefix: strings.Replace(cfg.IDService.Shoulder, "doi:", "", 1), // bare prefix numerals
		Titles: []TitleData{{Title: work.Title}},
		Descriptions: []DescriptionData{{
			Description:     work.Abstract,
			DescriptionType: "Abstract",
		}},
		Creators:          parseContributors(work.Authors),
		Contributors:      parseContributors(work.Contributors),
		Subjects:          parseKeywords(work.Keywords),
		RightsList:        []RightsData{{Rights: work.License}},
		FundingReferences: parseSponsors(work.Sponsors),

		Affiliation: UVAAffiliation(),
		Types: TypeData{
			ResourceTypeGeneral: getGeneralResourceType(cfg, work.ResourceType),
			ResourceType:        work.ResourceType,
		},
		Publisher: work.Publisher,
	}
	addDates(&payload, fields["publish-date"])
	return payload
}

func addDates(payload *DataciteData, publishDate string) {

	parsedDate, err := time.Parse(time.RFC3339, publishDate)
	if err != nil {
		fmt.Printf("WARN: unable to parse date %s\n", err.Error())
		return
	}

	payload.Data.Attributes.Dates = []DateData{{
		Date:     parsedDate.Format("2006-01-02"),
		DateType: "Issued",
	}}
	payload.Data.Attributes.PublicationYear = parsedDate.Format("2006")
}

func parseSponsors(s []string) []FundingData {
	fundingList := []FundingData{}
	for _, sponsor := range s {
		fundingList = append(fundingList, FundingData{FunderName: sponsor})
	}
	return fundingList
}

func parseKeywords(keywords []string) []SubjectData {
	// Keywords are mapped to subjects
	subjects := []SubjectData{}
	for _, keyword := range keywords {
		subjects = append(subjects, SubjectData{Subject: keyword})
	}
	return subjects

}

func sendToDatacite(cfg *Config, payload *DataciteData) (string, error) {
	var response []byte
	var httpMethod string
	if len(payload.Data.Attributes.DOI) == 0 {
		httpMethod = "POST"
	} else {
		httpMethod = "PUT"
	}

	jsonPayload, _ := json.Marshal(payload)

	req, err := http.NewRequest(httpMethod, cfg.IDService.BaseURL+"/dois", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}
	req.Header.Add("content-type", "application/json")
	req.SetBasicAuth(cfg.IDService.User, cfg.IDService.Password)

	response, err = httpPost(&cfg.httpClient, req)
	spew.Dump(response)
	if err != nil {
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

func parseContributors(contributors []librametadata.ContributorData) []PersonData {
	var contribList []PersonData
	for _, contrib := range contributors {
		contribList = append(contribList, parseContributor(contrib))
	}
	return contribList
}

func parseContributor(contributor librametadata.ContributorData) PersonData {
	var person PersonData
	person.GivenName = contributor.FirstName
	person.FamilyName = contributor.LastName
	person.NameType = "Personal"
	if len(contributor.ComputeID) > 0 {
		person.Affiliation = []AffiliationData{UVAAffiliation()}
	} else {
		person.Affiliation = []AffiliationData{{Name: contributor.Institution}}
	}

	// Check for ORCID Account
	if false && len(contributor.ORCID) > 0 {
		person.NameIdentifiers = []NameIdentifierData{{
			SchemeURI:            "https://orcid.org",
			NameIdentifier:       "Author's ORCID",
			NameIdentifierScheme: "ORCID",
		}}
	}
	return person

}
