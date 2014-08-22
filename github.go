package main

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"github.com/thoj/go-ircevent"
	"hash"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type GitHubAdapter struct {
	ircbot   *irc.Connection
	config   *GitHubConfig
	channels map[string]*GitHubRepositoryConfig
}

// Every GithubEvent has a method Strings that returns a slice of string with the text for the irc PRIVMSG.
type GithubEvent interface {
        Strings() []string
        GetRepository() string
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

type GithubDelete struct {
        Ref             string
        RefType         string `json:"ref_type"`
        PusherType      string
        Repository      GithubRepository
        Sender          GithubUser
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

// CheckMAC returns true if messageMAC is a valid HMAC tag for message.
func CheckMAC(message []byte, messageMAC string, key string) bool {
	var err error
	var mac hash.Hash
	var macdata []byte
	var macparts = strings.Split(messageMAC, "=")
	macdata, err = hex.DecodeString(macparts[1])
	if err != nil {
		log.Print("Error decoding hex digest: ", err)
		return false
	}
	switch macparts[0] {
	case "md5":
		mac = hmac.New(md5.New, []byte(key))
	case "sha1":
		mac = hmac.New(sha1.New, []byte(key))
	case "sha256":
		mac = hmac.New(sha256.New, []byte(key))
	case "sha512":
		mac = hmac.New(sha512.New, []byte(key))
	default:
		log.Print("Unsupported hash: ", macparts[0])
		return false
	}
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(macdata, expectedMAC)
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
	var lines []string = strings.Split(g.Message, "\n")                // Commit message
	var text string 
        text = g.Author.String() + " \x02" + g.Id[0:7] + "\x0f" // First 7 characters

	if len(g.Added) > 0 {
		text += " \x0303" + strings.Join(g.Added, " ") + "\x0f"
	}

	if len(g.Removed) > 0 {
		text += " \x0304" + strings.Join(g.Removed, " ") + "\x0f"
	}

	if len(g.Modified) > 0 {
		text += " \x0310" + strings.Join(g.Modified, " ") + "\x0f"
	}
	if len(lines) > 0 {
		text += " " + lines[0]
	}
	return text
}

func (g *GithubPush) Strings() []string {
	var refs []string = strings.Split(g.Ref, "/")
	var prefix string = "\x0303" + g.Repository.String() + "\x0f \x0305" + refs[len(refs)-1] + "\x0f"
	var pushes []string = make([]string, 0)

	for _, commit := range g.Commits {
		pushes = append(pushes, prefix+" "+commit.String())
		pushes = append(pushes, prefix+" "+commit.Url)
	}
	return pushes
}

func (g *GithubPush) GetRepository() string {
       return g.Repository.String()
}


func (g *GithubCreate) Strings() []string {
	return []string{"\x0303" + g.Sender.String() + "\x0f has pushed a new " + g.RefType + " \x0305" + g.Ref + "\x0f to \x0303" + g.Repository.String() + "\x0f"}
}

func (g *GithubCreate) GetRepository() string {
       return g.Repository.String()
}

func (g *GithubDelete) Strings() []string {
	return []string{"\x0303" + g.Sender.String() + "\x0f has deleted a " + g.RefType + " \x0305" + g.Ref + "\x0f from \x0303" + g.Repository.String() + "\x0f"}
}

func (g *GithubDelete) GetRepository() string {
       return g.Repository.String()
}

func (g *GitHubAdapter) WriteGithubEvent(event GithubEvent, body []byte, signature string) error {
        var ok bool
        var githubconf *GitHubRepositoryConfig
        var err error

        err = json.Unmarshal(body, &event)

        if err != nil {
                log.Print("Error decoding github event: ", err)
                return err
        }

        githubconf, ok = g.channels[event.GetRepository()]
        if !ok {
                log.Print("Repository ", event.GetRepository(), " not configured.")
                return err
        }
        if !CheckMAC(body, signature, githubconf.GetSecret()) {
                log.Print("DEBUG Spam, spam spam")
                return err
        }
        for _, channel := range githubconf.GetIrcChannel() {
                for _, commit := range event.Strings() {
                        g.ircbot.Privmsg(channel, commit)
                }
        }
        return err
}

func NewGitHubAdapter(ircbot *irc.Connection, config *GitHubConfig) *GitHubAdapter {
	var channels = make(map[string]*GitHubRepositoryConfig)
	for _, repo := range config.Repo {
		channels[repo.GetName()] = repo
	}
	return &GitHubAdapter{
		ircbot:   ircbot,
		config:   config,
		channels: channels,
	}
}

func (g *GitHubAdapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var body []byte
	var err error
	body, err = ioutil.ReadAll(req.Body)
	if err != nil {
		log.Print("Error reading body: ", err)
		return
	}

	switch req.Header.Get("X-GitHub-Event") {
	case "create":
                var create GithubCreate
                g.WriteGithubEvent(&create, body, req.Header.Get("X-Hub-Signature"))
        case "delete":
                var del GithubDelete
                g.WriteGithubEvent(&del, body, req.Header.Get("X-Hub-Signature"))
	case "push":
		var push GithubPush
                g.WriteGithubEvent(&push, body, req.Header.Get("X-Hub-Signature"))
	default:
		log.Print("Unknown GitHub event.", req.Header.Get("X-GitHub-Event"))
	}
}
