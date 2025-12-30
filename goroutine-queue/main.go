package main

import (
	"fmt"
	"net/http"
	"sync"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	set := map[string]string{
		"foo": "69",
		"bar": "420",
		"baz": "1337",
	}
	var wg sync.WaitGroup
	for k, v := range set {
		wg.Go(func() {
			cache.Put(k, v)
		})
	}
	wg.Wait()
	fmt.Fprintf(w, "Cache insertion completed")
}

func main() {
	// Run in the background
	go cache.Worker()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleIndex)
	http.ListenAndServe(":3000", mux)
}
