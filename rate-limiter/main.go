package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	limit       int
	window      time.Duration
	requests    int
	windowStart time.Time
	mu          sync.Mutex
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:  limit,
		window: window,
	}
}

func (rl *RateLimiter) Allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if time.Since(rl.windowStart) >= 20*time.Second {
		rl.requests = 0
		rl.windowStart = time.Now()
	}

	if rl.requests >= rl.limit {
		return false
	}

	rl.requests++

	return true
}

var limiter = NewRateLimiter(2, 20*time.Second)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		return
	}

	fmt.Fprintln(w, "Hello from sexy api ")
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	fmt.Printf("server is running in http://localhost:7000")

	http.ListenAndServe(":7000", nil)
}
