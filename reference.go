package copilot

import (
	"encoding/json"
	"regexp"
)

// ReferenceType is the type of copilot_reference.
type ReferenceType string

const (
	// ReferenceTypeGitHubRedacted is another reference that has been redacted.
	ReferenceTypeGitHubRedacted ReferenceType = "github.redacted"
	// ReferenceTypeGitHubAgent is a GitHub Copilot agent.
	ReferenceTypeGitHubAgent ReferenceType = "github.agent"
	// ReferenceTypeGitHubCurrentUrl is the current URL of the user.
	ReferenceTypeGitHubCurrentUrl ReferenceType = "github.current-url"
	// ReferenceTypeGitHubFile is a file in a GitHub repository.
	ReferenceTypeGitHubFile ReferenceType = "github.file"
	// ReferenceTypeGitHubRepository is a GitHub repository.
	ReferenceTypeGitHubRepository ReferenceType = "github.repository"
	// ReferenceTypeGitHubSnippet is a code snippet from a file in a GitHub repository.
	ReferenceTypeGitHubSnippet ReferenceType = "github.snippet"
	// ReferenceTypeClientFile is a file from a client like vscode.
	ReferenceTypeClientFile ReferenceType = "client.file"
	// ReferenceTypeClientSelection is a selection in a file from a client like vscode.
	ReferenceTypeClientSelection ReferenceType = "client.selection"
)

// Reference is a copilot_reference in chat messages
type Reference struct {
	// Type is the type of the copilot references.
	Type ReferenceType `json:"type"`
	// ID is the unique identifier of the copilot references.
	ID string `json:"id"`
	// IsImplicit specifies if the reference was passed implicitly or explicitly.
	IsImplicit bool `json:"is_implicit"`
	// Metadata includes any metadata to display in the user's environment.
	// If any of the Metadata required fields are missing, the reference will
	// not be rendered in the UI.
	Metadata ReferenceMetadata `json:"metadata"`
	// Data that is specific to the copilot references.
	Data    ReferenceData   `json:"-"`
	RawData json.RawMessage `json:"data"`
}

// ReferenceMetadata contains metadata about a copilot_reference to display in
// the user's environment. If any of the required fields are missing, the reference
// will not be rendered in the UI.
type ReferenceMetadata struct {
	DisplayName string `json:"display_name"`
	DisplayIcon string `json:"display_icon"`
	DisplayURL  string `json:"display_url"`
}

// ReferenceData is the data that is specific to the copilot_reference type.
type ReferenceData interface{}

// ReferenceDataOther contains information about a copilot_reference that is not
// included in the well-known types.
type ReferenceDataOther struct {
	Type string `json:"type"`
}

// ReferenceDataGitHubRedacted contains information about a redacted copilot_reference.
// Included on copilot references of type: "github.redacted"
type ReferenceDataGitHubRedacted struct {
	// Type is the type of the copilot references that was redacted.
	Type ReferenceType `json:"type"`
}

// ReferenceDataGitHubAgent contains information about a GitHub Copilot agent.
// Included on copilot references of type: "github.agent"
type ReferenceDataGitHubAgent struct {
	AvatarURL string `json:"avatarURL"`
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Type      string `json:"type" default:"github.agent"`
	URL       string `json:"url"`
}

// ReferenceDataGitHubCurrentUrl contains information about the users current URL.
// Included on copilot references of type: "github.current-url"
//
// Note: this is a special copilot_reference that only seems to be present in a
// message with role: "user" and name: "_session"
//
// Only the URL field is sent by the client, the rest of the fields are extracted
// during unmarshaling.
//
// "github.current-url" is ONLY sent by the github.com (web) chat client.
type ReferenceDataGitHubCurrentUrl struct {
	URL   string `json:"url"`
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
	Path  string `json:"path,omitempty"`
	Hash  string `json:"hash,omitempty"`
}

// ReferenceDataGitHubRepository contains information about a github repository.
// Included on copilot references of type: "github.repository"
//
// Note: fields set by different clients vary, and are not guaranteed to be present.
// ID, OwnerLogin, Name, and Type do seem to be present across clients.
//
// "github.repository" is sent by the vscode and github.com (web) chat clients.
type ReferenceDataGitHubRepository struct {
	CommitOID   string                                   `json:"commitOID"`
	Description string                                   `json:"description"`
	ID          int64                                    `json:"id"`
	Languages   []*ReferenceDataGitHubRepositoryLanguage `json:"languages"`
	Name        string                                   `json:"name"`
	OwnerLogin  string                                   `json:"ownerLogin"`
	OwnerType   string                                   `json:"ownerType"`
	ReadmePath  string                                   `json:"readmePath"`
	Ref         string                                   `json:"ref"`
	RefInfo     ReferenceDataGitHubRepositoryRefInfo     `json:"refInfo"`
	Type        string                                   `json:"type" default:"repository"`
	Visibility  string                                   `json:"visibility"`
}
type ReferenceDataGitHubRepositoryLanguage struct {
	Name    string  `json:"name"`
	Percent float64 `json:"percent"`
}
type ReferenceDataGitHubRepositoryRefInfo struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// ReferenceDataGitHubFile contains information about a file in a github repository.
// Included on copilot references of type: "github.file"
//
// "github.file" is ONLY sent by github.com (web) chat client.
type ReferenceDataGitHubFile struct {
	CommitOID    string `json:"commitOID"`
	LanguageID   int64  `json:"languageID"`
	LanguageName string `json:"languageName"`
	Path         string `json:"path"`
	Ref          string `json:"ref"`
	RepoID       int64  `json:"repoID"`
	RepoName     string `json:"repoName"`
	RepoOwner    string `json:"repoOwner"`
	Type         string `json:"type" default:"file"`
	URL          string `json:"url"`
}

