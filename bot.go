package main

import (
	"flag"
	"github.com/thoj/go-ircevent"
	"log"
        "net/http"
	"strings"
)

func launchhttpserver(bindto string){
        var err error
        err = http.ListenAndServe(bindto, nil)
        if err != nil {
                log.Fatal("Error starting http server: ", err)
        }
}

func main() {
	var myircbot *irc.Connection
	var botname *string
	var serveraddress *string
	var rawchannellist *string
	var channellist []string
	var channelname string
	var err error
	var extractor *URLTitleExtractor
        var github *GitHubAdapter
        var bindto *string

	botname = flag.String("botname", "justanotherbot", "Name of the bot")
	serveraddress = flag.String("server-address", "irc.freenode.org:6667", "Server Address")
	rawchannellist = flag.String("channels", "#ancient-solutions", "List of channels")
        bindto = flag.String("bind-to", ":8080", "IP:Port pair to bind the http-server to")
	flag.Parse()

	channellist = strings.Split(*rawchannellist, ",")

	myircbot = irc.IRC(*botname, *botname)
	if err = myircbot.Connect(*serveraddress); err != nil {
		log.Fatal("Error connecting to server: ", err)
	}

	extractor = &URLTitleExtractor{
		ircobject: myircbot,
	}
        github = &GitHubAdapter{
                ircbot: myircbot,
        }

	//Join all channels.
	for _, channelname = range channellist {
		myircbot.Join(channelname)
	}

	//Event handling
	myircbot.AddCallback("PRIVMSG", logprivmsgs)
	myircbot.AddCallback("PRIVMSG", extractor.WriteURLTitle)

        http.Handle("/github", github)

        // Start http server in a new thread
        go launchhttpserver(*bindto)

	myircbot.Loop()
}
