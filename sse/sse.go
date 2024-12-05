package sse

import (
	"encoding/json"
	"io"
	"net/http"

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

// WriteStop writes stop SSE data to the writer.
func WriteStop(w io.Writer) error {
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

// WriteStopAndFlush writes stop SSE data to the writer and flushes the writer.
func WriteStopAndFlush(w io.Writer) error {
	if err := WriteStop(w); err != nil {
		return err
	}
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
	return nil
}

// WriteStreamingHeaders writes the headers for a streaming response.
func WriteStreamingHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
}