// ReferenceDataGitHubSnippet contains information about a code snippet from a
// file in a github repository.
// Included on copilot references of type: "github.snippet"
//
// "github.snippet" is ONLY sent by github.com (web) chat client.
type ReferenceDataGitHubSnippet struct {
	CommitOID    string                          `json:"commitOID"`
	LanguageID   int64                           `json:"languageID"`
	LanguageName string                          `json:"languageName"`
	Path         string                          `json:"path"`
	Range        ReferenceDataGitHubSnippetRange `json:"range"`
	Ref          string                          `json:"ref"`
	RepoID       int64                           `json:"repoID"`
	RepoName     string                          `json:"repoName"`
	RepoOwner    string                          `json:"repoOwner"`
	Type         string                          `json:"type" default:"snippet"`
	URL          string                          `json:"url"`
}
type ReferenceDataGitHubSnippetRange struct {
	End   int32 `json:"end"`
	Start int32 `json:"start"`
}

// ReferenceDataClientFile contains information about a specific file on a client.
// Included on copilot references of type: "client.file"
//
// "client.file" is sent by the vscode chat client (and likely visual studio).
type ReferenceDataClientFile struct {
	Content  string `json:"content"`
	Language string `json:"language"`
}

// ReferenceDataClientSelection contains information about a selection in a file.
// Included on copilot references of type: "client.selection"
//
// "client.selection" is sent by the vscode chat client (and likely visual studio).
type ReferenceDataClientSelection struct {
	Content string                               `json:"content"`
	End     ReferenceDataClientSelectionLocation `json:"end"`
	Start   ReferenceDataClientSelectionLocation `json:"start"`
}
type ReferenceDataClientSelectionLocation struct {
	Col  int32 `json:"col"`
	Line int32 `json:"line"`
}

func (r *Reference) UnmarshalJSON(data []byte) error {
	type reference Reference

	if err := json.Unmarshal(data, (*reference)(r)); err != nil {
		return err
	}

	var d ReferenceData

	switch r.Type {
	case ReferenceTypeGitHubRedacted:
		d = &ReferenceDataGitHubRedacted{}
	case ReferenceTypeGitHubAgent:
		d = &ReferenceDataGitHubAgent{}
	case ReferenceTypeGitHubCurrentUrl:
		d = &ReferenceDataGitHubCurrentUrl{}
	case ReferenceTypeGitHubFile:
		d = &ReferenceDataGitHubFile{}
	case ReferenceTypeGitHubRepository:
		d = &ReferenceDataGitHubRepository{}
	case ReferenceTypeGitHubSnippet:
		d = &ReferenceDataGitHubSnippet{}
	case ReferenceTypeClientFile:
		d = &ReferenceDataClientFile{}
	case ReferenceTypeClientSelection:
		d = &ReferenceDataClientSelection{}
	default:
		d = &ReferenceDataOther{}
	}

	if err := json.Unmarshal(r.RawData, d); err != nil {
		return err
	}

	r.Data = d

	return nil
}

var repoRe = regexp.MustCompile(`https://github.com/(?:orgs/)?(?P<owner>[^/]+)/(?P<repo>[^/]+)?(?P<path>(?:/(?:[^/#]+))+)?(?:#(?P<hash>.+))?`)

func (u *ReferenceDataGitHubCurrentUrl) UnmarshalJSON(data []byte) error {
	type referenceDataGitHubCurrentUrl ReferenceDataGitHubCurrentUrl

	if err := json.Unmarshal(data, (*referenceDataGitHubCurrentUrl)(u)); err != nil {
		return err
	}

	if matches := repoRe.FindStringSubmatch(u.URL); matches != nil {
		if o := repoRe.SubexpIndex("owner"); o > -1 {
			u.Owner = matches[o]
		}
		if r := repoRe.SubexpIndex("repo"); r > -1 {
			u.Repo = matches[r]
		}
		if p := repoRe.SubexpIndex("path"); p > -1 {
			u.Path = matches[p]
		}
		if h := repoRe.SubexpIndex("hash"); h > -1 {
			u.Hash = matches[h]
		}
	}

	return nil
}
