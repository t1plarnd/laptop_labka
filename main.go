package main

import (
	"context"
	"log"
	"net/http"
	"os"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	store, err := NewPostgresStore(context.Background(), dsn)
	if err != nil {
		log.Fatalf("postgres: %v", err)
	}
	defer store.Close()

	h := NewHandler(store)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("GET /laptops", h.List)
	mux.HandleFunc("POST /laptops", h.Create)
	mux.HandleFunc("GET /laptops/{id}", h.Get)
	mux.HandleFunc("PUT /laptops/{id}", h.Update)
	mux.HandleFunc("PATCH /laptops/{id}", h.Patch)
	mux.HandleFunc("DELETE /laptops/{id}", h.Delete)

	addr := ":8080"
	log.Printf("server starting on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
