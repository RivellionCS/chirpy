package main

import (
	"log"
	"net/http"
)

func main() {
	buildAndStartServer()
}

func buildAndStartServer() {
	const filepathRoot = "."
	const port = "8080"
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(filepathRoot)))
	server := http.Server{
		Handler: mux,
		Addr: ":" + port,
	}
	log.Printf("Serving files from %s on port: %s", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}