package main

import (
	"fmt"
	"net/http"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// <link> suppresses Chrome's double request
	fmt.Fprintln(w, `<head><link rel="icon" href="data:,"></head><h1>Hello</h1>`)
}

func main() {
	var app App
	app.Mux = http.NewServeMux()
	app.Mux.HandleFunc("GET /", handleIndex)

	app.Use(logMiddleware)
	app.Use(authMiddleware)
	app.Use(errorHandlerMiddleware)

	http.ListenAndServe(":3000", app.Run())
}
