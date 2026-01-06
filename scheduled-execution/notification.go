package main

import (
	"fmt"
	"sync"
	"time"
)

type Notification struct {
	Id     int
	SendAt time.Time
}

type NotificationStore struct {
	ID    int // Saves next notification ID
	Store []Notification
	Mu    sync.Mutex
}

func (ns *NotificationStore) Dispatcher() {
	ticker := time.NewTicker(time.Second * 1)
	go func() {
		for range ticker.C {
			if len(notifs.Store) < 1 {
				continue
			}
			notifs.Mu.Lock()
			now := time.Now()
			var remaining []Notification
			for _, n := range notifs.Store {
				if now.Equal(n.SendAt) || now.After(n.SendAt) {
					fmt.Println("Sent notification at", n.SendAt.Format(time.RFC3339))
				} else {
					// Remove sent notification from slice by creating a new
					// slice of still-to-be-sent notifications
					// TODO: This could be improved to minimise allocations
					remaining = append(remaining, n)
				}
			}
			notifs.Store = remaining
			notifs.Mu.Unlock()
		}
	}()
}
