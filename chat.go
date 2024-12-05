package copilot

import "github.com/colbylwilliams/copilot-go/jsonschema"

type ChatRole string

const (
	ChatRoleUser      ChatRole = "user"
	ChatRoleAssistant ChatRole = "assistant"
	ChatRoleFunction  ChatRole = "function"
	ChatRoleTool      ChatRole = "tool"
	ChatRoleSystem    ChatRole = "system"
)

// Request is the request to the agent.
type Request struct {
	ThreadID         string     `json:"copilot_thread_id"`
	Messages         []*Message `json:"messages"`
	Stop             []string   `json:"stop"`
	TopP             float32    `json:"top_p"`
	Temperature      float32    `json:"temperature"`
	MaxTokens        int32      `json:"max_tokens"`
	PresencePenalty  float32    `json:"presence_penalty"`
	FrequencyPenalty float32    `json:"frequency_penalty"`
	Skills           []string   `json:"copilot_skills"`
	Agent            string     `json:"agent"`
}

// Message is a message in the request.
type Message struct {
	Role          ChatRole              `json:"role"`
	Content       string                `json:"content"`
	Name          string                `json:"name,omitempty"`
	References    []*Reference          `json:"copilot_references"`
	Confirmations []*ClientConfirmation `json:"copilot_confirmations"`
	FunctionCall  *ToolFunctionCall     `json:"function_call,omitempty"`
	ToolCalls     []*ToolCall           `json:"tool_calls,omitempty"`
	ToolCallID    string                `json:"tool_call_id,omitempty"`
}

type ToolCall struct {
	ID       string            `json:"id"`
	Type     string            `json:"type" default:"function"`
	Function *ToolFunctionCall `json:"function"`
	// Index    int
}

type ToolFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ToolFunctionDefinition struct {
	Name        string                `json:"name"`
	Description string                `json:"description,omitempty"`
	Parameters  jsonschema.Definition `json:"parameters"`
}

// Response is the response from the agent.
type Response struct {
	ID                string        `json:"id,omitempty"`
	Created           int64         `json:"created,omitempty"`
	Object            string        `json:"object,omitempty" default:"chat.completion.chunk"`
	Model             string        `json:"model,omitempty"`
	SystemFingerprint string        `json:"system_fingerprint,omitempty"`
	Choices           []ChatChoice  `json:"choices"`
	References        []*Reference  `json:"copilot_references,omitempty"`
	Confirmation      *Confirmation `json:"copilot_confirmation,omitempty"`
	Errors            []*Error      `json:"copilot_errors,omitempty"`
}

const (
	ChatFinishReasonStop         string = "stop"
	ChatFinishReasonToolCalls    string = "tool_calls"
	ChatFinishReasonFunctionCall string = "function_call"
)

type ChatChoice struct {
	Index        int64           `json:"index"`
	FinishReason string          `json:"finish_reason,omitempty"`
	Delta        ChatChoiceDelta `json:"delta"`
}

type ChatChoiceDelta struct {
	Content      string                       `json:"content"`
	Role         string                       `json:"role,omitempty"`
	Name         string                       `json:"name,omitempty"`
	FunctionCall *ChatChoiceDeltaFunctionCall `json:"function_call,omitempty"`
}

type ChatChoiceDeltaFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}
