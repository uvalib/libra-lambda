//
//
//

package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/uvalib/easystore/uvaeasystore"
	"regexp"
	"strings"
	"time"
)

func titleSort(title string, languages []string) string {

	// first, to lower case
	str := strings.ToLower(title)

	// remove all non numeric/alpha/space
	nonAlphanumericRegex := regexp.MustCompile(`[^\p{L}\p{N} ]+`)
	str = nonAlphanumericRegex.ReplaceAllString(str, "")

	// remove leading transitional words
	str = removeLeadingTransitional(str, languages)

	// remove duplicate whitespace
	whitespace := regexp.MustCompile(`\s+`)
	str = whitespace.ReplaceAllString(str, " ")

	// finally convert spaces to underscores
	str = strings.Replace(str, " ", "_", -1)

	return str
}

func titleSuffix(first string, last string) string {
	return "/" + strings.ToLower(last) + "/" + strings.ToLower(first) + "/Thesis"
}

func removeLeadingTransitional(title string, languages []string) string {

	if listContains("English", languages) == true {
		return removeLeadingEnglish(title)
	}

	if listContains("French", languages) == true {
		return removeLeadingFrench(title)
	}

	if listContains("Italian", languages) == true {
		return removeLeadingItalian(title)
	}

	if listContains("Spanish", languages) == true {
		return removeLeadingSpanish(title)
	}

	if listContains("German", languages) == true {
		return removeLeadingGerman(title)
	}

	// default behavior
	return removeLeadingEnglish(title)
}

func removeLeadingEnglish(title string) string {
	regx := regexp.MustCompile(`(^the |^a |^an )`)
	return regx.ReplaceAllString(title, "")
}

func removeLeadingFrench(title string) string {
	regx := regexp.MustCompile(`(^la |^le |^l&apos;|^les |^une |^un |^des )`)
	return regx.ReplaceAllString(title, "")
}

func removeLeadingItalian(title string) string {
	regx := regexp.MustCompile(`(^uno |^una |^un |^un&apos;|^lo |^gli |^il |^i |^l&apos;|^la |^le )`)
	return regx.ReplaceAllString(title, "")
}

func removeLeadingSpanish(title string) string {
	regx := regexp.MustCompile(`(^el |^los |^las |^un |^una |^unos |^unas )`)
	return regx.ReplaceAllString(title, "")
}

func removeLeadingGerman(title string) string {
	regx := regexp.MustCompile(`(^der |^die |^das |^den |^dem |^des |^ein |^eine[mnr]?|^keine |^[k]?einer )`)
	return regx.ReplaceAllString(title, "")
}

func listContains(key string, list []string) bool {
	k := strings.ToLower(key)
	for _, v := range list {
		if strings.ToLower(v) == k {
			return true
		}
	}
	return false
}

// attempt to clean up the date for indexing (crap, I know)
func cleanupDate(date string) string {

	// remove periods, commas and a trailing 'th' on the date
	clean := strings.Replace(date, ".", "", -1)
	clean = strings.Replace(clean, "th,", "", -1)
	clean = strings.Replace(clean, ",", "", -1)

	// first try "YYYY"
	format := "2006"
	str, err := makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "YYYY-MM-DD"
	format = "2006-01-02"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "Month (short) Day, YYYY"
	format = "Jan 2 2006"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "Month (long) Day, YYYY"
	format = "January 2 2006"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "Month (short) YYYY"
	format = "Jan 2006"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "Month (long) YYYY"
	format = "January 2006"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "MM/DD/YYYY"
	format = "01/02/2006"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "YYYY/MM/DD"
	format = "2006/01/02"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "Day Month (short) YYYY"
	format = "2 Jan 2006"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "Day Month (long) YYYY"
	format = "2 January 2006"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "YYYY-MM"
	format = "2006-01"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "M/D/YYYY"
	format = "1/2/2006"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "M/D/YY"
	format = "1/2/06"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next try "M-D-YYYY"
	format = "1-2-2006"
	str, err = makeDate(clean, format)
	if err == nil {
		return str
	}

	// next this
	if len(clean) > 19 {
		format = "2006-01-02T15:04:05"
		str, err = makeDate(clean[:19], format)
		if err == nil {
			return str
		}
	}

	// really finally
	str = extractYYYY(clean)
	if len(str) != 0 {
		return fmt.Sprintf("%s-01-01T00:00:00Z", str)
	}

	fmt.Printf("ERROR: unable intrpret date [%s]\n", date)

	return date
}

// make a fixed format date given a date and expected format
func makeDate(date string, format string) (string, error) {
	tm, err := time.Parse(format, date)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ",
		tm.Year(), tm.Month(), tm.Day(), tm.Hour(), tm.Minute(), tm.Second()), nil
}

// attempt to extract a 4 digit year from the date string (crap, I know)
func extractYYYY(date string) string {
	if len(date) == 0 {
		return ""
	}

	re := regexp.MustCompile("\\d{4}")
	if re.MatchString(date) == true {
		return re.FindAllString(date, 1)[0]
	}
	return ""
}

// taken from https://github.com/uvalib/v4-libra-indexer/blob/master/libraoc/LibraOCToVirgo4.xsl
func poolAdditional(resourceType string) string {

	switch resourceType {
	case "Audio":
		return "sound_recordings"
	case "Book":
		return "catalog"
	case "Image":
		return "Visual Materials"
	case "Journals":
		return "serials"
	case "Map or Cartographic Material":
		return "maps"
	case "Part of Book":
		return "catalog"
	case "Video":
		return "video"
	default:
		return ""
	}
}

// visibility in the index
func workVisibility(fields uvaeasystore.EasyStoreObjectFields) string {

	// possible results
	hidden := "HIDDEN"
	visible := "VISIBLE"

	// all draft works are hidden
	if fields["draft"] == "true" {
		return hidden
	}

	// restricted works are hidden
	if fields["default-visibility"] == "restricted" {
		return hidden
	}

	// must be visible then
	return visible
}

func XmlEncode(str string) string {
	var b bytes.Buffer
	err := xml.EscapeText(&b, []byte(str))
	if err != nil {
		fmt.Printf("ERROR: escaping (%s)\n", err.Error())
		return str
	}
	return string(b.Bytes())
}

//
// end of file
//
