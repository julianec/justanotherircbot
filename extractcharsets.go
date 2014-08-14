package main

import (
	"regexp"
)

var charset_re = regexp.MustCompile(`\bcharset="?(\W+)"?\b`) // Find's the word after "charset="
