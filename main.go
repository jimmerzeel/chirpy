package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	serveMux := http.NewServeMux()

	// create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	// start the server
	log.Printf("Server running on port :%s\n", port)
	log.Fatal(server.ListenAndServe())
}
