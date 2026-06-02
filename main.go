package main

import (
	"fmt"
	"net/http"
)

func main() {
	buildAndStartServer()
}

func buildAndStartServer() error {
	mux := http.NewServeMux()
	indexFilepath := http.Dir(".")
	mux.Handle("/", http.FileServer(indexFilepath))
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}
	return nil
}