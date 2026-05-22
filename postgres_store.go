package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(ctx context.Context, dsn string) (*PostgresStore, error) {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("create pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping: %w", err)
	}
	return &PostgresStore{pool: pool}, nil
}

func (s *PostgresStore) Close() {
	s.pool.Close()
}

func scanLaptop(row pgx.Row) (Laptop, error) {
	var l Laptop
	err := row.Scan(&l.ID, &l.Brand, &l.Model, &l.CPU, &l.RAM, &l.Storage, &l.Price, &l.Year, &l.SerialNumber)
	return l, err
}

func (s *PostgresStore) Create(l Laptop) (Laptop, error) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
		INSERT INTO laptops (brand, model, cpu, ram, storage, price, year, serial_number)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
		RETURNING id, brand, model, cpu, ram, storage, price, year, serial_number
	`, l.Brand, l.Model, l.CPU, l.RAM, l.Storage, l.Price, l.Year, l.SerialNumber)
	return scanLaptop(row)
}

func (s *PostgresStore) GetAll() ([]Laptop, error) {
	ctx := context.Background()
	rows, err := s.pool.Query(ctx, `
		SELECT id, brand, model, cpu, ram, storage, price, year, serial_number
		FROM laptops ORDER BY id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]Laptop, 0)
	for rows.Next() {
		l, err := scanLaptop(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}

func (s *PostgresStore) GetByID(id int) (Laptop, error) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
		SELECT id, brand, model, cpu, ram, storage, price, year, serial_number
		FROM laptops WHERE id = $1
	`, id)
	l, err := scanLaptop(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Laptop{}, ErrNotFound
	}
	return l, err
}

func (s *PostgresStore) Update(id int, l Laptop) (Laptop, error) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
		UPDATE laptops
		SET brand=$1, model=$2, cpu=$3, ram=$4, storage=$5, price=$6, year=$7, serial_number=$8
		WHERE id=$9
		RETURNING id, brand, model, cpu, ram, storage, price, year, serial_number
	`, l.Brand, l.Model, l.CPU, l.RAM, l.Storage, l.Price, l.Year, l.SerialNumber, id)
	updated, err := scanLaptop(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Laptop{}, ErrNotFound
	}
	return updated, err
}

func (s *PostgresStore) Patch(id int, p LaptopPatch) (Laptop, error) {
	ctx := context.Background()
	row := s.pool.QueryRow(ctx, `
		UPDATE laptops SET
			brand         = COALESCE($1, brand),
			model         = COALESCE($2, model),
			cpu           = COALESCE($3, cpu),
			ram           = COALESCE($4, ram),
			storage       = COALESCE($5, storage),
			price         = COALESCE($6, price),
			year          = COALESCE($7, year),
			serial_number = COALESCE($8, serial_number)
		WHERE id = $9
		RETURNING id, brand, model, cpu, ram, storage, price, year, serial_number
	`, p.Brand, p.Model, p.CPU, p.RAM, p.Storage, p.Price, p.Year, p.SerialNumber, id)
	updated, err := scanLaptop(row)
	if errors.Is(err, pgx.ErrNoRows) {
		return Laptop{}, ErrNotFound
	}
	return updated, err
}

func (s *PostgresStore) Delete(id int) error {
	ctx := context.Background()
	tag, err := s.pool.Exec(ctx, `DELETE FROM laptops WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
