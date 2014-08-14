package main

import (
	"regexp"
)

// Global regexp vars
var url_re = regexp.MustCompile(`\b(https?://\S+)[\s\.$]`) // Finds URLs

func FindURLs(input string) (ret []string) {
        var matches = url_re.FindAllStringSubmatch(input, 5)
        for _, value := range matches {
                if len(value) >= 2 {
                        ret = append(ret, value[1])
                }
        }
        return
}
