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

func (c *CarRepository) CreateCars(ctx context.Context, input []domain.Car) error {

	tx, err := c.db.Begin(ctx)

	if err != nil {
		return fmt.Errorf("error begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	for _, car := range input {
		var ownerID int
		row := tx.QueryRow(ctx, "INSERT INTO owners (name, surname, patronymic) VALUES ($1, $2, $3) RETURNING id", car.Owner.Name, car.Owner.Surname, car.Owner.Patronymic)

		if err := row.Scan(&ownerID); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, "INSERT INTO cars (reg_num, mark, model, ownerid) VALUES ($1, $2, $3, $4)", car.RegNum, car.Mark, car.Model, ownerID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
