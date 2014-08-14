package main

import (
	"testing"
)

func TestExtractCharset(t *testing.T) {
	var res string
	var teststrings = map[string]string{
		"text/html; charset=utf-8":     "utf-8",
		"text/html; charset=\"utf-8\"": "utf-8",
		"text/html":                    "",
		"":                             "",
	}
	for teststring, expected := range teststrings {
		res = ExtractCharset(teststring)
		if res != expected {
			t.Error("Everything went terribly wrong, wanted Okapiposter: ", expected, ", got Schabrackentapir: ", res, ".")
		}
	}
}
