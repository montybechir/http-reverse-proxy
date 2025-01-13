package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Simulate processing time
		time.Sleep(100 * time.Millisecond)
		fmt.Fprintf(w, "Response from Backend A")
	})

	fmt.Println("Backend A is running on http://localhost:9090")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("Backend A failed to start: %v", err)
	}
}
