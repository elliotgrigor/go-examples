package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleIndex)

	app := rateLimiter(mux)

	http.ListenAndServe(":3000", app)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<head><link rel="icon" href="data:,"></head><h1>Hello, World!</h1>`)
}
