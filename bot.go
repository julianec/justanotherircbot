package main

import (
	"flag"
	"github.com/thoj/go-ircevent"
	"log"
	"net/url"
	"regexp"
	"strings"
        "net/http"
        "io/ioutil"
//        "code.google.com/p/go.net/html"
)

var url_re = regexp.MustCompile(`\bhttps?://[^\s]+\b`)

func logprivmsgs(event *irc.Event) {
	log.Print(event.Nick+": ", event.Arguments)
}

func writeurltitle(event *irc.Event) {
	var urls []string = FindURLs(event.Arguments[1])
        var err error
        var resp *http.Response
        var contentType string
        var respbody []byte

        // URL valid?
        for _, oneurl := range urls {
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
                if contentType != "text/html" && contentType != "application/xhtml+xml" {
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
                log.Print(respbody[0], "1. Byte vom Body")
        }
}

func FindURLs(input string) []string {
	return url_re.FindAllString(input, 5)
}

func main() {
	var myircbot *irc.Connection
	var botname *string
	var serveraddress *string
	var rawchannellist *string
	var channellist []string
	var channelname string
	var err error

	botname = flag.String("botname", "justanotherbot", "Name of the bot")
	serveraddress = flag.String("server-address", "irc.freenode.org:6667", "Server Address")
	rawchannellist = flag.String("channels", "#ancient-solutions", "List of channels")
	flag.Parse()

	channellist = strings.Split(*rawchannellist, ",")

	myircbot = irc.IRC(*botname, *botname)
	if err = myircbot.Connect(*serveraddress); err != nil {
		log.Fatal("Error connecting to server: ", err)
	}

	//Join all channels.
	for _, channelname = range channellist {
		myircbot.Join(channelname)
	}

	//Event handling
	myircbot.AddCallback("PRIVMSG", logprivmsgs)
	myircbot.AddCallback("PRIVMSG", writeurltitle)

	myircbot.Loop()
}
