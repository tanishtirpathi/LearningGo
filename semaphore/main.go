package main

import (
	"fmt"
	"net/http"
	"time"
)

var semaphore = make(chan struct{}, 4)

func processRequest(w http.ResponseWriter, r *http.Request) {

	semaphore <- struct{}{}

	fmt.Println("Processing request...")
	time.Sleep(2 * time.Second)

	fmt.Fprintln(w, "Processing complete")

	fmt.Println("Request done ")
	<-semaphore
}

func main() {
	http.HandleFunc("/", processRequest)
	fmt.Println("Server started on port 8000")
	http.ListenAndServe(":8000", nil)

}
