package main

import (
	"log"
	"net/http"
	"time"
)

const (
	requestLimit = 10
	windowSizeMs = 1 * 1000 // 1s
	// 10 requests/sec
)

type rateLimit struct {
	requests    int
	windowStart time.Time
}

var ipLimits = map[string]*rateLimit{}

// rateLimiter implements an IP-based, fixed-window limiter
func rateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL.Path)

		rl, exists := ipLimits[r.RemoteAddr]
		if !exists {
			rl = &rateLimit{
				requests:    0,
				windowStart: time.Now().UTC(),
			}
			ipLimits[r.RemoteAddr] = rl
		}

		// Reset limiter if time window exceeded
		nowMilli := time.Now().UTC().UnixMilli()
		if nowMilli > rl.windowStart.UnixMilli()+windowSizeMs {
			rl.requests = 0
			rl.windowStart = time.Now().UTC()
		}

		rl.requests += 1
		nowMilli = time.Now().UTC().UnixMilli()
		if nowMilli > rl.windowStart.UnixMilli() && rl.windowStart.UnixMilli() >= nowMilli-windowSizeMs {
			if rl.requests > requestLimit {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
