package copilot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CopilotModel represents the model to use for completions.
type CopilotModel string

const (
	CopilotModelGPT35      CopilotModel = "gpt-3.5-turbo"
	CopilotModelGPT4       CopilotModel = "gpt-4"
	CopilotModelGPT4o      CopilotModel = "gpt-4o"
	CopilotModelEmbeddings CopilotModel = "text-embedding-ada-002"
)

// CompletionsRequest is a request to the Copilot API to get completions.
type CompletionsRequest struct {
	Model    CopilotModel       `json:"model" default:"gpt-4o"`
	Messages []*Message         `json:"messages"`
	Tools    []*CompletionsTool `json:"tools,omitempty"`
	Stream   bool               `json:"stream"`
}

// CompletionsTool represents a tool to use for completions.
type CompletionsTool struct {
	Type     string                 `json:"type" default:"function"`
	Function ToolFunctionDefinition `json:"function"`
}

const endpoint = "https://api.githubcopilot.com/chat/completions"

// ChatCompletionsStream is a convenience function that sets the Stream field to true
// and calls ChatCompletions.
func ChatCompletionsStream(ctx context.Context, token string, r CompletionsRequest, w io.Writer) (io.ReadCloser, error) {
	r.Stream = true
	return ChatCompletions(ctx, token, r, w)
}

// ChatCompletions sends a request to the the Copilot API to get completions.
func ChatCompletions(ctx context.Context, token string, r CompletionsRequest, w io.Writer) (io.ReadCloser, error) {

	b, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	return res.Body, nil
}
