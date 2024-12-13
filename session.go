package copilot

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// SessionInfo represents the context of the chat session based on the copilot
// references in the chat messages
type SessionInfo struct {
	URL         *ReferenceDataGitHubCurrentUrl `json:"url,omitempty"`
	Issue       *Issue                         `json:"issue,omitempty"`
	PullRequest *PullRequest                   `json:"pull_request,omitempty"`
	Repo        *ReferenceDataGitHubRepository `json:"repo,omitempty"`
	Agent       *ReferenceDataGitHubAgent      `json:"agent,omitempty"`
}

// IsSessionMessage returns true if the message has the role of "system" and
// the name of "_session"
func (msg *Message) IsSessionMessage() bool {
	return msg.Name == "_session"
}

// GetSessionMessage returns the message with the role of "system" and the name of
// "_session" which is the message that the github.com chat interface sends to
// communicate the current url and other context of the chat session
func (req *Request) GetSessionMessage() *Message {
	// iterate over the messages in reverse order
	for i := len(req.Messages) - 1; i >= 0; i-- {
		msg := req.Messages[i]
		if msg.Role == ChatRoleSystem {
			if msg.IsSessionMessage() {
				return msg
			}
		}
	}
	return nil
}

// GetSessionInfo returns the context of the chat session, including the
// current url, the relevant repository, agent details, and the associated issue
// or pull request if the current url is a valid issue or pull request url
func (req *Request) GetSessionInfo() (*SessionInfo, error) {
	// iterate over the messages in reverse order
	// var session *Message
	var url *ReferenceDataGitHubCurrentUrl
	var item *repoItemRef
	var repo *ReferenceDataGitHubRepository
	var agent *ReferenceDataGitHubAgent

	fmt.Println("resolving session info")

	for i := len(req.Messages) - 1; i >= 0; i-- {
		msg := req.Messages[i]

		switch msg.Role {
		case ChatRoleUser:
			if repo == nil {
				for _, ref := range msg.References {
					switch data := ref.Data.(type) {
					case *ReferenceDataGitHubRepository:
						repo = data
					}

					if ref.Type == ReferenceTypeGitHubRepository {
						repo = ref.Data.(*ReferenceDataGitHubRepository)
					}
				}
			}

		case ChatRoleAssistant:
			if agent == nil {
				for _, ref := range msg.References {
					switch data := ref.Data.(type) {
					case *ReferenceDataGitHubAgent:
						agent = data
					}
				}
			}

		case ChatRoleSystem:
			// resolve the session url context
			if msg.IsSessionMessage() {
				// session = msg
				if urlData := cetCurrentURLData(msg); urlData != nil {
					url = urlData
					if itemRefData, err := resolveRepoItemRef(url.URL); err == nil {
						item = itemRefData
					}
				}
			}

			if repo == nil {
				for _, ref := range msg.References {
					switch data := ref.Data.(type) {
					case *ReferenceDataGitHubRepository:
						repo = data
					}

					if ref.Type == ReferenceTypeGitHubRepository {
						repo = ref.Data.(*ReferenceDataGitHubRepository)
					}
				}
			}

		default:
			continue
		}
	}

	if url == nil {
		// this will happen if the user is not using the web (dotcom)
		// chat interface, or if the current url reference is redacted
		fmt.Println(" - no session url context found")
	}

	if item == nil {
		// item will ONLY have a value (potentially) if url has a value
		// and that url value is a valid issue or pull request url
		fmt.Println(" - no session issue or pull request context found")
	}

	if repo == nil {
		// repo may be nil if the user is not interacting in the context
		// of a repository, or if the current repository reference is redacted
		fmt.Println(" - no session repo context found")
	}

	if agent == nil {
		// we won't have an agent reference on the first message
		// from the user, so we'll create one with the agent login
		agent = &ReferenceDataGitHubAgent{
			Login: req.Agent,
			Type:  string(ReferenceTypeGitHubAgent),
			URL:   fmt.Sprintf("https://github.com/apps/%s", req.Agent),
		}
	}

	// do a little validation

	if !strings.EqualFold(agent.Login, req.Agent) {
		return nil, fmt.Errorf("warning: agent login %s does not match session agent login %s", req.Agent, agent.Login)
	}

	if url != nil {
		if item != nil {
			if url.Owner != "" && item.Owner != "" && !strings.EqualFold(url.Owner, item.Owner) {
				return nil, fmt.Errorf("warning: session url owner %s does not match item ref owner %s", url.Owner, item.Owner)
			}
			if url.Repo != "" && item.Repo != "" && !strings.EqualFold(url.Repo, item.Repo) {
				return nil, fmt.Errorf("warning: session url repo %s does not match item ref repo %s", url.Repo, item.Repo)
			}
		}
		if repo != nil {
			if url.Owner != "" && repo.OwnerLogin != "" && !strings.EqualFold(url.Owner, repo.OwnerLogin) {
				return nil, fmt.Errorf("warning: session url owner %s does not match repo owner %s", url.Owner, repo.OwnerLogin)
			}
			if url.Repo != "" && repo.Name != "" && !strings.EqualFold(url.Repo, repo.Name) {
				return nil, fmt.Errorf("warning: session url repo %s does not match repo name %s", url.Repo, repo.Name)
			}
		}
	}

	c := &SessionInfo{
		// Item:  item,
		URL:   url,
		Repo:  repo,
		Agent: agent,
	}

	if item != nil {
		switch item.Type {
		case RepoItemRefTypeIssue:
			c.Issue = &Issue{repoItemRef: *item}
		case RepoItemRefTypePull:
			c.PullRequest = &PullRequest{repoItemRef: *item}
		}
	}

	return c, nil
}

