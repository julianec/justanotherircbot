package main

import (
	"code.google.com/p/go-charset/charset"
        _ "code.google.com/p/go-charset/data"
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"github.com/thoj/go-ircevent"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func logprivmsgs(event *irc.Event) {
	log.Print(event.Nick+": ", event.Arguments)
}

type URLTitleExtractor struct {
	ircobject *irc.Connection
}

func (t *URLTitleExtractor) WriteURLTitle(event *irc.Event) {
	var urls []string = FindURLs(event.Arguments[1])
	var err error
	var resp *http.Response
	var contentType string
	var foundcharset string
	var ureader io.Reader
	var htmlnode *html.Node

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
		// Get the charset
		foundcharset = ExtractCharset(contentType)

		// Get the Body
		resp, err = http.Get(oneurl)
		if err != nil {
			log.Print("Error during HTTP GET: ", err)
			continue
		}
		// Close later
		defer resp.Body.Close()

		if strings.ToLower(foundcharset) != "utf-8" && strings.ToLower(foundcharset) != "utf8" {
			log.Print("Converting from ", foundcharset, " to UTF-8")
			ureader, err = charset.NewReader(foundcharset, resp.Body)
			if err != nil {
				log.Print("Error during utf-8 transformation: ", err)
				continue
			}
		} else {
			ureader = resp.Body
		}
		// Get the top HTML node
		htmlnode, err = html.Parse(ureader)
		if err != nil {
			log.Print("Error parsing HTML file: ", err)
			continue
		}
		var htmltag *html.Node = htmlnode.FirstChild // doctype, if well formed

		// Advance until we find the html tag or until no elements are left.
		for htmltag != nil && (htmltag.Type != html.ElementNode || htmltag.DataAtom != atom.Html) {
			htmltag = htmltag.NextSibling
		}
		// In case of broken HTML where everything is a top level element:
		if htmltag == nil {
			htmltag = htmlnode.FirstChild
		} else {
			htmlnode = htmltag // If head is missing we can continue from here
			htmltag = htmltag.FirstChild
		}
		for htmltag != nil && (htmltag.Type != html.ElementNode || htmltag.DataAtom != atom.Head) {
			htmltag = htmltag.NextSibling
		}
		// In case of even more broken HTML where even the Head is missing
		if htmltag == nil {
			htmltag = htmlnode.FirstChild
		} else {
			htmlnode = htmltag
			htmltag = htmltag.FirstChild // Go into head's first child
		}
		// Continue until finding title element or no elements are left
		for htmltag != nil && (htmltag.Type != html.ElementNode || htmltag.DataAtom != atom.Title) {
			htmltag = htmltag.NextSibling
		}
		if htmltag != nil && htmltag.FirstChild != nil && htmltag.FirstChild.Type == html.TextNode {
			log.Print(htmltag.FirstChild.Data)
                        t.ircobject.Privmsg(event.Arguments[0], "Title: "+htmltag.FirstChild.Data)
		}
	}
}
