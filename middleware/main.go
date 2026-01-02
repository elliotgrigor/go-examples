package main

import (
	"errors"
	"fmt"
	"net/http"
)

func handleIndex(w http.ResponseWriter, r *http.Request) error {
	err := errors.New("Simulated error from `handleIndex`")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, `<head><link rel="icon" href="data:,"></head><pre>Internal Server Error</pre>`)
		return err
	}

	// <link> suppresses Chrome's double request
	fmt.Fprintln(w, `<head><link rel="icon" href="data:,"></head><h1>Hello</h1>`)
	return nil
}

func main() {
	var app App
	app.Mux = http.NewServeMux()
	make(app.Mux, "GET /", handleIndex)

	app.Use(loggerMiddleware)
	// app.Use(authMiddleware) // removed to avoid polluting stdout

	http.ListenAndServe(":3000", app.Run())
}
