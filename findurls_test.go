package main

import (
	"testing"
)

func TestFindURLs(t *testing.T) {
	var res []string
	var teststrings = map[string][]string{
		"Wie sieht es aus wenn man http://starship-factory.ch/blah?foo=1bar sagt?":                     []string{"http://starship-factory.ch/blah?foo=1bar"},
		"Wie sieht es aus wenn man https://starship-factory.ch/blah?foo=2bar sagt?":                    []string{"https://starship-factory.ch/blah?foo=2bar"},
		"Wie sieht es aus wenn man http://starship-factory.ch/blah?foo=3bar http://foo?baz=quux sagt?": []string{"http://starship-factory.ch/blah?foo=3bar", "http://foo?baz=quux"},
		"http://starship-factory.ch/blah?foo=b1ar.":                                                    []string{"http://starship-factory.ch/blah?foo=b1ar"},
		"Dieser freche Text enthaelt einfach keine URLs!":                                              []string{},
		"http://eins/ http://zwei/ http://drei/ http://vier/ http://fuenf/ http://sechs/ http://sie/":  []string{"http://eins/", "http://zwei/", "http://drei/", "http://vier/", "http://fuenf/"},
		"http://www.duroehre.com/?bewegtesbild=foobarbaz":                                              []string{"http://www.duroehre.com/?bewegtesbild=foobarbaz"},
		"Error getting Head: Head http://ipv4scan.com):: dial tcp: lookup ipv4scan.com): no":           []string{"http://ipv4scan.com"},
		"Error getting Head: Head )http://ipv4scan.com(:: dial tcp: lookup ipv4scan.com): no":          []string{"http://ipv4scan.com"},
		"Was machen wir eigentlich jetzt mit versehentlichen URLshttp://www.google.ch/ ?":              []string{"http://www.google.ch/"},
		"Wie sieht es aus wenn man http://starship-factory.ch:8080/ sagt?":                             []string{"http://starship-factory.ch:8080/"},
		"Wie sieht es aus wenn man http://starship-factory.ch:8080 sagt?":                              []string{"http://starship-factory.ch:8080"},
		"Wie sieht es aus wenn man http://starship-factory.ch:8080/foo/bar?baz=foo?bar=baz sagt?":      []string{"http://starship-factory.ch:8080/foo/bar?baz=foo?bar=baz"},
	}
	for teststring, expected := range teststrings {
		res = FindURLs(teststring)
		if len(res) != len(expected) {
			t.Error("Wrong number of URLs. Expected: ", expected, ", got ", res, ".")
		}
		for i := range res {
			if res[i] != expected[i] {
				t.Error("Wrong URL. Expected: ", expected[i], " got ", res[i], ".")
			}
		}
	}
}
