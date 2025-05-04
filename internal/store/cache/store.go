package cache

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/mubashshir3767/currencyExchange/internal/store"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*store.User, error)
		Set(context.Context, *store.User) error
	}

	Employees interface {
		Get(context.Context, int64) (*store.Employee, error)
		Set(context.Context, *store.Employee) error
	}
}

func NewRedisStorage(redis *redis.Client) Storage {
	return Storage{
		Users:     &UsersStore{rdb: redis},
		Employees: &EmployeeStore{rdb: redis},
	}
}
