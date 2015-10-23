package main

import (
	"github.com/golang/protobuf/proto"
	"flag"
	"github.com/thoj/go-ircevent"
	"io/ioutil"
	"log"
	"net/http"
)

func logErrors(c chan error) {
	var err error

	for err = range c {
		log.Print("IRC error: ", err)
	}
}

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
	if myircbot == nil {
                log.Fatal("Error calling IRC(nick, user string) *Connection. Nick or User empty.")
        }
	go logErrors(myircbot.ErrorChan()) // collect irc errors and log

        // 1: RPL_WELCOME "Welcome to the Internet Relay Network
        // <nick>!<user>@<host>"
        myircbot.AddCallback("001", func(e *irc.Event) {
                //Join all channels.
                for _, channelname = range config.GetIrcChannel() {
                        myircbot.Join(channelname)
                }
        })

	if err = myircbot.Connect(config.GetServerAddress()); err != nil {
		log.Fatal("Error connecting to server: ", err)
	}
	msgbuffer = NewMessageBuffer(myircbot, config.GetSendQueueLength())

	extractor = &URLTitleExtractor{
		msgbuffer: msgbuffer,
	}
	github = NewGitHubAdapter(msgbuffer, config.GetGithub())

	//Event handling
	myircbot.AddCallback("PRIVMSG", logprivmsgs)
	myircbot.AddCallback("PRIVMSG", extractor.WriteURLTitle)

	// Write GitHub status messages to the specified channels
	http.Handle("/github", github)

	// Start http server in a new thread
	go launchhttpserver(config.GetHttpServerAddress())

	myircbot.Loop()
}
