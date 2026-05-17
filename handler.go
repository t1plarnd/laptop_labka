package main

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	store Store
}

func NewHandler(s Store) *Handler {
	return &Handler{store: s}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var l Laptop
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	l.ID = 0
	if err := l.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	created, err := h.store.Create(l)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create laptop")
		return
	}
	writeJSON(w, http.StatusCreated, created)
}
