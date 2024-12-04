package copilot

import (
	"context"
)

var (
	// SessionCtxKey is the context.Context key to store the request context.
	SessionCtxKey = &contextKey{"SessionContext"}
)

// Session returns the session Context object from a
// http.Request Context.
func GetSession(ctx context.Context) *SessionContext {
	val, _ := ctx.Value(SessionCtxKey).(*SessionContext)
	return val
}

func AddSession(ctx context.Context, data *SessionContext) context.Context {
	return context.WithValue(ctx, SessionCtxKey, data)
}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "copilot context value " + k.name
}
