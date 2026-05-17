package main

import (
	"errors"
	"sync"
)

var ErrNotFound = errors.New("laptop not found")

type LaptopPatch struct {
	Brand        *string  `json:"brand,omitempty"`
	Model        *string  `json:"model,omitempty"`
	CPU          *string  `json:"cpu,omitempty"`
	RAM          *int     `json:"ram,omitempty"`
	Storage      *int     `json:"storage,omitempty"`
	Price        *float64 `json:"price,omitempty"`
	Year         *int     `json:"year,omitempty"`
	SerialNumber *string  `json:"serial_number,omitempty"`
}

type Store interface {
	Create(l Laptop) (Laptop, error)
	GetAll() ([]Laptop, error)
	GetByID(id int) (Laptop, error)
	Update(id int, l Laptop) (Laptop, error)
	Patch(id int, p LaptopPatch) (Laptop, error)
	Delete(id int) error
}

type MemoryStore struct {
	mu     sync.Mutex
	items  map[int]Laptop
	nextID int
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		items:  make(map[int]Laptop),
		nextID: 1,
	}
}

func (s *MemoryStore) Create(l Laptop) (Laptop, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l.ID = s.nextID
	s.nextID++
	s.items[l.ID] = l
	return l, nil
}

func (s *MemoryStore) GetAll() ([]Laptop, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	out := make([]Laptop, 0, len(s.items))
	for _, l := range s.items {
		out = append(out, l)
	}
	return out, nil
}

func (s *MemoryStore) GetByID(id int) (Laptop, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, ok := s.items[id]
	if !ok {
		return Laptop{}, ErrNotFound
	}
	return l, nil
}

func (s *MemoryStore) Update(id int, l Laptop) (Laptop, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return Laptop{}, ErrNotFound
	}
	l.ID = id
	s.items[id] = l
	return l, nil
}

func (s *MemoryStore) Patch(id int, p LaptopPatch) (Laptop, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	l, ok := s.items[id]
	if !ok {
		return Laptop{}, ErrNotFound
	}
	if p.Brand != nil {
		l.Brand = *p.Brand
	}
	if p.Model != nil {
		l.Model = *p.Model
	}
	if p.CPU != nil {
		l.CPU = *p.CPU
	}
	if p.RAM != nil {
		l.RAM = *p.RAM
	}
	if p.Storage != nil {
		l.Storage = *p.Storage
	}
	if p.Price != nil {
		l.Price = *p.Price
	}
	if p.Year != nil {
		l.Year = *p.Year
	}
	if p.SerialNumber != nil {
		l.SerialNumber = *p.SerialNumber
	}
	s.items[id] = l
	return l, nil
}

func (s *MemoryStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
