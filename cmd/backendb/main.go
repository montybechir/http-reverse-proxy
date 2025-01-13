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
		fmt.Fprintf(w, "Response from Backend B")
	})

	fmt.Println("Backend B is running on http://localhost:9091")
	if err := http.ListenAndServe(":9091", nil); err != nil {
		log.Fatalf("Backend B failed to start: %v", err)
	}
}
