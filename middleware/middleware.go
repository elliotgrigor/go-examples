package main

import (
	"fmt"
	"net/http"
)

// Implement a custom http.HandlerFunc which returns errors to be caught by
// the error handling middleware
type appHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func errorHandlerMiddleware(next appHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			logger.Error(err.Error())
		}
	}
}

// make wraps each route to centralise error handling
func make(mux *http.ServeMux, pattern string, handler appHandlerFunc) {
	mux.HandleFunc(pattern, errorHandlerMiddleware(handler))
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info(fmt.Sprintf("%s %s %s", r.RemoteAddr, r.Method, r.URL.Path))
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("Auth Middleware")
		next.ServeHTTP(w, r)
	})
}
