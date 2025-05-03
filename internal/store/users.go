package store

import (
	"context"
	"database/sql"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	Phone     string `json:"phone"`
	Password  string `json:"password"`
	CreatedAt string `json:"created_at"`
}

type UserStorage struct {
	db *sql.DB
}

func (s *UserStorage) Create(ctx context.Context, user *User) error {

	return nil
}

func (s *UserStorage) Login(ctx context.Context, user *User) error {

	return nil
}

func (s *UserStorage) GetAll(ctx context.Context) ([]User, error) {

	return nil, nil
}

func (s *UserStorage) Update(ctx context.Context, user *User) error {

	return nil
}
