package client

import (
	"fmt"
	"net/http"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			fmt.Println("Sending request", id)

			resp, err := http.Get("http://localhost:8000")
			if err != nil {
				fmt.Println("Error:", err)
				return
			}

			defer resp.Body.Close()

			fmt.Println("Request", id, "finished with", resp.Status)
		}(i)
	}

	wg.Wait()

	fmt.Println("All requests finished")
}
