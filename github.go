package main
import(
        "net/http"
        "github.com/thoj/go-ircevent"
        "encoding/json"
        "log"
        "flag"
)

type GitHubAdapter struct {
        ircbot *irc.Connection
}

type GithubCreate struct{
        Ref string
        RefType string `json:"ref_type"`
        MasterBranch string
        Description string
        PusherType string
        Repository GithubRepository
        Sender GithubUser
}
type GithubRepository struct{
        Id uint64
        Name string
        FullName string `json:"full_name"`
        Owner GithubUser
        Private bool
        HtmlUrl string
        Description string
        Fork bool
        Url string
        ForksUrl string
        KeysUrl string
        CollaboratorsUrl string
        TeamsUrl string
        HooksUrl string
        IssueEventsUrl string
        EventsUrl string
        AssigneesUrl string
        BranchesUrl string
        TagsUrl string
        BlobsUrl string
        GitTagsUrl string
        GitRefsUrl string
        TreesUrl string
        StatusesUrl string
        LanguagesUrl string
        StargazersUrl string
        ContributorsUrl string
        SubscribersUrl string
        SubscriptionUrl string
        CommitsUrl string
        GitCommitsUrl string
        CommentsUrl string
        CompareUrl string
        MergesUrl string
        ArchiveUrl string
        DownloadsUrl string
        IssuesUrl string
        PullsUrl string
        MilestonesUrl string
        NotificationsUrl string
        LabelsUrl string
        ReleasesUrl string
        CreatedAt string
        UpdatedAt string
        PushedAt string
        GitUrl string
        SshUrl string
        CloneUrl string
        SvnUrl string
        Homepage string
        Size uint64
        StargazersCount uint64
        WatchersCount uint64
        Language string
        HasIssues bool
        HasDownloads bool
        HasWiki bool
        ForksCount uint64
        MirrorUrl string
        OpenIssuesCount uint64
        Forks uint64
        OpenIssues uint64
        Watchers uint64
        DefaultBranch string
}
type GithubUser struct{
        Login string
        Id uint64
        AvatarURL string
        GravatarId string
        Url string
        HtmlUrl string
        FollowersUrl string
        FollowingUrl string
        GistsUrl string
        StarredUrl string
        SubscriptionsUrl string
        OrganizationsUrl string
        ReposUrl string
        EventsUrl string
        ReceivedEventsUrl string
        Type string
        SiteAdmin bool
        Name string
        Email string
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

func (g *GithubCreate) String() string {
        return g.Sender.String() + " has pushed a new " + g.RefType + " " + g.Ref + " to " + g.Repository.String()
}

func (g *GitHubAdapter) ServeHTTP(rw http.ResponseWriter, req *http.Request){
        var jsondecoder = json.NewDecoder(req.Body)
        switch req.Header.Get("X-GitHub-Event") {
        case "create":
                var create GithubCreate
                var err error
                err = jsondecoder.Decode(&create)
                if err != nil {
                        log.Print("Error decoding github create: ", err)
                }
                log.Print(create.String())
                g.ircbot.Privmsg(githubchannel, create.String())
        default:
                log.Print("Unknown GitHub event.", req.Header.Get("X-GitHub-Event"))
        }
}
