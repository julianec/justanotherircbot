package main

import (
	"regexp"
)

// Global regexp vars
var url_re = regexp.MustCompile(`\bhttps?://[^\s]+\b`) // Finds URLs

func FindURLs(input string) []string {
	return url_re.FindAllString(input, 5)
}
