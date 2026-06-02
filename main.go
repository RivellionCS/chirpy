package main

import (
	"fmt"
	"net/http"
)

func main() {
	buildAndStartServer()
}

func buildAndStartServer() error {
	serve := http.NewServeMux()
	indexFilepath := http.Dir(".")
	serve.Handle("/", http.FileServer(indexFilepath))
	server := http.Server{
		Handler: serve,
		Addr: ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		return fmt.Errorf("error starting server: %v", err)
	}
	return nil
}