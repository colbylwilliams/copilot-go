package copilot

import (
	"context"
)

// GetSessionInfo returns the SessionInfo object from the context.
func GetSessionInfo(ctx context.Context) *SessionInfo {
	val, _ := ctx.Value(sessionCtxKey).(*SessionInfo)
	return val
}

// AddSessionInfo adds the SessionInfo object to the context.
func AddSessionInfo(ctx context.Context, data *SessionInfo) context.Context {
	return context.WithValue(ctx, sessionCtxKey, data)
}

// GetGetHubToken returns the SessionInfo object from the context.
func GetGetHubToken(ctx context.Context) string {
	val, _ := ctx.Value(githubTokenCtxKey).(string)
	return val
}

// AddGetHubToken adds the SessionInfo object to the context.
func AddGetHubToken(ctx context.Context, data string) context.Context {
	return context.WithValue(ctx, githubTokenCtxKey, data)
}

var (
	// sessionCtxKey is the context.Context key to store the session context.
	sessionCtxKey     = &contextKey{"SessionInfo"}
	githubTokenCtxKey = &contextKey{"GithubToken"}
)

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "copilot context value " + k.name
}
