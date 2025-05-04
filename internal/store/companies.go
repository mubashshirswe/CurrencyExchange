package store

import (
	"context"
	"database/sql"
	"errors"
)

type Company struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Details   string `json:"details"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type CompanyStorage struct {
	db *sql.DB
}

func (s *CompanyStorage) Create(ctx context.Context, company *Company) error {
	query := `INSERT INTO companies(name, details, password) VALUES($1, $2, $3) RETURNING id, created_at`
	err := s.db.QueryRowContext(ctx, query, company.Name, company.Details, company.Password).Scan(
		&company.ID,
		&company.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *CompanyStorage) Update(ctx context.Context, company *Company) error {
	query := `UPDATE companies SET name = $1, details = $2, password = $3 WHERE id $4`

	rows, err := s.db.ExecContext(ctx, query, company.Name, company.Details, company.Password, company.ID)

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

func (s *CompanyStorage) GetAll(ctx context.Context) ([]Company, error) {
	query := `SELECT * FROM companies`

	var companies []Company
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		company := &Company{}

		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.Details,
			&company.Password,
			&company.CreatedAt,
		)

		if err != nil {
			return nil, err
		}
		companies = append(companies, *company)
	}

	return companies, nil
}

func (s *CompanyStorage) GetById(ctx context.Context, id *int64) (*Company, error) {
	query := `SELECT * FROM companies WHERE id = $1`

	company := &Company{}

	err := s.db.QueryRowContext(ctx, query).Scan(
		&company.ID,
		&company.Name,
		&company.Details,
		&company.Password,
		&company.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return company, nil
}

func (s *CompanyStorage) Delete(ctx context.Context, id *int64) error {
	query := `DELETE FROM companies WHERE id = $1`

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
