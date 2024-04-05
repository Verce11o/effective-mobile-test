package repository

import (
	"context"
	"fmt"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CarRepository struct {
	db *pgxpool.Pool
}

func NewCarRepository(db *pgxpool.Pool) *CarRepository {
	return &CarRepository{db: db}
}

func (c *CarRepository) CreateCar(ctx context.Context, input domain.CreateCarsRequest) error {

	tx, err := c.db.Begin(ctx)

	if err != nil {
		return fmt.Errorf("error begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	for _, regNum := range input.RegNums {
		_, err = tx.Exec(ctx, "INSERT INTO cars (reg_num) VALUES ($1)", regNum)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
