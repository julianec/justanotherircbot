package main

import (
	"encoding/json"
	"github.com/thoj/go-ircevent"
	"log"
	"net/http"
	"strings"
)

type GitHubAdapter struct {
	ircbot   *irc.Connection
	config   *GitHubConfig
	channels map[string]*GitHubRepositoryConfig
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
	var text string = g.Author.String() + " \x02" + g.Id[0:7] + "\x0f" // First 7 characters

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

func (g *GithubCreate) String() string {
	return "\x0303" + g.Sender.String() + "\x0f has pushed a new " + g.RefType + " \x0305" + g.Ref + "\x0f to \x0303" + g.Repository.String() + "\x0f"
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
	var jsondecoder = json.NewDecoder(req.Body)
	switch req.Header.Get("X-GitHub-Event") {
	case "create":
		var create GithubCreate
                var githubconf *GitHubRepositoryConfig
                var ok bool
		var err error
		err = jsondecoder.Decode(&create)
		if err != nil {
			log.Print("Error decoding github create: ", err)
			return
		}
		log.Print(create.String())
                githubconf, ok = g.channels[create.Repository.String()]
                if !ok {
                        log.Print("Repository ", create.Repository.String(), " not configured.")
                        return
                }
                for _, channel := range githubconf.GetIrcChannel() {
                        g.ircbot.Privmsg(channel, create.String())
                }
	case "push":
		var push GithubPush
                var ok bool
                var githubconf *GitHubRepositoryConfig
		var err error

		err = jsondecoder.Decode(&push)

		if err != nil {
			log.Print("Error decoding github push: ", err)
			return
		}

                githubconf, ok = g.channels[push.Repository.String()]
                if !ok {
                        log.Print("Repository ", push.Repository.String(), " not configured.")
                        return
                }
                for _, channel := range githubconf.GetIrcChannel() {
                        for _, commit := range push.Strings() {
                                g.ircbot.Privmsg(channel, commit)
                        }
                }
	default:
		log.Print("Unknown GitHub event.", req.Header.Get("X-GitHub-Event"))
	}
}