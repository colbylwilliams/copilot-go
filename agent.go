package copilot

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	GitHubTokenHeader         = "X-Github-Token"
	PublicKeyIdentifierHeader = "Github-Public-Key-Identifier"
	PublicKeySignatureHeader  = "Github-Public-Key-Signature"
)

type Agent interface {
	Execute(ctx context.Context, token string, req *Request, w http.ResponseWriter) error
}

func AgentHandler(v PayloadVerifier, a Agent) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = getRequiredHeader(r, PublicKeyIdentifierHeader)
		signature := getRequiredHeader(r, PublicKeySignatureHeader)

		b, err := io.ReadAll(r.Body)
		if err != nil {
			fmt.Println(fmt.Errorf("failed to read request body: %w", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		isValid, err := v.Verify(b, signature)
		if err != nil {
			fmt.Printf("failed to validate payload signature: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if !isValid {
			http.Error(w, "invalid payload signature", http.StatusUnauthorized)
			return
		}

		ctx := r.Context()

		token := getRequiredHeader(r, GitHubTokenHeader)
		ctx = AddGetHubToken(ctx, token)

		var req Request
		if err := json.Unmarshal(b, &req); err != nil {
			fmt.Printf("failed to unmarshal request: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		session, err := req.GetSessionInfo()
		if err != nil {
			fmt.Println("error getting session context: ", err)
		}

		ctx = AddSessionInfo(ctx, session)

		if err := a.Execute(ctx, token, &req, w); err != nil {
			fmt.Printf("failed to execute agent: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func getRequiredHeader(r *http.Request, key string) string {
	value := r.Header.Get(key)
	if value == "" {
		panic(fmt.Errorf("missing required header %s", key))
	}
	return value
}
