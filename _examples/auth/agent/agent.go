package agent

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/colbylwilliams/copilot-go"
	"github.com/colbylwilliams/copilot-go/sse"
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
	// change this to "confirmation" or "reference" to test those events
	event := "error"

	// write the sse headers
	sse.WriteStreamingHeaders(w)

	switch event {
	case "error":

		copilotError := &copilot.Error{
			Type:       copilot.ErrorTypeAgent,
			Code:       "test.error",
			Message:    "this is a test error",
			Identifier: req.Agent,
		}

		// write the error to the stream
		sse.WriteErrorAndFlush(w, copilotError)

	case "confirmation":

		// if the last message in req.Messages has items in the confirmations array
		// then the user just responded to a confirmation
		if last := req.Messages[len(req.Messages)-1]; last != nil && last.Role == copilot.ChatRoleUser && last.Confirmations != nil && len(last.Confirmations) > 0 {

			// confirmation := last.Confirmations[0]
			// the user has confirmed the action
			// TODO: print the outcome of the action

			sse.WriteStopAndFlush(w, "")
			return nil
		} else {
			confirmation := &copilot.Confirmation{
				Type:    copilot.ConfirmationTypeAction,
				Title:   "Hello there",
				Message: "Are you sure you want to continue?",
				Confirmation: map[string]interface{}{
					"id":  "123",
					"key": "value",
				},
			}

			// write the confirmation to the stream
			sse.WriteConfirmationAndFlush(w, confirmation)
		}

	case "reference":

		reference := &copilot.Reference{
			Type: "my.custom.type",
			ID:   "123",
			Metadata: copilot.ReferenceMetadata{
				DisplayName: "My Custom Reference",
				DisplayIcon: "https://example.com/icon.png",
				DisplayURL:  "https://example.com",
			},
			RawData: json.RawMessage(`{"key": "value", "key2": "value2"}`),
		}

		// write the reference to the stream
		sse.WriteReferenceAndFlush(w, reference)
	}

	sse.WriteStopAndFlush(w, "")

	return nil
}
