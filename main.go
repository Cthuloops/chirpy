package main

import (
	"net/http"
)

func main() {
	newMux := http.NewServeMux()
	server := http.Server{
		Addr:    ":8080",
		Handler: newMux,
	}
	server.ListenAndServe()
}
