// Package sse provides utilities for Server-Sent Events (SSE) in the context
// of the GitHub Copilot Extensions APIs.
package sse

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/colbylwilliams/copilot-go"
)

const (
	SseEventNameConfirmation string = "copilot_confirmation"
	SseEventNameReferences   string = "copilot_references"
	SseEventNameErrors       string = "copilot_errors"
)

// WriteDone writes a [DONE] SSE message to the writer.
func WriteDone(w io.Writer) {
	_, _ = w.Write([]byte("data: [DONE]\n\n"))
}

// WriteDone writes a [DONE] SSE message to the writer.
func WriteDoneAndFlush(w io.Writer) {
	_, _ = w.Write([]byte("data: [DONE]\n\n"))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

// WriteData writes a data SSE message to the writer.
func WriteData(w io.Writer, v any) error {
	_, _ = w.Write([]byte("data: "))
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return err
	}
	_, _ = w.Write([]byte("\n")) // Encode() adds one newline, so add only one more here.
	return nil
}

// WriteDataAndFlush writes a data SSE message to the writer and flushes the writer.
func WriteDataAndFlush(w io.Writer, v any) error {
	if err := WriteData(w, v); err != nil {
		return err
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

func WriteRawData(w io.Writer, raw string) error {
	_, _ = w.Write([]byte("data: " + raw + "\n\n"))
	return nil
}

// WriteEvent writes a data SSE event to the writer.
func WriteEvent(w io.Writer, name string) error {
	_, err := w.Write([]byte("event: " + name))
	if err != nil {
		return err
	}
	_, err = w.Write([]byte("\n"))
	if err != nil {
		return err
	}
	return nil
}

// WriteEventAndFlush writes a data SSE message to the writer and flushes the writer.
func WriteEventAndFlush(w io.Writer, name string) error {
	if err := WriteEvent(w, name); err != nil {
		return err
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// WriteErrors writes error SSE events to the writer.
func WriteErrors(w io.Writer, errs []*copilot.Error) error {
	if len(errs) == 0 {
		return nil
	}

	if err := WriteEvent(w, SseEventNameErrors); err != nil {
		return err
	}
	return WriteData(w, errs)
}

// WriteErrorsAndFlush writes error SSE events to the writer and flushes the writer.
func WriteErrorsAndFlush(w io.Writer, errs []*copilot.Error) error {
	if len(errs) == 0 {
		return nil
	}

	if err := WriteEvent(w, SseEventNameErrors); err != nil {
		return err
	}
	return WriteDataAndFlush(w, errs)
}

// WriteError writes error SSE events to the writer.
func WriteError(w io.Writer, errs *copilot.Error) error {
	return WriteErrors(w, []*copilot.Error{errs})
}

// WriteErrorAndFlush writes error SSE events to the writer and flushes the writer.
func WriteErrorAndFlush(w io.Writer, errs *copilot.Error) error {
	return WriteErrorsAndFlush(w, []*copilot.Error{errs})
}

// WriteReferences writes reference SSE events to the writer.
func WriteReferences(w io.Writer, refs []*copilot.Reference) error {
	if len(refs) == 0 {
		return nil
	}

	if err := WriteEvent(w, SseEventNameReferences); err != nil {
		return err
	}
	return WriteData(w, refs)
}

// WriteReferencesAndFlush writes reference SSE events to the writer and flushes the writer.
func WriteReferencesAndFlush(w io.Writer, refs []*copilot.Reference) error {
	if len(refs) == 0 {
		return nil
	}

	if err := WriteEvent(w, SseEventNameReferences); err != nil {
		return err
	}
	return WriteDataAndFlush(w, refs)
}

// WriteReference writes reference SSE events to the writer.
func WriteReference(w io.Writer, ref *copilot.Reference) error {
	return WriteReferences(w, []*copilot.Reference{ref})
}

// WriteReferenceAndFlush writes reference SSE events to the writer and flushes the writer.
func WriteReferenceAndFlush(w io.Writer, ref *copilot.Reference) error {
	return WriteReferencesAndFlush(w, []*copilot.Reference{ref})
}

// WriteConfirmation writes confirmation SSE events to the writer.
func WriteConfirmation(w io.Writer, c *copilot.Confirmation) error {
	if err := WriteEvent(w, SseEventNameConfirmation); err != nil {
		return err
	}

	return WriteData(w, c)
}

// WriteConfirmationAndFlush writes confirmation SSE events to the writer and flushes the writer.
func WriteConfirmationAndFlush(w io.Writer, c *copilot.Confirmation) error {
	if err := WriteEvent(w, SseEventNameConfirmation); err != nil {
		return err
	}

	return WriteDataAndFlush(w, c)
}

// WriteDelta writes a custom message to the writer and flushes the writer.
//
// The id must match the id set in previous messages, and match the id used later
// with WriteStop, otherwise some clients will drop the "stickiness" of your agent
// and attribute the messages to copilot.
//
// IMPORTANT: You must call WriteStop/WriteStopAndFlush after the last message
// to ensure the chat session is properly closed.
func WriteDelta(ctx context.Context, w io.Writer, id string, delta string) error {
	if err := WriteDeltaWithoutFlush(ctx, w, id, delta); err != nil {
		return err
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// WriteDeltaWithoutFlush writes a custom message to the writer.
//
// The id must match the id set in previous messages, and match the id used later
// with WriteStop, otherwise some clients will drop the "stickiness" of your agent
// and attribute the messages to copilot.
//
// IMPORTANT: You must call WriteStop/WriteStopAndFlush after the last message
// to ensure the chat session is properly closed.
func WriteDeltaWithoutFlush(ctx context.Context, w io.Writer, id string, delta string) error {
	resp := copilot.Response{
		ID:      id,
		Created: time.Now().UTC().Unix(),
		Choices: []copilot.ChatChoice{{
			Delta: copilot.ChatChoiceDelta{
				Content: delta,
				Role:    string(copilot.ChatRoleAssistant),
			},
		}},
	}

	if err := WriteData(w, resp); err != nil {
		return err
	}
	return nil
}

// deprecated: use WriteStop instead
func WriteStopAndFlush(w io.Writer, id string) error {
	return WriteStop(w, id)
}

// WriteStop writes stop SSE data to the writer and flushes the writer.
//
// The id must match the id set in previous messages, specifically the id used
// with WriteDeltaWithoutFlush, otherwise some clients will drop the "stickiness"
// of your agent and attribute the messages to copilot.
//
// If you haven't set an id in previous messages, you can use an empty string.
func WriteStop(w io.Writer, id string) error {
	if err := WriteStopWithoutFlush(w, id); err != nil {
		return err
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// WriteStopWithoutFlush writes stop SSE data to the writer.
//
// The id must match the id set in previous messages, specifically the id used
// with WriteDeltaWithoutFlush, otherwise some clients will drop the "stickiness"
// of your agent and attribute the messages to copilot.
//
// If you haven't set an id in previous messages, you can use an empty string.

func WriteStopWithoutFlush(w io.Writer, id string) error {
	stop := copilot.Response{
		Choices: []copilot.ChatChoice{{
			FinishReason: copilot.ChatFinishReasonStop,
		}},
	}
	if err := WriteData(w, stop); err != nil {
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
