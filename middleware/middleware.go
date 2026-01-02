package main

import (
	"fmt"
	"log"
	"net/http"
)

type appHandlerFunc func(w http.ResponseWriter, r *http.Request) error

func errorHandlerMiddleware(next appHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := next(w, r); err != nil {
			log.Println("[ERROR]", err)
		}
	}
}

// make wraps each route to centralise error handling
func make(mux *http.ServeMux, pattern string, handler appHandlerFunc) {
	mux.HandleFunc(pattern, errorHandlerMiddleware(handler))
}

func loggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Auth Middleware")
		next.ServeHTTP(w, r)
	})
}
