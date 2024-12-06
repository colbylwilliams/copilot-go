package agent

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/colbylwilliams/copilot-go"
	"github.com/colbylwilliams/copilot-go/sse"
	"github.com/google/go-github/v67/github"
	"github.com/openai/openai-go"
)

const (
	PromptStart string = `You are a helpful AI assistant. You are here to help the user with their questions.
You are not a human, so you don't have to worry about being polite or making small talk.
You can be direct and to the point. You can also be funny and clever, but you don't have to be.
You can be as creative as you like, but you must be helpful.`
)

type MyAgent struct {
	cfg *copilot.Config
	oai *openai.Client
}

func NewAgent(cfg *copilot.Config, oai *openai.Client) *MyAgent {
	return &MyAgent{
		cfg: cfg,
		oai: oai,
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

	// write the sse headers
	sse.WriteStreamingHeaders(w)

	messages := []openai.ChatCompletionMessageParamUnion{openai.SystemMessage(PromptStart)}

	for _, m := range req.Messages {
		// for now, we won't send the _session message
		// web (dotcom) chat sends us downstream to openai
		if m.IsSessionMessage() {
			continue
		}

		// we'll also skip adding messages with no content
		if m.Content == "" {
			continue
		}

		switch m.Role {
		case copilot.ChatRoleSystem:
			messages = append(messages, openai.SystemMessage(m.Content))

		case copilot.ChatRoleUser:
			// if the message begins with @agent-name then remove it
			messages = append(messages, openai.UserMessage(strings.TrimPrefix(m.Content, fmt.Sprintf("@%s ", req.Agent))))

		case copilot.ChatRoleAssistant:
			messages = append(messages, openai.AssistantMessage(m.Content))

		default:
			return fmt.Errorf("unhandled role: %s", m.Role)
		}
	}

	stream := a.oai.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Messages: openai.F(messages),
		Model:    openai.F(a.cfg.ChatModel),
	})

	for stream.Next() {
		evt := stream.Current()
		chunk, err := agentResponse(&evt)
		if err != nil {
			return err
		}

		sse.WriteDataAndFlush(w, chunk)
	}

	if err := stream.Err(); err != nil {
		return err
	}

	return nil
}

// agentResponse converts a openai ChatCompletionChunk
// to the copilot representation of agent.Response.
func agentResponse(in *openai.ChatCompletionChunk) (*copilot.Response, error) {
	out := copilot.Response{
		ID:                in.ID,
		Created:           in.Created,
		Object:            string(in.Object),
		Model:             in.Model,
		SystemFingerprint: in.SystemFingerprint,
		Choices:           make([]copilot.ChatChoice, len(in.Choices)),
	}

	for i, c := range in.Choices {
		newChoice := copilot.ChatChoice{
			Index:        c.Index,
			FinishReason: string(c.FinishReason),
			Delta: copilot.ChatChoiceDelta{
				Role:    string(c.Delta.Role),
				Content: c.Delta.Content,
			},
		}
		out.Choices[i] = newChoice
	}

	return &out, nil
}

func getGitHubClient(token string) *github.Client {
	client := github.NewClient(nil)
	if envURL := os.Getenv("GITHUB_API_URL"); envURL != "" {
		client.BaseURL, _ = url.Parse(envURL + "/")
	}
	return client.WithAuthToken(token)
}
