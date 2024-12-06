package agent

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/colbylwilliams/copilot-go"
	"github.com/colbylwilliams/copilot-go/_examples/github/github"
	"github.com/colbylwilliams/copilot-go/sse"
	"github.com/go-chi/chi/v5/middleware"
	gh "github.com/google/go-github/v67/github"
)

type MyAgent struct{ cfg *copilot.Config }

func NewAgent(cfg *copilot.Config) *MyAgent {
	return &MyAgent{cfg: cfg}
}

func (a *MyAgent) Execute(ctx context.Context, token string, req *copilot.Request, w http.ResponseWriter) error {

	// write the sse headers
	sse.WriteStreamingHeaders(w)

	// get the request reqId from the context
	reqId := middleware.GetReqID(ctx)
	if reqId == "" {
		return errors.New("request id not found in context")
	}

	// the appClient is authenticated as the app itself and can only be used to
	// get information about the app, its installations, and other app-level data.
	appClient, err := github.NewAppClient(ctx, a.cfg)
	if err != nil {
		return err
	}

	app, _, err := appClient.Apps.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("error getting app: %w", err)
	}

	sse.WriteDelta(w, reqId, "# App\n")

	sse.WriteDelta(w, reqId, "|  |  |\n")
	sse.WriteDelta(w, reqId, "| --- | --- |\n")
	sse.WriteDelta(w, reqId, fmt.Sprintf("| ID   | %d |\n", app.GetID()))
	sse.WriteDelta(w, reqId, fmt.Sprintf("| Name | %s |\n", app.GetName()))
	sse.WriteDelta(w, reqId, fmt.Sprintf("| Slug | %s |\n", app.GetSlug()))
	sse.WriteDelta(w, reqId, fmt.Sprintf("| Owner| %s |\n", app.GetOwner().GetLogin()))

	installs, _, err := appClient.Apps.ListInstallations(context.Background(), &gh.ListOptions{PerPage: 100})
	if err != nil {
		return fmt.Errorf("error listing installations: %w", err)
	}

	sse.WriteDelta(w, reqId, "#### Installations\n")
	sse.WriteDelta(w, reqId, "| Login | ID |\n")
	sse.WriteDelta(w, reqId, "| --- | --- |\n")

	for _, install := range installs {
		sse.WriteDelta(w, reqId, fmt.Sprintf("| %s | %d |\n", install.GetAccount().GetLogin(), install.GetID()))
	}

	// The client is authenticated using the token provided by github when it called
	// the app. This token's permissions are a union of the user's permissions and
	// the app's permissions and is only valid for a short period of time.
	client := gh.NewClient(nil).WithAuthToken(token)

	// get the current user
	me, _, err := client.Users.Get(ctx, "")
	if err != nil {
		return fmt.Errorf("error getting current user: %w", err)
	}

	sse.WriteDelta(w, reqId, "# User\n")
	sse.WriteDelta(w, reqId, "|  |  |\n")
	sse.WriteDelta(w, reqId, "| --- | --- |\n")
	sse.WriteDelta(w, reqId, fmt.Sprintf("| Login | %s |\n", me.GetLogin()))
	sse.WriteDelta(w, reqId, fmt.Sprintf("| Name  | %s |\n", me.GetName()))

	// The session is populated by the copilot middleware and contains information
	// about the current user's chat session, for example the relevant repository.
	//
	// If the user is chatting with the app on github.com (web) UI, this session
	// will include the user's current.
	//
	// Additionally, if the current url is a github issue or pull request, the
	// session will include details about the issue or pull request.
	session := copilot.GetSessionInfo(ctx)
	if session == nil {
		return errors.New("session not found in context")
	}

	sse.WriteDelta(w, reqId, "# Session\n")

	if session.URL != nil {
		sse.WriteDelta(w, reqId, "#### Current URL\n")
		sse.WriteDelta(w, reqId, "|  |  |\n")
		sse.WriteDelta(w, reqId, "| --- | --- |\n")
		sse.WriteDelta(w, reqId, fmt.Sprintf("| URL   | %s |\n", session.URL.URL))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Owner | %s |\n", session.URL.Owner))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Repo  | %s |\n", session.URL.Repo))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Path  | %s |\n", session.URL.Path))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Hash  | %s |\n", session.URL.Hash))
	}

	if session.Agent != nil {
		sse.WriteDelta(w, reqId, "#### Agent\n")
		sse.WriteDelta(w, reqId, "|  |  |\n")
		sse.WriteDelta(w, reqId, "| --- | --- |\n")
		sse.WriteDelta(w, reqId, fmt.Sprintf("| ID    | %d |\n", session.Agent.ID))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Login | %s |\n", session.Agent.Login))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| URL   | %s |\n", session.Agent.URL))
	}

	if session.Repo != nil {
		sse.WriteDelta(w, reqId, "#### Repository\n")
		sse.WriteDelta(w, reqId, "|  |  |\n")
		sse.WriteDelta(w, reqId, "| --- | --- |\n")
		sse.WriteDelta(w, reqId, fmt.Sprintf("| ID         | %d |\n", session.Repo.ID))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Name       | %s |\n", session.Repo.Name))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| OwnerLogin | %s |\n", session.Repo.OwnerLogin))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| OwnerType  | %s |\n", session.Repo.OwnerType))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Visibility | %s |\n", session.Repo.Visibility))
	}

	if session.Issue != nil {
		sse.WriteDelta(w, reqId, "#### Issue\n")
		sse.WriteDelta(w, reqId, "|  |  |\n")
		sse.WriteDelta(w, reqId, "| --- | --- |\n")
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Number | [#%d](%s) |\n", session.Issue.Number, session.Issue.URL))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Repo   | %s |\n", session.Issue.Repo))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Owner  | %s |\n", session.Issue.Owner))
	}

	if session.PullRequest != nil {
		sse.WriteDelta(w, reqId, "#### Pull Request\n")
		sse.WriteDelta(w, reqId, "|  |  |\n")
		sse.WriteDelta(w, reqId, "| --- | --- |\n")
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Number | [#%d](%s) |\n", session.PullRequest.Number, session.PullRequest.URL))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Repo   | %s |\n", session.PullRequest.Repo))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Owner  | %s |\n", session.PullRequest.Owner))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Page   | %s |\n", session.PullRequest.Page))
	}

	// if issue or pull request is present, use the github client to get the details
	if session.Issue != nil {
		issue, _, err := client.Issues.Get(ctx, session.Repo.OwnerLogin, session.Repo.Name, session.Issue.Number)
		if err != nil {
			return fmt.Errorf("error getting issue: %w", err)
		}

		sse.WriteDelta(w, reqId, "# Issue\n")
		sse.WriteDelta(w, reqId, "|  |  |\n")
		sse.WriteDelta(w, reqId, "| --- | --- |\n")
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Title | %s |\n", issue.GetTitle()))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Body  | %s |\n", issue.GetBody()))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| State | %s |\n", issue.GetState()))
	}

	if session.PullRequest != nil {
		pr, _, err := client.PullRequests.Get(ctx, session.Repo.OwnerLogin, session.Repo.Name, session.PullRequest.Number)
		if err != nil {
			return fmt.Errorf("error getting pull request: %w", err)
		}

		sse.WriteDelta(w, reqId, "# Pull Request\n")
		sse.WriteDelta(w, reqId, "|  |  |\n")
		sse.WriteDelta(w, reqId, "| --- | --- |\n")
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Title | %s |\n", pr.GetTitle()))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| Body  | %s |\n", pr.GetBody()))
		sse.WriteDelta(w, reqId, fmt.Sprintf("| State | %s |\n", pr.GetState()))
	}

	sse.WriteStop(w, reqId)

	return nil
}
