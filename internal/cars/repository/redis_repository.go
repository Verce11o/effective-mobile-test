package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Verce11o/effective-mobile-test/internal/models"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"time"
)

const (
	productTTL = 3600
)

type CarCacheRepository struct {
	client *redis.Client
	tracer trace.Tracer
}

func NewCarCacheRepository(client *redis.Client, tracer trace.Tracer) *CarCacheRepository {
	return &CarCacheRepository{client: client, tracer: tracer}
}

func (c *CarCacheRepository) GetCarList(ctx context.Context, hash string) (*models.CarList, error) {
	ctx, span := c.tracer.Start(ctx, "carRedis.GetCarList")
	defer span.End()

	carListBytes, err := c.client.Get(ctx, c.createKey(hash)).Bytes()

	if err != nil {
		return nil, err
	}

	var carList models.CarList

	if err = json.Unmarshal(carListBytes, &carList); err != nil {
		return nil, err
	}

	return &carList, nil
}

func (c *CarCacheRepository) SetByIDCtx(ctx context.Context, hash string, cars models.CarList) error {
	ctx, span := c.tracer.Start(ctx, "carRedis.SetByIDCtx")
	defer span.End()

	carListBytes, err := json.Marshal(cars)

	if err != nil {
		return err
	}

	return c.client.Set(ctx, c.createKey(hash), carListBytes, time.Second*time.Duration(productTTL)).Err()
}

func (c *CarCacheRepository) DeleteCarList(ctx context.Context) error {
	ctx, span := c.tracer.Start(ctx, "carRedis.DeleteCarList")
	defer span.End()

	return c.client.FlushDB(ctx).Err()
}

func (c *CarCacheRepository) createKey(hash string) string {
	return fmt.Sprintf("cars:%s", hash)
}
