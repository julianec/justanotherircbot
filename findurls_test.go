package main

import(
        "testing"
)

func TestFindURLs(t *testing.T) {
        var res []string
        var teststrings = map[string][]string{
                "Wie sieht es aus wenn man http://starship-factory.ch/blah?foo=bar sagt?":[]string{"http://starship-factory.ch/blah?foo=bar"},
                "Wie sieht es aus wenn man http://starship-factory.ch/blah?foo=bar http://foo?baz=quux sagt?":[]string{"http://starship-factory.ch/blah?foo=bar", "http://foo?baz=quux"},
                "http://starship-factory.ch/blah?foo=bar.":[]string{"http://starship-factory.ch/blah?foo=bar"},
        }
        for teststring, expected := range teststrings {
                res = FindURLs(teststring)
                if len(res) != len(expected) {
                        t.Error("Wrong number of URLs. Expected: ", expected, ", got ", res, ".")
                }
        }
}
