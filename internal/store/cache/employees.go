package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/mubashshir3767/currencyExchange/internal/store"
)

const EmployeeExpTime = time.Minute

type EmployeeStore struct {
	rdb *redis.Client
}

func (s *EmployeeStore) Set(ctx context.Context, employee *store.Employee) error {
	if employee.ID == 0 {
		return fmt.Errorf("employee id not found")
	}
	cacheKey := fmt.Sprintf("employee-%v", employee.ID)

	json, err := json.Marshal(employee)
	if err != nil {
		return err
	}
	log.Println("CACHE SET METHOD USED")

	return s.rdb.Set(ctx, cacheKey, json, UserExpTime).Err()
}

func (s *EmployeeStore) Get(ctx context.Context, employeeId int64) (*store.Employee, error) {
	cacheKey := fmt.Sprintf("employee-%v", employeeId)
	log.Println("CACHE GET METHOD USED")

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var employee store.Employee
	if data != "" {
		err := json.Unmarshal([]byte(data), &employee)
		if err != nil {
			return nil, err
		}
	}

	return &employee, nil
}
