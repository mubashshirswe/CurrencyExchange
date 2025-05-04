package store

import (
	"context"
	"database/sql"
	"errors"
)

type City struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	SubName   *string `json:"sub_name"`
	CompanyId int64   `json:"company_id"`
	CreatedAt *string `json:"created_at"`
}

type CityStorage struct {
	db *sql.DB
}

func (s *CityStorage) Create(ctx context.Context, city *City) error {
	query := `INSERT INTO cities(name, sub_name, company_id) VALUES($1, $2, $3) RETURNING id, created_at`
	err := s.db.QueryRowContext(ctx, query, city.Name, city.SubName, city.CompanyId).Scan(&city.ID, &city.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *CityStorage) GetAll(ctx context.Context) ([]City, error) {
	query := `SELECT id, name, sub_name, company_id, created_at FROM cities`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cities []City
	for rows.Next() {
		city := &City{}
		rows.Scan(
			city.ID,
			city.Name,
			city.SubName,
			city.CompanyId,
			city.CreatedAt,
		)

		cities = append(cities, *city)
	}

	return cities, nil
}

func (s *CityStorage) Update(ctx context.Context, city *City) error {
	query := `UPDATE cities SET name = $1, sub_name $2 WHERE id = $3`
	rows, err := s.db.ExecContext(ctx, query, city.Name, city.SubName, city.ID)
	if err != nil {
		return err
	}

	res, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if res == 0 {
		return errors.New("NOT FOUND")
	}

	return nil
}

func (s *CityStorage) Delete(ctx context.Context, id *int64) error {
	query := `DELETE FROM cities  WHERE id = $1`
	rows, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	res, err := rows.RowsAffected()
	if err != nil {
		return err
	}

	if res == 0 {
		return errors.New("NOT FOUND")
	}

	return nil
}
