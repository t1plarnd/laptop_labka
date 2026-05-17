package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewMemoryStore()
	h := NewHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /laptops", h.List)
	mux.HandleFunc("POST /laptops", h.Create)
	mux.HandleFunc("GET /laptops/{id}", h.Get)

	addr := ":8080"
	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
