package main

import (
	"errors"
	"testing"
)

func TestMemoryStore_CreateAssignsID(t *testing.T) {
	s := NewMemoryStore()
	l1, _ := s.Create(validLaptop())
	l2, _ := s.Create(validLaptop())
	if l1.ID != 1 || l2.ID != 2 {
		t.Fatalf("expected ids 1 and 2, got %d and %d", l1.ID, l2.ID)
	}
}

func TestMemoryStore_GetByID(t *testing.T) {
	s := NewMemoryStore()
	created, _ := s.Create(validLaptop())

	got, err := s.GetByID(created.ID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Brand != created.Brand {
		t.Fatalf("expected brand %q, got %q", created.Brand, got.Brand)
	}

	_, err = s.GetByID(999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStore_GetAll(t *testing.T) {
	s := NewMemoryStore()
	all, _ := s.GetAll()
	if len(all) != 0 {
		t.Fatalf("expected empty, got %d", len(all))
	}
	s.Create(validLaptop())
	s.Create(validLaptop())
	all, _ = s.GetAll()
	if len(all) != 2 {
		t.Fatalf("expected 2 items, got %d", len(all))
	}
}

func TestMemoryStore_Update(t *testing.T) {
	s := NewMemoryStore()
	created, _ := s.Create(validLaptop())

	upd := validLaptop()
	upd.Brand = "Dell"
	got, err := s.Update(created.ID, upd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Brand != "Dell" {
		t.Fatalf("expected Dell, got %q", got.Brand)
	}
	if got.ID != created.ID {
		t.Fatalf("id changed: expected %d, got %d", created.ID, got.ID)
	}

	_, err = s.Update(999, upd)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStore_Patch(t *testing.T) {
	s := NewMemoryStore()
	created, _ := s.Create(validLaptop())

	newBrand := "HP"
	newRAM := 32
	got, err := s.Patch(created.ID, LaptopPatch{Brand: &newBrand, RAM: &newRAM})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Brand != "HP" || got.RAM != 32 {
		t.Fatalf("patch not applied: %+v", got)
	}
	if got.Model != created.Model {
		t.Fatalf("model should be unchanged, got %q", got.Model)
	}

	_, err = s.Patch(999, LaptopPatch{Brand: &newBrand})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	s := NewMemoryStore()
	created, _ := s.Create(validLaptop())

	if err := s.Delete(created.ID); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := s.GetByID(created.ID); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}

	if err := s.Delete(999); !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
