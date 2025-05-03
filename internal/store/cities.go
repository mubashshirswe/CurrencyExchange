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

	return nil
}

func (s *CityStorage) GetAll(ctx context.Context) ([]City, error) {

	return nil, nil
}
