package main

import (
	"fmt"
	"net/http"
	"time"
)

var notifs NotificationStore

func HandleScheduler(w http.ResponseWriter, r *http.Request) {
	n := Notification{
		Id:     notifs.ID,
		SendAt: time.Now().Add(time.Second * 10),
	}
	notifs.Mu.Lock()
	notifs.Store = append(notifs.Store, n)
	notifs.Mu.Unlock()
	fmt.Fprintln(w,
		`<head><link rel="icon" href="data:,"></head>`,
		`<p>Notification added, ID:`, notifs.ID, "Quantity:", len(notifs.Store), "</p>",
	)
	notifs.ID += 1
}

func main() {
	notifs = NotificationStore{
		ID:    1,
		Store: []Notification{},
	}

	// Start the background worker
	notifs.Dispatcher()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", HandleScheduler)
	http.ListenAndServe(":3000", mux)
}
