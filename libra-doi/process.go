//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/uvalib/easystore/uvaeasystore"
	librametadata "github.com/uvalib/libra-metadata"
	"github.com/uvalib/librabus-sdk/uvalibrabus"
)

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
	Funding string `json:"funding"`
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
	// For a new record, Datacite generate a DOI when empty
	DOI               string            `json:"doi,omitempty"`
	Prefix            string            `json:"prefix"`
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
	URL               string            `json:"url"`
	Publisher         string            `json:"publisher"`

	Affiliation AffiliationData `json:"affiliation"`
}

type PayloadData struct {
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

func process(messageId string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to librabus event
	ev, err := uvalibrabus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra bus event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("EVENT %s from:%s -> %s\n", messageId, messageSrc, ev.String())

	// load configuration
	cfg, err := loadConfiguration()
	if err != nil {
		return err
	}

	// easystore access
	es, err := newEasystore(cfg)
	if err != nil {
		fmt.Printf("ERROR: creating easystore (%s)\n", err.Error())
		return err
	}

	// important, cleanup properly
	defer es.Close()

	eso, err := getEasystoreObject(es, ev.Namespace, ev.Identifier)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	fields := eso.Fields()

	mdBytes, err := eso.Metadata().Payload()
	if err != nil {
		fmt.Printf("ERROR: unable to get metadata paload from respose: %s", err.Error())
		return err
	}
	work, err := librametadata.ETDWorkFromBytes(mdBytes)
	if err != nil {
		fmt.Printf("ERROR: unable to process paypad from work %s", err.Error())
		return err
	}

	fmt.Printf("Metadata: %+v\n", work)
	fmt.Printf("Fields: %+v\n", fields)

	cfg.httpClient = *newHTTPClient(1, 30)

	payload := PayloadData{}
	payload.Data.TypeName = "dois"
	payload.Data.Attributes = AttributesData{
		// string replace doi with blank
		DOI:    strings.Replace(fields["doi"], "doi:", "", 1),          // bare DOI
		Prefix: strings.Replace(cfg.IDService.Shoulder, "doi:", "", 1), // bare prefix numerals
		Titles: []TitleData{{Title: work.Title}},
		Descriptions: []DescriptionData{{
			Description:     work.Abstract,
			DescriptionType: "Abstract",
		}},
		Creators:     parseAuthor(work.Author),
		Contributors: parseAdvisors(work.Advisors),
		Subjects:     parseKeywords(work.Keywords),
		RightsList:   []RightsData{{}},

		Affiliation: UVAAffiliation(),
	}
	//      if work.description.present?
	//        attributes['descriptions'] = [{
	//          description: work.description,
	//          descriptionType: 'Abstract'
	//        }]
	//      end
	//
	//      attributes[:creators] = authors_construct( work )
	//      attributes[:contributors] = contributors_construct( work )
	//      attributes[:subjects] = work.keyword.map{|k| {subject: k}} if work.keyword.present?
	//      attributes[:rightsList] = [{rights: work.rights.first}] if work.rights.present? && work.rights.first.present?
	//      attributes[:fundingReferences] = work.sponsoring_agency.map{|f| {funderName: f}} if work.sponsoring_agency.present?
	//      attributes[:types] = {resourceTypeGeneral: DC_GENERAL_TYPE_TEXT, resourceType: RESOURCE_TYPE_DISSERTATION}
	//
	//      yyyymmdd = extract_yyyymmdd_from_datestring( work.date_published )
	//      yyyymmdd = extract_yyyymmdd_from_datestring( work.date_created ) if yyyymmdd.blank?
	//      attributes[:dates] = [{date: yyyymmdd, dateType: 'Issued'}] if yyyymmdd.present?
	//      attributes[:publicationYear] = yyyymmdd.first(4) if yyyymmdd.present?
	//
	//      attributes[:url] = fully_qualified_work_url( work.id ) # 'http://google.com'
	//      attributes[:publisher] = work.publisher if work.publisher.present?
	//
	//      #puts "==> #{h.to_json}"
	//      payload = {
	//        data: {
	//          type: 'dois',
	//          attributes: attributes
	//        }
	//      }

	fmt.Printf("%+v\n", payload)

	// Check DOI
	if len(fields["doi"]) == 0 {
		// No DOI present. Create one.
		fmt.Printf("INFO: DOI blank\n")
		doi, err := createDOI(cfg, eso)
		if err != nil {
			panic(err)
		}

		fmt.Println("New DOI: " + doi)

	} else {
		// Update DOI
		fmt.Printf("INFO: DOI for %s = %s\n", ev.Identifier, fields["doi"])
	}

	return nil
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

func parseAuthor(author librametadata.StudentData) []PersonData {
	var person PersonData
	person.GivenName = author.FirstName
	person.FamilyName = author.LastName
	person.NameType = "Personal"
	person.Affiliation = AffiliationData{Name: author.Institution}

	// Check for ORCID Accoun
	if false && len(author.ORCID) > 0 {
		person.NameIdentifiers = NameIdentifierData{
			SchemeURI:            "https://orcid.org",
			NameIdentifier:       "Author's ORCID",
			NameIdentifierScheme: "ORCID",
		}
	}

	person.Affiliation = UVAAffiliation()

	return []PersonData{person}

}

func parseAdvisors(contributors []librametadata.ContributorData) []PersonData {

	var person PersonData
	//Check for ORCID Account here
	if false {
		person.NameIdentifiers = NameIdentifierData{
			SchemeURI:            "https://orcid.org",
			NameIdentifier:       "https://orcid.org/0000-0002-2222-3333",
			NameIdentifierScheme: "ORCID",
		}
	}
	return []PersonData{person}

}

//
// end of file
//
