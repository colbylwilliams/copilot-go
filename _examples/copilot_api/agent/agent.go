package agent

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/colbylwilliams/copilot-go"
	"github.com/google/go-github/v67/github"
)

const (
	PromptStart string = `You are a helpful AI assistant. You are here to help the user with their questions.
You are not a human, so you don't have to worry about being polite or making small talk.
You can be direct and to the point. You can also be funny and clever, but you don't have to be.
You can be as creative as you like, but you must be helpful.`
)

type MyAgent struct {
	cfg *copilot.Config
}

func NewAgent(cfg *copilot.Config) *MyAgent {
	return &MyAgent{
		cfg: cfg,
	}
}

func (a *MyAgent) Execute(ctx context.Context, token string, req *copilot.Request, w http.ResponseWriter) error {

	gh := getGitHubClient(token)

	me, _, err := gh.Users.Get(ctx, "")
	if err != nil {
		return err
	}
	fmt.Println("me: ", me.GetName())

	// session := copilot.GetSession(ctx)
	// if session == nil {
	// 	return fmt.Errorf("session not found in context")
	// }

	systemMessage := &copilot.Message{Role: copilot.ChatRoleSystem, Content: PromptStart}

	chatReq := &copilot.CompletionsRequest{
		Messages: []*copilot.Message{systemMessage},
		Model:    copilot.CopilotModelGPT4o,
	}

	chatReq.Messages = append(chatReq.Messages, req.Messages...)

	stream, err := copilot.ChatCompletionsStream(ctx, token, *chatReq, w)
	if err != nil {
		return fmt.Errorf("failed to get chat completions stream: %w", err)
	}

	// Write the response to the stream
	_, err = io.Copy(w, stream)
	if err != nil {
		return fmt.Errorf("failed to write response to stream: %w", err)
	}

	return nil
}

func getGitHubClient(token string) *github.Client {
	client := github.NewClient(nil)
	if envURL := os.Getenv("GITHUB_API_URL"); envURL != "" {
		client.BaseURL, _ = url.Parse(envURL + "/")
	}
	return client.WithAuthToken(token)
}
