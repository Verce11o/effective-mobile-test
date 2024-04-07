package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

const (
	productTTL = 3600
)

type CarCacheRepository struct {
	client *redis.Client
}

func NewCarCacheRepository(client *redis.Client) *CarCacheRepository {
	return &CarCacheRepository{client: client}
}

func (c *CarCacheRepository) GetCarList(ctx context.Context, cursor string) (*domain.CarList, error) {
	carListBytes, err := c.client.Get(ctx, c.createKey(cursor)).Bytes()

	if err != nil {
		return nil, err
	}

	var carList domain.CarList

	if err = json.Unmarshal(carListBytes, &carList); err != nil {
		return nil, err
	}

	return &carList, nil
}

func (c *CarCacheRepository) SetByIDCtx(ctx context.Context, cursor string, cars domain.CarList) error {
	carListBytes, err := json.Marshal(cars)

	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.createKey(cursor), carListBytes, time.Second*time.Duration(productTTL)).Err()
}

func (c *CarCacheRepository) DeleteCarList(ctx context.Context) error {
	return c.client.FlushDB(ctx).Err()
}

func (c *CarCacheRepository) createKey(cursor string) string {
	return fmt.Sprintf("cars:%s", cursor)
}
