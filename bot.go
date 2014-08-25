package main

import (
	"code.google.com/p/goprotobuf/proto"
	"flag"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"log"
	"net/http"
)

func launchhttpserver(bindto string) {
	var err error
	err = http.ListenAndServe(bindto, nil)
	if err != nil {
		log.Fatal("Error starting http server: ", err)
	}
}

func main() {
	var myircbot *irc.Connection
	var channelname string
	var err error
	var extractor *URLTitleExtractor
	var github *GitHubAdapter
	var configpath string // path of config file
	var configdata []byte
	var config IRCBotConfig
	var msgbuffer *MessageBuffer

	flag.StringVar(&configpath, "config", "", "Specify the path to the configuration file.")
	flag.Parse()

	configdata, err = ioutil.ReadFile(configpath)
	if err != nil {
		log.Fatal("Error reading config file: ", err)
	}

	err = proto.Unmarshal(configdata, &config)
	if err != nil {
		err = proto.UnmarshalText(string(configdata), &config)
	}
	if err != nil {
		log.Fatal("Error parsing config: ", err)
	}

	myircbot = irc.IRC(config.GetBotName(), config.GetBotName())
	if err = myircbot.Connect(config.GetServerAddress()); err != nil {
		log.Fatal("Error connecting to server: ", err)
	}
	msgbuffer = NewMessageBuffer(myircbot, config.GetSendQueueLength())

	extractor = &URLTitleExtractor{
		msgbuffer: msgbuffer,
	}
	github = NewGitHubAdapter(msgbuffer, config.GetGithub())

	//Join all channels.
	for _, channelname = range config.GetIrcChannel() {
		myircbot.Join(channelname)
	}

	//Event handling
	myircbot.AddCallback("PRIVMSG", logprivmsgs)
	myircbot.AddCallback("PRIVMSG", extractor.WriteURLTitle)

	// Write GitHub status messages to the specified channels
	http.Handle("/github", github)

	// Start http server in a new thread
	go launchhttpserver(config.GetHttpServerAddress())

	myircbot.Loop()
}
