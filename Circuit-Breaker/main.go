package main

// so this is like the M-seal which trip when something fail
// imagine you went in hotel and like call for order and nobody come now what u will do is u wait for sometime and then try again

// this work kinda same

// this also has three state

// open : if some request fails the circuit will open up

// half open : when it is in the state of waiting for some time and then it will try to make request again

// close : when our service is working fine so our circuit is closed

// Lets get into the code now

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type CircuitBreaker struct {
	mu            sync.Mutex
	state         string
	failures      int
	maxFailures   int
	cooldown      time.Duration
	lastFailure   time.Time
	probeInFlight bool
	serviceFails  bool
}

func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{
		state:        "closed",
		maxFailures:  3,
		cooldown:     5 * time.Second,
		serviceFails: true,
	}

}

func (cb *CircuitBreaker) Allow() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	switch cb.state {
	case "open":
		if time.Since(cb.lastFailure) < cb.cooldown {
			return false
		}
		cb.state = "half-open"
		return true
	default:
		return true

	}
}

func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	cb.state = "closed"
	cb.probeInFlight = false
	fmt.Println("Circuit Closed: service is healthy ")
}

func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()
	cb.probeInFlight = false

	if cb.state == "HALF-OPEN" || cb.failures >= cb.maxFailures {
		cb.state = "OPEN"
		fmt.Println("Circuit OPEN: blocking calls")
	}
}

func (cb *CircuitBreaker) ServiceFails() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	return cb.serviceFails
}

func (cb *CircuitBreaker) SetServiceFails(fails bool) {
	cb.mu.Lock()
	defer cb.mu.Unlock()
	cb.serviceFails = fails
}

func (cb *CircuitBreaker) Status() map[string]interface{} {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return map[string]interface{}{
		"state":           cb.state,
		"failures":        cb.failures,
		"service_mode":    map[bool]string{true: "fail", false: "success"}[cb.serviceFails],
		"probe_in_flight": cb.probeInFlight,
	}
}

func main() {
	cb := NewCircuitBreaker()
	http.HandleFunc("/call", func(w http.ResponseWriter, r *http.Request) {
		if !cb.Allow() {
			http.Error(w, "Circuit is OPEN or test request unavailable", http.StatusServiceUnavailable)
			return
		}
		if cb.ServiceFails() {
			cb.RecordFailure()
			http.Error(w, "Downstream service failed", http.StatusBadGateway)
			return
		}

		cb.RecordSuccess()
		fmt.Fprintln(w, "Service call succeeded!")
	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(cb.Status())
	})

	http.HandleFunc("/mode", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Use POST with ?value=fail or ?value=success", http.StatusMethodNotAllowed)
			return
		}

		switch r.URL.Query().Get("value") {
		case "fail":
			cb.SetServiceFails(true)
		case "success":
			cb.SetServiceFails(false)
		default:
			http.Error(w, "value must be fail or success", http.StatusBadRequest)
			return
		}

		fmt.Fprintln(w, "Service mode updated")
	})

	fmt.Println("Circuit Breaker API running on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))

}

// so what exactly happening is the code is like

// it create 3 routes
// status : jaha we can check the status of circuit breaker
// call : jo ki block or allow karega request ko
// mode : jaha hum service ko fail ya success kar sakte hai

// simple

// and like when we first run status will be closed and service mode will be fail means our service is working fine

// then we do call it make our service fail and now the status will be open and like it will now show our service is failing
// then we do mode it will recover our service within 5 second of cool down time
