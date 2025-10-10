package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from example-api!")
	})

	port := ":8080"
	fmt.Printf("🐝 Server running on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
