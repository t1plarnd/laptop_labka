package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func doRequest(t *testing.T, h http.HandlerFunc, method, target string, body string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, target, nil)
	} else {
		r = httptest.NewRequest(method, target, strings.NewReader(body))
	}
	w := httptest.NewRecorder()
	h(w, r)
	return w
}

func TestCreate_OK(t *testing.T) {
	m := &mockStore{
		createFn: func(l Laptop) (Laptop, error) {
			l.ID = 42
			return l, nil
		},
	}
	h := NewHandler(m)

	body, _ := json.Marshal(validLaptop())
	w := doRequest(t, h.Create, http.MethodPost, "/laptops", string(body))

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d, body: %s", w.Code, w.Body.String())
	}
	var got Laptop
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid response json: %v", err)
	}
	if got.ID != 42 {
		t.Fatalf("expected id 42, got %d", got.ID)
	}
}

func TestCreate_InvalidJSON(t *testing.T) {
	h := NewHandler(&mockStore{})
	w := doRequest(t, h.Create, http.MethodPost, "/laptops", "{not json")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreate_ValidationError(t *testing.T) {
	h := NewHandler(&mockStore{})
	bad := validLaptop()
	bad.Brand = ""
	body, _ := json.Marshal(bad)
	w := doRequest(t, h.Create, http.MethodPost, "/laptops", string(body))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("brand is required")) {
		t.Fatalf("expected validation error in body, got %s", w.Body.String())
	}
}

func TestCreate_StoreError(t *testing.T) {
	m := &mockStore{
		createFn: func(l Laptop) (Laptop, error) {
			return Laptop{}, errors.New("db down")
		},
	}
	h := NewHandler(m)
	body, _ := json.Marshal(validLaptop())
	w := doRequest(t, h.Create, http.MethodPost, "/laptops", string(body))
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestCreate_IgnoresClientID(t *testing.T) {
	var seen Laptop
	m := &mockStore{
		createFn: func(l Laptop) (Laptop, error) {
			seen = l
			l.ID = 7
			return l, nil
		},
	}
	h := NewHandler(m)
	l := validLaptop()
	l.ID = 999
	body, _ := json.Marshal(l)
	w := doRequest(t, h.Create, http.MethodPost, "/laptops", string(body))
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	if seen.ID != 0 {
		t.Fatalf("client-provided id must be ignored, got %d", seen.ID)
	}
}
