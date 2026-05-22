package main

import "errors"

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
