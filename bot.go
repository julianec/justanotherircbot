package main

import (
	"flag"
	"github.com/thoj/go-ircevent"
	"log"
	"strings"
	//        "code.google.com/p/go.net/html"
)

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
