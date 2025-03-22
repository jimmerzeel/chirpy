package main

import (
	"log"
	"net/http"
)

func main() {
	// create request handler
	serveMux := http.NewServeMux()

	// create server
	server := http.Server{
		Addr:    ":8080",
		Handler: serveMux,
	}

	// start the server
	err := server.ListenAndServe()
	if err != nil {
		log.Printf("error starting the HTTP server: %v", err)
	}
}
