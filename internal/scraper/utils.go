package scraper

import (
	"regexp"
	"strings"
)

// containsAny checks if any of the substrings are found within the provided text.
func containsAny(text string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.Contains(text, substring) {
			return true
		}
	}
	return false
}

func extractExperience(text string) string {
	// Regex to match numbers (integers and floats) followed by "year" or "years"
	re := regexp.MustCompile(`(\d+(\.\d+)?)\s+year[s]?`)

	if strings.Contains(text, "year") {
		matches := re.FindStringSubmatch(text)
		if len(matches) > 0 {
			return matches[1] // Return the matched number as a string
		}
	}

	return "" // Return empty string if no experience found
}