package main

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type Laptop struct {
	ID           int     `json:"id"`
	Brand        string  `json:"brand"`
	Model        string  `json:"model"`
	CPU          string  `json:"cpu"`
	RAM          int     `json:"ram"`
	Storage      int     `json:"storage"`
	Price        float64 `json:"price"`
	Year         int     `json:"year"`
	SerialNumber string  `json:"serial_number"`
}

var serialRegex = regexp.MustCompile(`^SN-[A-Z0-9]{6}-[0-9]{4}$`)

func (l *Laptop) Validate() error {
	if l.Brand == "" {
		return errors.New("brand is required")
	}
	if len(l.Brand) > 50 {
		return errors.New("brand is too long")
	}
	if l.Model == "" {
		return errors.New("model is required")
	}
	if len(l.Model) > 50 {
		return errors.New("model is too long")
	}
	if l.CPU == "" {
		return errors.New("cpu is required")
	}
	if l.RAM <= 0 {
		return errors.New("ram must be positive")
	}
	if l.RAM > 256 {
		return errors.New("ram is too large")
	}
	if l.Storage <= 0 {
		return errors.New("storage must be positive")
	}
	if l.Price <= 0 {
		return errors.New("price must be positive")
	}
	maxYear := time.Now().Year() + 1
	if l.Year < 2000 || l.Year > maxYear {
		return fmt.Errorf("year must be between 2000 and %d", maxYear)
	}
	if !serialRegex.MatchString(l.SerialNumber) {
		return errors.New("serial_number must match format SN-XXXXXX-XXXX")
	}
	return nil
}
