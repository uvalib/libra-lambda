//
// main message processing
//

package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

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
	// For a new record, Datacite generate a DOI when empty
	DOI               string            `json:"doi,omitempty"`
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

func process(messageID string, messageSrc string, rawMsg json.RawMessage) error {

	// convert to librabus event
	ev, err := uvalibrabus.MakeBusEvent(rawMsg)
	if err != nil {
		fmt.Printf("ERROR: unmarshaling libra bus event (%s)\n", err.Error())
		return err
	}

	fmt.Printf("EVENT %s from:%s -> %s\n", messageID, messageSrc, ev.String())

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

	eso, err := getEasystoreObjectByKey(es, ev.Namespace, ev.Identifier, uvaeasystore.Fields+uvaeasystore.Metadata)
	if err != nil {
		fmt.Printf("ERROR: getting object ns/oid [%s/%s] (%s)\n", ev.Namespace, ev.Identifier, err.Error())
		return err
	}

	fields := eso.Fields()

	if eso.Metadata() == nil {
		fmt.Printf("ERROR: unable to get metadata payload for ns/oid [%s/%s]\n", ev.Namespace, ev.Identifier)
		return ErrNoMetadata
	}

	mdBytes, err := eso.Metadata().Payload()
	if err != nil {
		fmt.Printf("ERROR: unable to get metadata payload from response: %s\n", err.Error())
		return err
	}

	fmt.Printf("Fields: %+v\n", fields)

	var payload PayloadData
	if ev.Namespace == cfg.ETDNamespace.Name {
		work, err := librametadata.ETDWorkFromBytes(mdBytes)
		if err != nil {
			fmt.Printf("ERROR: unable to process ETD Work %s\n", err.Error())
			return err
		}

		fmt.Printf("Metadata: %+v\n", work)
		payload = createETDPayload(work, cfg, fields)

	} else if ev.Namespace == cfg.OpenNamespace.Name {
		work, err := librametadata.OAWorkFromBytes(mdBytes)
		if err != nil {
			fmt.Printf("ERROR: unable to process OA Work  %s\n", err.Error())
			return err
		}

		fmt.Printf("Metadata: %+v\n", work)
		payload = createOAPayload(work, cfg, fields)
	}
	cfg.httpClient = *newHttpClient(1, 30)
	fmt.Printf("Payload: %+v\n", payload)

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

func getGeneralResourceType(cfg *Config, resourceTypeStr string) string {
	for _, rt := range cfg.ResourceTypes {
		if rt.Value == resourceTypeStr {
			return rt.Category
		}
	}
	return ""
}

func createETDPayload(work *librametadata.ETDWork, cfg *Config, fields uvaeasystore.EasyStoreObjectFields) PayloadData {
	var payload = PayloadData{}
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
	}
	return payload
}
func createOAPayload(work *librametadata.OAWork, cfg *Config, fields uvaeasystore.EasyStoreObjectFields) PayloadData {
	var payload = PayloadData{}
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

func addDates(payload *PayloadData, dateStr string) {

	parsedDate, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		fmt.Printf("ERROR: unable to parse date %s\n", err.Error())
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

//
// end of file
//
