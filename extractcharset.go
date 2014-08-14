package main

import (
	"regexp"
)

var charset_re = regexp.MustCompile(`\bcharset="?(\S+)"?\b`) // Find's the word after "charset="

func ExtractCharset(contentType string) string {
	var charset = charset_re.FindStringSubmatch(contentType)
	if len(charset) < 2 { // no match
		return ""
	}
	return charset[1]
}
