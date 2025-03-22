package main

import (
	"log"
	"net/http"
)

func main() {
	const filepathRoot = "."
	const port = "8080"

	serveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(filepathRoot))
	serveMux.Handle("/", fileServer)

	// create server
	server := &http.Server{
		Addr:    ":" + port,
		Handler: serveMux,
	}
	// start the server
	log.Printf("Server running on port :%s\n", port)
	log.Fatal(server.ListenAndServe())
}
