package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
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

func parseID(r *http.Request) (int, error) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	l, err := h.store.GetByID(id)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "laptop not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch laptop")
		return
	}
	writeJSON(w, http.StatusOK, l)
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	laptops, err := h.store.GetAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to fetch laptops")
		return
	}
	writeJSON(w, http.StatusOK, laptops)
}

func (h *Handler) Patch(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var p LaptopPatch
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if err := p.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, err := h.store.Patch(id, p)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "laptop not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update laptop")
		return
	}
	writeJSON(w, http.StatusOK, updated)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	var l Laptop
	if err := json.NewDecoder(r.Body).Decode(&l); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	l.ID = id
	if err := l.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	updated, err := h.store.Update(id, l)
	if errors.Is(err, ErrNotFound) {
		writeError(w, http.StatusNotFound, "laptop not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update laptop")
		return
	}
	writeJSON(w, http.StatusOK, updated)
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
