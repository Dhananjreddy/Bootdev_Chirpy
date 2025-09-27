package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(http.Dir("."))))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("content-type", "text/plain; charset=utf-8")
		w.Write([]byte("OK"))
	})

	Server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	Server.ListenAndServe()
}