package main

import (
	"net/http"
	"slices"
)

type MiddlewareFunc func(http.Handler) http.Handler

type App struct {
	Mux             *http.ServeMux
	MiddlewareChain []MiddlewareFunc
}

// Use inserts the middleware in reverse order to a slice
// Reverse order results in first inserted == first executed
func (a *App) Use(mw MiddlewareFunc) {
	a.MiddlewareChain = slices.Insert(a.MiddlewareChain, 0, mw)
}

// Run wraps the root ServeMux in the middleware chain
func (a *App) Run() http.Handler {
	if len(a.MiddlewareChain) < 1 {
		return http.HandlerFunc(a.Mux.ServeHTTP)
	}
	var wrapper http.Handler
	for i, mw := range a.MiddlewareChain {
		if i == 0 {
			wrapper = mw(a.Mux)
			continue
		}
		wrapper = mw(wrapper)
	}
	return wrapper
}
