package main

import (
	"encoding/json"
	"flag"
	"github.com/thoj/go-ircevent"
	"log"
	"net/http"
        "strings"
)

type GitHubAdapter struct {
	ircbot *irc.Connection
}

type GithubCreate struct {
	Ref          string
	RefType      string `json:"ref_type"`
	MasterBranch string
	Description  string
	PusherType   string
	Repository   GithubRepository
	Sender       GithubUser
}

type GithubRepository struct {
	Id               uint64
	Name             string
	FullName         string `json:"full_name"`
	Owner            GithubUser
	Private          bool
	HtmlUrl          string
	Description      string
	Fork             bool
	Url              string
	ForksUrl         string
	KeysUrl          string
	CollaboratorsUrl string
	TeamsUrl         string
	HooksUrl         string
	IssueEventsUrl   string
	EventsUrl        string
	AssigneesUrl     string
	BranchesUrl      string
	TagsUrl          string
	BlobsUrl         string
	GitTagsUrl       string
	GitRefsUrl       string
	TreesUrl         string
	StatusesUrl      string
	LanguagesUrl     string
	StargazersUrl    string
	ContributorsUrl  string
	SubscribersUrl   string
	SubscriptionUrl  string
	CommitsUrl       string
	GitCommitsUrl    string
	CommentsUrl      string
	CompareUrl       string
	MergesUrl        string
	ArchiveUrl       string
	DownloadsUrl     string
	IssuesUrl        string
	PullsUrl         string
	MilestonesUrl    string
	NotificationsUrl string
	LabelsUrl        string
	ReleasesUrl      string
	CreatedAt        string
	UpdatedAt        string
	PushedAt         string
	GitUrl           string
	SshUrl           string
	CloneUrl         string
	SvnUrl           string
	Homepage         string
	Size             uint64
	StargazersCount  uint64
	WatchersCount    uint64
	Language         string
	HasIssues        bool
	HasDownloads     bool
	HasWiki          bool
	ForksCount       uint64
	MirrorUrl        string
	OpenIssuesCount  uint64
	Forks            uint64
	OpenIssues       uint64
	Watchers         uint64
	DefaultBranch    string
}

type GithubUser struct {
	Login             string
	Id                uint64
	AvatarURL         string
	GravatarId        string
	Url               string
	HtmlUrl           string
	FollowersUrl      string
	FollowingUrl      string
	GistsUrl          string
	StarredUrl        string
	SubscriptionsUrl  string
	OrganizationsUrl  string
	ReposUrl          string
	EventsUrl         string
	ReceivedEventsUrl string
	Type              string
	SiteAdmin         bool
	Name              string
	Email             string
}

type GithubPush struct {
	Ref        string
	After      string
	Before     string
	Created    bool
	Deleted    bool
	Forced     bool
	Compare    string
	Commits    []GithubCommit
	HeadCommit GithubCommit `json:"head_commit"`
	Repository GithubRepository
	Pusher     GithubUser
}

type GithubCommit struct {
	Id        string
	Distinct  bool
	Message   string
	Timestamp string
	Url       string
	Author    GithubUser
	Committer GithubUser
	Added     []string
	Removed   []string
	Modified  []string
}

var githubchannel string

func init() {
	flag.StringVar(&githubchannel, "github-channel", "#ancient-solutions", "Channel to post github messages to.")
}

func (g *GithubUser) String() string {
	if len(g.Name) > 0 && len(g.Email) > 0 {
		return g.Name + " <" + g.Email + ">"
	} else if len(g.Name) > 0 {
		return g.Name
	} else if len(g.Email) > 0 {
		return g.Email
	}
	return g.Login
}

func (g *GithubRepository) String() string {
	if len(g.FullName) > 0 {
		return g.FullName
	}
	return g.Name
}

func (g *GithubCommit) String() string {
        var lines []string = strings.Split(g.Message, "\n") // Commit message
        var text string = g.Author.String() + " " + g.Id[0:7] // First 7 characters

        if len(g.Added) > 0 {
                text += " a: " + strings.Join(g.Added, " ")
        }

        if len(g.Removed) > 0 {
                text += " d: " + strings.Join(g.Removed, " ")
        }

        if len(g.Modified) > 0 {
                text += " m: " + strings.Join(g.Modified, " ")
        }
        if len(lines) > 0 {
                text += " " + lines[0]
        }
        return text
}

func (g *GithubPush) Strings() []string {
        var refs []string = strings.Split(g.Ref, "/")
        var prefix string = g.Repository.String()+ " " + refs[len(refs)-1]
        var pushes []string = make([]string, 0)

        for _, commit := range g.Commits{
                pushes = append(pushes, prefix + " " + commit.String())
                pushes = append(pushes, prefix + " " + commit.Url)
        }
        return pushes
}

func (g *GithubCreate) String() string {
	return g.Sender.String() + " has pushed a new " + g.RefType + " " + g.Ref + " to " + g.Repository.String()
}

func (g *GitHubAdapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var jsondecoder = json.NewDecoder(req.Body)
	switch req.Header.Get("X-GitHub-Event") {
	case "create":
		var create GithubCreate
		var err error
		err = jsondecoder.Decode(&create)
		if err != nil {
			log.Print("Error decoding github create: ", err)
                        return
		}
		log.Print(create.String())
		g.ircbot.Privmsg(githubchannel, create.String())
        case "push":
                var push GithubPush
                var err error

                err = jsondecoder.Decode(&push)

                if err != nil {
                        log.Print("Error decoding github push: ", err)
                        return
                }

                for _, commit := range push.Strings() {
                        g.ircbot.Privmsg(githubchannel, commit)
                }
	default:
		log.Print("Unknown GitHub event.", req.Header.Get("X-GitHub-Event"))
	}
}
