package main

import(
        "testing"
)

func TestFindURLs(t *testing.T) {
        if len(FindURLs("Wie sieht es aus wenn man http://starship-factory.ch/blah?foo=bar http://foo?baz=quux sagt?")) == 0 {
                t.Error("no urls found")
        }
        if len(FindURLs("Wie sieht es aus wenn man http://starship-factory.ch/blah?foo=bar http://foo?baz=quux sagt?")) == 0 {
                t.Error("no urls found")
        }
}
