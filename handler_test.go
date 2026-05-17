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

func doRequestWithID(t *testing.T, h http.HandlerFunc, method, target, id, body string) *httptest.ResponseRecorder {
	t.Helper()
	var r *http.Request
	if body == "" {
		r = httptest.NewRequest(method, target, nil)
	} else {
		r = httptest.NewRequest(method, target, strings.NewReader(body))
	}
	r.SetPathValue("id", id)
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

func TestDelete_OK(t *testing.T) {
	called := false
	m := &mockStore{
		deleteFn: func(id int) error {
			called = true
			if id != 4 {
				t.Fatalf("expected id 4, got %d", id)
			}
			return nil
		},
	}
	h := NewHandler(m)
	w := doRequestWithID(t, h.Delete, http.MethodDelete, "/laptops/4", "4", "")
	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
	if !called {
		t.Fatalf("store.Delete was not called")
	}
}

func TestDelete_NotFound(t *testing.T) {
	m := &mockStore{
		deleteFn: func(id int) error {
			return ErrNotFound
		},
	}
	h := NewHandler(m)
	w := doRequestWithID(t, h.Delete, http.MethodDelete, "/laptops/9", "9", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestDelete_BadID(t *testing.T) {
	h := NewHandler(&mockStore{})
	w := doRequestWithID(t, h.Delete, http.MethodDelete, "/laptops/abc", "abc", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDelete_StoreError(t *testing.T) {
	m := &mockStore{
		deleteFn: func(id int) error {
			return errors.New("db down")
		},
	}
	h := NewHandler(m)
	w := doRequestWithID(t, h.Delete, http.MethodDelete, "/laptops/1", "1", "")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestPatch_OK(t *testing.T) {
	m := &mockStore{
		patchFn: func(id int, p LaptopPatch) (Laptop, error) {
			l := Laptop{ID: id, Brand: "Lenovo", RAM: 16}
			if p.RAM != nil {
				l.RAM = *p.RAM
			}
			return l, nil
		},
	}
	h := NewHandler(m)

	w := doRequestWithID(t, h.Patch, http.MethodPatch, "/laptops/1", "1", `{"ram": 32}`)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	var got Laptop
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.RAM != 32 {
		t.Fatalf("expected ram 32, got %d", got.RAM)
	}
}

func TestPatch_NotFound(t *testing.T) {
	m := &mockStore{
		patchFn: func(id int, p LaptopPatch) (Laptop, error) {
			return Laptop{}, ErrNotFound
		},
	}
	h := NewHandler(m)
	w := doRequestWithID(t, h.Patch, http.MethodPatch, "/laptops/9", "9", `{"ram": 32}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestPatch_InvalidJSON(t *testing.T) {
	h := NewHandler(&mockStore{})
	w := doRequestWithID(t, h.Patch, http.MethodPatch, "/laptops/1", "1", "{not json")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestPatch_ValidationError(t *testing.T) {
	h := NewHandler(&mockStore{})
	w := doRequestWithID(t, h.Patch, http.MethodPatch, "/laptops/1", "1", `{"ram": -1}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("ram must be positive")) {
		t.Fatalf("expected validation error, got %s", w.Body.String())
	}
}

func TestPatch_EmptyBodyAllowed(t *testing.T) {
	m := &mockStore{
		patchFn: func(id int, p LaptopPatch) (Laptop, error) {
			return Laptop{ID: id, Brand: "Lenovo"}, nil
		},
	}
	h := NewHandler(m)
	w := doRequestWithID(t, h.Patch, http.MethodPatch, "/laptops/1", "1", `{}`)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestUpdate_OK(t *testing.T) {
	m := &mockStore{
		updateFn: func(id int, l Laptop) (Laptop, error) {
			l.ID = id
			return l, nil
		},
	}
	h := NewHandler(m)

	l := validLaptop()
	l.Brand = "Asus"
	body, _ := json.Marshal(l)
	w := doRequestWithID(t, h.Update, http.MethodPut, "/laptops/3", "3", string(body))

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d, body: %s", w.Code, w.Body.String())
	}
	var got Laptop
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.ID != 3 || got.Brand != "Asus" {
		t.Fatalf("unexpected response: %+v", got)
	}
}

func TestUpdate_NotFound(t *testing.T) {
	m := &mockStore{
		updateFn: func(id int, l Laptop) (Laptop, error) {
			return Laptop{}, ErrNotFound
		},
	}
	h := NewHandler(m)

	body, _ := json.Marshal(validLaptop())
	w := doRequestWithID(t, h.Update, http.MethodPut, "/laptops/9", "9", string(body))
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdate_InvalidJSON(t *testing.T) {
	h := NewHandler(&mockStore{})
	w := doRequestWithID(t, h.Update, http.MethodPut, "/laptops/1", "1", "{not json")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdate_ValidationError(t *testing.T) {
	h := NewHandler(&mockStore{})
	bad := validLaptop()
	bad.RAM = 0
	body, _ := json.Marshal(bad)
	w := doRequestWithID(t, h.Update, http.MethodPut, "/laptops/1", "1", string(body))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !bytes.Contains(w.Body.Bytes(), []byte("ram must be positive")) {
		t.Fatalf("expected validation error in body, got %s", w.Body.String())
	}
}

func TestUpdate_BadID(t *testing.T) {
	h := NewHandler(&mockStore{})
	body, _ := json.Marshal(validLaptop())
	w := doRequestWithID(t, h.Update, http.MethodPut, "/laptops/abc", "abc", string(body))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGet_OK(t *testing.T) {
	m := &mockStore{
		getByIDFn: func(id int) (Laptop, error) {
			return Laptop{ID: id, Brand: "Lenovo"}, nil
		},
	}
	h := NewHandler(m)

	w := doRequestWithID(t, h.Get, http.MethodGet, "/laptops/5", "5", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got Laptop
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if got.ID != 5 || got.Brand != "Lenovo" {
		t.Fatalf("unexpected laptop: %+v", got)
	}
}

func TestGet_NotFound(t *testing.T) {
	m := &mockStore{
		getByIDFn: func(id int) (Laptop, error) {
			return Laptop{}, ErrNotFound
		},
	}
	h := NewHandler(m)
	w := doRequestWithID(t, h.Get, http.MethodGet, "/laptops/99", "99", "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestGet_BadID(t *testing.T) {
	h := NewHandler(&mockStore{})
	cases := []string{"abc", "0", "-1", ""}
	for _, id := range cases {
		w := doRequestWithID(t, h.Get, http.MethodGet, "/laptops/"+id, id, "")
		if w.Code != http.StatusBadRequest {
			t.Fatalf("id %q: expected 400, got %d", id, w.Code)
		}
	}
}

func TestGet_StoreError(t *testing.T) {
	m := &mockStore{
		getByIDFn: func(id int) (Laptop, error) {
			return Laptop{}, errors.New("db down")
		},
	}
	h := NewHandler(m)
	w := doRequestWithID(t, h.Get, http.MethodGet, "/laptops/1", "1", "")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestList_OK(t *testing.T) {
	m := &mockStore{
		getAllFn: func() ([]Laptop, error) {
			return []Laptop{
				{ID: 1, Brand: "Lenovo"},
				{ID: 2, Brand: "Dell"},
			}, nil
		},
	}
	h := NewHandler(m)

	w := doRequest(t, h.List, http.MethodGet, "/laptops", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got []Laptop
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
}

func TestList_Empty(t *testing.T) {
	m := &mockStore{
		getAllFn: func() ([]Laptop, error) {
			return []Laptop{}, nil
		},
	}
	h := NewHandler(m)

	w := doRequest(t, h.List, http.MethodGet, "/laptops", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var got []Laptop
	if err := json.Unmarshal(w.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected empty, got %d", len(got))
	}
}

func TestList_StoreError(t *testing.T) {
	m := &mockStore{
		getAllFn: func() ([]Laptop, error) {
			return nil, errors.New("db down")
		},
	}
	h := NewHandler(m)

	w := doRequest(t, h.List, http.MethodGet, "/laptops", "")
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
