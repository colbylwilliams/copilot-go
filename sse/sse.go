// Package sse provides utilities for Server-Sent Events (SSE) in the context
// of the GitHub Copilot Extensions APIs.
package sse

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/colbylwilliams/copilot-go"
)

const (
	sseEventNameConfirmation string = "copilot_confirmation"
	sseEventNameReferences   string = "copilot_references"
	sseEventNameErrors       string = "copilot_errors"
)

func flush(w io.Writer) {
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// WriteDone writes a [DONE] SSE message to the writer and flushes the writer.
//
// example output:
//
//	data: [DONE]
func WriteDone(w io.Writer) {
	_, _ = w.Write([]byte("data: [DONE]\n\n"))
	flush(w)
}

// WriteData writes a data SSE message to the writer and flushes the writer.
func WriteData(w io.Writer, v any) error {
	_, _ = w.Write([]byte("data: "))
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	_, _ = w.Write([]byte("\n")) // Encode() adds one newline, so add only one more here.
	flush(w)
	return nil
}

// WriteEvent writes a data SSE event to the writer and flushes the writer.
//
// example output:
//
//	event: event-name
func WriteEvent(w io.Writer, name string) error {
	if _, err := w.Write([]byte("event: " + name + "\n")); err != nil {
		return err
	}
	flush(w)
	return nil
}

// WriteEvent writes a data SSE event and data to the writer and flushes the writer.
//
// example output:
//
//	event: event-name
//	data: {"key": "value"}
func WriteEventData(w io.Writer, name string, data any) error {
	if _, err := w.Write([]byte("event: " + name + "\n")); err != nil {
		return err
	}
	return WriteData(w, data)
}

// WriteErrors writes error SSE events and data to the writer and flushes the writer.
//
// example output:
//
//	event: copilot_errors
//	data: [{
//		"type": "agent",
//		"code": "my.error",
//		"message": "Something failed",
//		"identifier": "my/id"
//	}]
func WriteErrors(w io.Writer, errs []*copilot.Error) error {
	if len(errs) == 0 {
		return nil
	}
	return WriteEventData(w, sseEventNameErrors, errs)
}

// WriteError writes error SSE events and data to the writer and flushes the writer.
//
// example output:
//
//	event: copilot_errors
//	data: [{
//		"type": "agent",
//		"code": "my.error",
//		"message": "Something failed",
//		"identifier": "my/id"
//	}]
func WriteError(w io.Writer, errs *copilot.Error) error {
	return WriteErrors(w, []*copilot.Error{errs})
}

// WriteReferences writes reference SSE events and data to the writer and flushes the writer.
//
// example output:
//
//	event: copilot_references
//	data: [{
//		"type": "custom.reference",
//		"id": "123",
//		"is_implicit": false,
//		"metadata": {
//			"display_name": "My Ref",
//			"display_icon": "",
//			"display_url": ""
//		},
//		"data": {
//			"key": "value"
//		}
//	}]
func WriteReferences(w io.Writer, refs []*copilot.Reference) error {
	if len(refs) == 0 {
		return nil
	}
	return WriteEventData(w, sseEventNameReferences, refs)
}

// WriteReference writes reference SSE events and data to the writer and flushes the writer.
//
// example output:
//
//	event: copilot_references
//	data: [{
//		"type": "custom.reference",
//		"id": "123",
//		"is_implicit": false,
//		"metadata": {
//			"display_name": "My Ref",
//			"display_icon": "",
//			"display_url": ""
//		},
//		"data": {
//			"key": "value"
//		}
//	}]
func WriteReference(w io.Writer, ref *copilot.Reference) error {
	return WriteReferences(w, []*copilot.Reference{ref})
}

// WriteConfirmation writes confirmation SSE events and data to the writer and flushes the writer.
//
// example output:
//
//	event: copilot_confirmation
//	data: {
//		"type": "action",
//		"title": "Proceed?",
//		"message": "You sure?",
//		"confirmation": {
//			"id": "id-123"
//		}
//	}
func WriteConfirmation(w io.Writer, c *copilot.Confirmation) error {
	return WriteEventData(w, sseEventNameConfirmation, c)
}

// WriteDelta writes a custom message to the writer and flushes the writer.
//
// The id must match the id set in previous messages, and match the id used later
// with WriteStop, otherwise some clients will drop the "stickiness" of your agent
// and attribute the messages to copilot.
//
// IMPORTANT: You must call WriteStop after the last message to ensure the chat
// session is properly closed.
//
// example output:
//
//	data: {"id": "123", "created": 1234567890, "choices": [{"delta": {"content": "Hello, world!", "role": "assistant"}}]}
func WriteDelta(w io.Writer, id string, delta string) error {
	return WriteData(w, copilot.Response{
		ID:      id,
		Created: time.Now().UTC().Unix(),
		Choices: []copilot.ChatChoice{{
			Delta: copilot.ChatChoiceDelta{
				Content: delta,
				Role:    string(copilot.ChatRoleAssistant),
			},
		}},
	})
}

// WriteStop writes stop SSE data to the writer and flushes the writer.
//
// The id must match the id set in previous messages, specifically the id used
// with WriteDelta, otherwise some clients will drop the "stickiness" of your
// agent and attribute the messages to copilot.
//
// If you haven't set an id in previous messages, you can use an empty string.
//
// example output:
//
//	data: {"id": "123", "choices": [{"finish_reason": "stop"}]}
//	data: [DONE]
func WriteStop(w io.Writer, id string) error {
	if err := WriteData(w, copilot.Response{
		Choices: []copilot.ChatChoice{{
			FinishReason: copilot.ChatFinishReasonStop,
		}},
	}); err != nil {
		return err
	}
	WriteDone(w)
	return nil
}

// WriteStreamingHeaders writes the headers for a streaming response.
func WriteStreamingHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
}