// cetCurrentURLData returns the current url reference data from the _session message
func cetCurrentURLData(msg *Message) *ReferenceDataGitHubCurrentUrl {
	// we only care about the current url reference
	// on the _session message that dotcom sends
	if !msg.IsSessionMessage() {
		return nil
	}

	// _session message indicates a web (dotcom) chat session
	// and should have a reference of type "github.current-url"
	for _, ref := range msg.References {
		switch d := ref.Data.(type) {
		case *ReferenceDataGitHubCurrentUrl:
			return d
		case *ReferenceDataGitHubRedacted:
			// the current URL reference may be redacted
			if d.Type == ReferenceTypeGitHubCurrentUrl {
				fmt.Println("warning: current URL reference is redacted")
				return nil
			}
		}
	}

	return nil
}

// resolveRepoItemRef resolves the owner, repo, and number from a github issue
// or pull request url
func resolveRepoItemRef(url string) (*repoItemRef, error) {
	var i = &repoItemRef{}
	var re *regexp.Regexp
	var matches []string

	if matches = issueRe.FindStringSubmatch(url); matches != nil {
		re = issueRe
		i.Type = RepoItemRefTypeIssue
	} else if matches = pullRe.FindStringSubmatch(url); matches != nil {
		re = pullRe
		i.Type = RepoItemRefTypePull
	} else {
		return nil, nil
	}

	o := re.SubexpIndex("owner")
	r := re.SubexpIndex("repo")
	n := re.SubexpIndex("number")
	p := re.SubexpIndex("page")
	h := re.SubexpIndex("hash")

	if o == -1 || r == -1 || n == -1 {
		return nil, fmt.Errorf("invalid %s url: %s", i.Type.Singular, url)
	}

	i.Owner = matches[o]
	i.Repo = matches[r]

	// convert the number to an int
	num, err := strconv.Atoi(matches[n])
	if err != nil {
		return nil, fmt.Errorf("invalid %s url: %s", i.Type.Singular, url)
	}
	i.Number = num

	if p != -1 {
		i.Page = matches[p]
	}

	if h != -1 {
		i.Hash = matches[h]
	}

	if idx := re.FindStringSubmatchIndex(url); idx != nil {
		urlTemplate := fmt.Sprintf("https://github.com/$owner/$repo/%s/$number", i.Type.Plural)
		urlBytes := []byte{}
		urlBytes = re.ExpandString(urlBytes, urlTemplate, url, idx)
		i.URL = string(urlBytes)

		apiTemplate := fmt.Sprintf("https://api.github.com/repos/$owner/$repo/%s/$number", i.Type.Plural)
		apiBytes := []byte{}
		apiBytes = re.ExpandString(apiBytes, apiTemplate, url, idx)
		i.API = string(apiBytes)
	}

	return i, nil
}

// repoItemRef represents a reference to a github issue or pull request
type repoItemRef struct {
	Type   RepoItemRefType `json:"type"`
	Owner  string          `json:"owner"`
	Repo   string          `json:"repo"`
	Number int             `json:"number"`
	Page   string          `json:"page,omitempty"`
	Hash   string          `json:"hash,omitempty"`
	URL    string          `json:"url"`
	API    string          `json:"api"`
}

// Issue represents a github issue
type Issue struct {
	repoItemRef
}

// PullRequest represents a github pull request
type PullRequest struct {
	repoItemRef
}

var (
	issueRe = regexp.MustCompile(`https://github.com/(?P<owner>[^/]+)/(?P<repo>[^/]+)/issues/(?P<number>\d+)(?:#(?P<hash>.+))?`)
	pullRe  = regexp.MustCompile(`https://github.com/(?P<owner>[^/]+)/(?P<repo>[^/]+)/pull/(?P<number>\d+)(?:/(?P<page>commits|checks|files))?(?:#(?P<hash>.+))?`)
)

type RepoItemRefType struct {
	Singular string
	Plural   string
}

var (
	RepoItemRefTypeIssue = RepoItemRefType{"issue", "issues"}
	RepoItemRefTypePull  = RepoItemRefType{"pull", "pulls"}
)

func (t *RepoItemRefType) UnmarshalJSON(data []byte) error {
	var singular string
	if err := json.Unmarshal(data, &singular); err != nil {
		return err
	}

	if singular == RepoItemRefTypeIssue.Singular {
		t.Plural = RepoItemRefTypeIssue.Plural
	} else if singular == RepoItemRefTypePull.Singular {
		t.Plural = RepoItemRefTypePull.Plural
	} else {
		return fmt.Errorf("invalid repo item ref type: %s", singular)
	}

	t.Singular = singular
	return nil
}

func (t RepoItemRefType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Singular)
}
