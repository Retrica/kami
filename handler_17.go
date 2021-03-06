// +build go1.7,!go1.9

package kami

import (
	"context"
	"fmt"
	"net/http"

	netcontext "golang.org/x/net/context"
)

// OldContextHandler is like ContextHandler but uses the old x/net/context.
type OldContextHandler interface {
	ServeHTTPContext(netcontext.Context, http.ResponseWriter, *http.Request)
}

// wrap tries to turn a HandlerType into a ContextHandler
func wrap(h HandlerType) ContextHandler {
	switch x := h.(type) {
	case ContextHandler:
		return x
	case func(context.Context, http.ResponseWriter, *http.Request):
		return HandlerFunc(x)
	case func(netcontext.Context, http.ResponseWriter, *http.Request):
		return HandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			x(ctx, w, r)
		})
	case http.Handler:
		return HandlerFunc(func(_ context.Context, w http.ResponseWriter, r *http.Request) {
			x.ServeHTTP(w, r)
		})
	case func(http.ResponseWriter, *http.Request):
		return HandlerFunc(func(_ context.Context, w http.ResponseWriter, r *http.Request) {
			x(w, r)
		})
	}
	panic(fmt.Errorf("unsupported HandlerType: %T", h))
}
