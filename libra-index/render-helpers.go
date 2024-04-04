//
//
//

package main

import (
	"regexp"
	"strings"
)

func titleSort(title string, languages []string) string {

	// first, to lower case
	str := strings.ToLower(title)

	// remove all non numeric/alpha/space
	nonAlphanumericRegex := regexp.MustCompile(`[^\p{L}\p{N} ]+`)
	str = nonAlphanumericRegex.ReplaceAllString(str, "")

	// remove leading transitional words
	str = removeLeadingTransitional(str, languages)

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

//
// end of file
//
