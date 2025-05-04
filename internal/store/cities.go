package store

import (
	"context"
	"database/sql"
)

type City struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	SubName   *string `json:"sub_name"`
	CreatedAt *string `json:"created_at"`
}

type CityStorage struct {
	db *sql.DB
}

func (s *CityStorage) Create(ctx context.Context, city *City) error {
	query := `INSERT INTO cities(name, sub_name) VALUES($1, $2) RETURNING id, created_at`
	err := s.db.QueryRowContext(ctx, query, city.Name, city.SubName).Scan(&city.ID, &city.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *CityStorage) GetAll(ctx context.Context) ([]City, error) {
	// query := `SELECT * FROM cities`
	return nil, nil
}

func (s *CityStorage) Update(ctx context.Context, city *City) error {

	return nil
}

func (s *CityStorage) Delete(ctx context.Context, id *int64) error {

	return nil
}
