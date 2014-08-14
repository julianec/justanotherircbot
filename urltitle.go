package main

import (
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func logprivmsgs(event *irc.Event) {
	log.Print(event.Nick+": ", event.Arguments)
}

func writeurltitle(event *irc.Event) {
	var urls []string = FindURLs(event.Arguments[1])
	var err error
	var resp *http.Response
	var contentType string
	var respbody []byte

	for _, oneurl := range urls {
		// URL valid?
		_, err = url.Parse(oneurl)
		if err != nil {
			continue
		}
		resp, err = http.Head(oneurl)
		if err != nil {
			log.Print("Error getting Head: ", err)
			continue
		}

		// No HTML?
		contentType = resp.Header.Get("Content-Type")
		// Content type does not start with "text/html" or "application/xhtml+xml"?
		if !strings.HasPrefix(contentType, "text/html") && !strings.HasPrefix(contentType, "application/xhtml+xml") {
			log.Print("Wrong content type: ", contentType, " Expecting application/xhtml+xml or text/html")
			continue
		}

		// Get the Body
		resp, err = http.Get(oneurl)
		if err != nil {
			log.Print("Error during HTTP GET: ", err)
			continue
		}
		// Close later
		defer resp.Body.Close()

		// Create a slice of bytes from the body.
		respbody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Print("Error reading the body: ", err)
			continue
		}

		// What's the charset of the website?
		log.Print(respbody[0], "1. Byte vom Body")
		log.Print(contentType)
	}
}
