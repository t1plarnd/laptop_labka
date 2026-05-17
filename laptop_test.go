package main

import (
	"strings"
	"testing"
)

func validLaptop() Laptop {
	return Laptop{
		Brand:        "Lenovo",
		Model:        "ThinkPad X1",
		CPU:          "Intel i7-1260P",
		RAM:          16,
		Storage:      512,
		Price:        1500.50,
		Year:         2023,
		SerialNumber: "SN-ABC123-0001",
	}
}

func TestValidate_OK(t *testing.T) {
	l := validLaptop()
	if err := l.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestValidate_Errors(t *testing.T) {
	cases := []struct {
		name   string
		modify func(l *Laptop)
		want   string
	}{
		{"empty brand", func(l *Laptop) { l.Brand = "" }, "brand is required"},
		{"long brand", func(l *Laptop) { l.Brand = strings.Repeat("x", 51) }, "brand is too long"},
		{"empty model", func(l *Laptop) { l.Model = "" }, "model is required"},
		{"long model", func(l *Laptop) { l.Model = strings.Repeat("x", 51) }, "model is too long"},
		{"empty cpu", func(l *Laptop) { l.CPU = "" }, "cpu is required"},
		{"zero ram", func(l *Laptop) { l.RAM = 0 }, "ram must be positive"},
		{"negative ram", func(l *Laptop) { l.RAM = -4 }, "ram must be positive"},
		{"huge ram", func(l *Laptop) { l.RAM = 1000 }, "ram is too large"},
		{"zero storage", func(l *Laptop) { l.Storage = 0 }, "storage must be positive"},
		{"zero price", func(l *Laptop) { l.Price = 0 }, "price must be positive"},
		{"negative price", func(l *Laptop) { l.Price = -100 }, "price must be positive"},
		{"old year", func(l *Laptop) { l.Year = 1995 }, "year must be between"},
		{"future year", func(l *Laptop) { l.Year = 3000 }, "year must be between"},
		{"bad serial", func(l *Laptop) { l.SerialNumber = "bad" }, "serial_number must match"},
		{"empty serial", func(l *Laptop) { l.SerialNumber = "" }, "serial_number must match"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := validLaptop()
			c.modify(&l)
			err := l.Validate()
			if err == nil {
				t.Fatalf("expected error containing %q, got nil", c.want)
			}
			if !strings.Contains(err.Error(), c.want) {
				t.Fatalf("expected error containing %q, got %q", c.want, err.Error())
			}
		})
	}
}
