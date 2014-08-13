package main

import(
        "flag"
        "github.com/thoj/go-ircevent"
        "log"
        "strings"
        "regexp"
)

//var url_re = regexp.MustCompile(`\b(?:ht|f)tps?://[a-zA-Z0-9%_\+-]+\?[a-zA-Z0-9%_\+-=&]+\b`)
//var url_re = regexp.MustCompile(`\b(?:ht|f)tps?://[a-zA-Z0-9%_\+-]+\?[a-zA-Z0-9%_\+-=&]+`)
var url_re = regexp.MustCompile(`\bhttps?://[^\s]+\b`)

func logprivmsgs(event *irc.Event){
        log.Print(event.Nick + ": ", event.Arguments)
}
func writeurltitle(event *irc.Event){
        //var url string
        log.Print(url_re.FindAllString(event.Arguments[1], -1))
}

func FindURLs(input string) []string {
        return url_re.FindAllString(input, -1)
}

func main(){
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
