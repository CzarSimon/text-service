package context

import (
	"context"
	"fmt"
)

type contextKey struct{}

// ContextIDKey id key for context.
var ContextIDKey = contextKey{}

// Context request context.
type Context struct {
	ID       string
	Language string
	context.Context
}

// New creates new context from parent.
func New(parent context.Context, id, language string) *Context {
	return &Context{
		ID:       id,
		Language: language,
		Context:  context.WithValue(parent, ContextIDKey, id),
	}
}

func (c *Context) String() string {
	return fmt.Sprintf("Context(id=[%s], language=%s)", c.ID, c.Language)
}
