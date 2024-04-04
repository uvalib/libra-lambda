package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/uvalib/easystore/uvaeasystore"
	librametadata "github.com/uvalib/libra-metadata"
)

// https://support.datacite.org/reference/introduction

type AffiliationData struct {
	Name                        string `json:"name"`
	SchemeURI                   string `json:"schemeUri"`
	AffiliationIdentifier       string `json:"affiliationIdentifier"`
	AffiliationIdentifierScheme string `json:"affiliationIdentifierScheme"`
}

type TitleData struct {
	Title string `json:"title"`
}
type DescriptionData struct {
	Description     string `json:"description"`
	DescriptionType string `json:"descriptionType"`
}
type PersonData struct {
	GivenName       string             `json:"givenName"`
	FamilyName      string             `json:"familyName"`
	NameType        string             `json:"nameType"`
	ContributorType string             `json:"contributorType,omitempty"`
	Affiliation     AffiliationData    `json:"affiliation,omitempty"`
	NameIdentifiers NameIdentifierData `json:"nameIdentifier,omitempty"`
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
	ResourceType        string `json:"resourceType"`
	ResourceTypeGeneral string `json:"resourceTypeGeneral"`
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
	Descriptions      []DescriptionData `json:"descriptions"`
	Creators          []PersonData      `json:"creators"`
	Contributors      []PersonData      `json:"contributors"`
	Subjects          []SubjectData     `json:"subjects"`
	RightsList        []RightsData      `json:"rightsList"`
	FundingReferences []FundingData     `json:"fundingReferences"`
	Types             []TypeData        `json:"types"`
	Dates             []DateData        `json:"dates"`
	PublicationYear   string            `json:"publicationYear"`
	Publisher         string            `json:"publisher"`

	Affiliation AffiliationData `json:"affiliation"`
}

type DataciteData struct {
	Data struct {
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
		Types: []TypeData{{
			ResourceTypeGeneral: "Text",
			ResourceType:        "Dissertation",
		}},
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
		Types: []TypeData{{
			ResourceTypeGeneral: getGeneralResourceType(cfg, work.ResourceType),
			ResourceType:        work.ResourceType,
		}},
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

func createDOI(cfg *Config, obj uvaeasystore.EasyStoreObject) (string, error) {
	return "todo", nil

}

func updateMetadata(cfg *Config, obj uvaeasystore.EasyStoreObject) (string, error) {
	return "todo", nil
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
		person.Affiliation = UVAAffiliation()
	} else {
		person.Affiliation = AffiliationData{Name: contributor.Institution}
	}

	// Check for ORCID Account
	if false && len(contributor.ORCID) > 0 {
		person.NameIdentifiers = NameIdentifierData{
			SchemeURI:            "https://orcid.org",
			NameIdentifier:       "Author's ORCID",
			NameIdentifierScheme: "ORCID",
		}
	}
	return person

}
