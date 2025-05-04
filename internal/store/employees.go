package store

import (
	"context"
	"database/sql"
	"errors"
)

type Employee struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Role      int64  `json:"role"`
	CompanyId int64  `json:"company_id"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type EmployeeStorage struct {
	db *sql.DB
}

func (s *EmployeeStorage) Create(ctx context.Context, employee *Employee) error {
	query := `INSERT INTO employees(username, phone, password, role, company_id)
				VALUES($1, $2, $3, $4) RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		employee.Username,
		employee.Phone,
		employee.Password,
		employee.Role,
		employee.CompanyId).Scan(
		&employee.ID,
		&employee.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *EmployeeStorage) Login(ctx context.Context, employee *Employee) error {
	query := `SELECT id, username, phone, password, role, company_id, created_at FROM employees WHERE phone = $1 AND password = $2`

	err := s.db.QueryRowContext(
		ctx,
		query,
		employee.Phone,
		employee.Password).Scan(
		&employee.ID,
		&employee.Username,
		&employee.Phone,
		&employee.Password,
		&employee.Role,
		&employee.CompanyId,
		&employee.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *EmployeeStorage) GetById(ctx context.Context, id *int64) (*Employee, error) {
	query := `SELECT id, username, phone, password, role, company_id, created_at FROM employees WHERE id = $1`

	employee := &Employee{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&employee.ID,
		&employee.Username,
		&employee.Phone,
		&employee.Password,
		&employee.Role,
		&employee.CompanyId,
		&employee.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (s *EmployeeStorage) GetAll(ctx context.Context) ([]Employee, error) {
	query := `SELECT id, username, phone, password, role, company_id, created_at FROM employees`
	var employees []Employee

	rows, err := s.db.QueryContext(
		ctx,
		query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		employee := &Employee{}
		err := rows.Scan(
			&employee.ID,
			&employee.Username,
			&employee.Phone,
			&employee.Password,
			&employee.Role,
			&employee.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		employees = append(employees, *employee)
	}

	return employees, nil
}

func (s *EmployeeStorage) Update(ctx context.Context, employee *Employee) error {
	query := `UPDATE employees SET username = $1, phone = $2, password = $3, role = $4 WHERE id = $5
				RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		employee.Username,
		employee.Phone,
		employee.Password,
		employee.Role).Scan(
		&employee.ID,
		&employee.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *EmployeeStorage) Delete(ctx context.Context, id *int64) error {
	query := `DELETE FROM employees WHERE id = $1`

	res, err := s.db.ExecContext(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	result, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if result == 0 {
		return errors.New("NOT FOUND")
	}

	return nil
}
