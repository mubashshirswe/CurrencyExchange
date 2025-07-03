package store

import (
	"context"
	"errors"
)

type User struct {
	ID        int64   `json:"id"`
	Username  string  `json:"username"`
	Phone     string  `json:"phone"`
	Role      int64   `json:"role"`
	Password  string  `json:"password"`
	Avatar    *string `json:"avatar"`
	CompanyId int64   `json:"company_id"`
	CreatedAt string  `json:"created_at"`
}

type UserStorage struct {
	db DBTX
}

func NewUserStorage(db DBTX) *UserStorage {
	return &UserStorage{db: db}
}

func (s *UserStorage) Create(ctx context.Context, user *User) error {
	query := `INSERT INTO users(username, phone, password, role, company_id)
				VALUES($1, $2, $3, $4, $5) RETURNING id, created_at`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Phone,
		user.Password,
		user.Role,
		user.CompanyId).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) Login(ctx context.Context, user *User) error {
	query := `SELECT * FROM users WHERE phone = $1 AND password = $2`

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Phone,
		user.Password).Scan(
		&user.ID,
		&user.Phone,
		&user.Role,
		&user.Avatar,
		&user.Username,
		&user.Password,
		&user.CompanyId,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *UserStorage) GetById(ctx context.Context, id *int64) (*User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	user := &User{}

	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Phone,
		&user.Role,
		&user.Avatar,
		&user.Username,
		&user.Password,
		&user.CompanyId,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserStorage) GetAll(ctx context.Context) ([]User, error) {
	query := `SELECT * FROM users`
	var users []User

	rows, err := s.db.QueryContext(
		ctx,
		query)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user := &User{}
		rows.Scan(
			&user.ID,
			&user.Phone,
			&user.Role,
			&user.Avatar,
			&user.Username,
			&user.Password,
			&user.CompanyId,
			&user.CreatedAt,
		)

		users = append(users, *user)
	}

	return users, nil
}

func (s *UserStorage) Update(ctx context.Context, user *User) error {
	query := `UPDATE users SET username = $1, password = $2, role = $3, avatar = $4, company_id = $5 WHERE id = $6`

	result, err := s.db.ExecContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Role,
		user.Avatar,
		user.CompanyId,
		user.ID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (s *UserStorage) Delete(ctx context.Context, id *int64) error {
	query := `DELETE FROM users WHERE id = $1`

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
