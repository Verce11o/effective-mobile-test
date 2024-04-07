package repository

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"github.com/Verce11o/effective-mobile-test/internal/lib/pagination"
	"github.com/Verce11o/effective-mobile-test/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/trace"
)

const (
	paginationLimit = 10
)

type CarRepository struct {
	db     *pgxpool.Pool
	tracer trace.Tracer
}

func NewCarRepository(db *pgxpool.Pool, tracer trace.Tracer) *CarRepository {
	return &CarRepository{db: db, tracer: tracer}
}

func (c *CarRepository) CreateCars(ctx context.Context, input []domain.Car) error {
	ctx, span := c.tracer.Start(ctx, "carRepository.CreateCars")
	defer span.End()

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
		_, err = tx.Exec(ctx, "INSERT INTO cars (reg_num, mark, model, year, ownerid) VALUES ($1, $2, $3, $4, $5)", car.RegNum, car.Mark, car.Model, car.Year, ownerID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (c *CarRepository) GetCars(ctx context.Context, input domain.GetCarsRequest) (models.CarList, error) {
	ctx, span := c.tracer.Start(ctx, "carRepository.GetCars")
	defer span.End()

	var id int
	var err error

	if input.Cursor != "" {
		id, err = pagination.DecodeCursor(input.Cursor)

		if err != nil {
			return models.CarList{}, err
		}
	}

	query := sq.Select("cars.id, cars.reg_num, cars.mark, cars.model, cars.year, cars.created_at, o.id, o.name, o.surname, o.patronymic").
		From("cars").InnerJoin("owners o on o.id = cars.ownerid").
		OrderBy("cars.created_at, cars.id").Limit(paginationLimit)

	if id != 0 {
		query = query.Where(sq.Gt{"cars.id": id})
	}

	if input.RegNum != "" {
		query = query.Where(sq.Eq{"cars.reg_num": input.RegNum})
	}

	if input.Mark != "" {
		query = query.Where(sq.Eq{"cars.mark": input.Mark})
	}

	if input.Model != "" {
		query = query.Where(sq.Eq{"cars.model": input.Model})
	}

	if input.Year != 0 {
		query = query.Where(sq.Eq{"cars.year": input.Year})
	}

	if input.OwnerID != 0 {
		query = query.Where(sq.Eq{"cars.ownerid": input.OwnerID})
	}

	sql, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return models.CarList{}, err
	}

	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		return models.CarList{}, err
	}

	defer rows.Close()

	cars := make([]models.Car, 0)

	for rows.Next() {
		var car models.Car

		err = rows.Scan(&car.ID, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.CreatedAt, &car.Owner.ID, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic)
		if err != nil {
			return models.CarList{}, err
		}

		cars = append(cars, car)

	}

	var nextCursor string

	if len(cars) > 0 {
		last := cars[len(cars)-1]
		nextCursor = pagination.EncodeCursor(last.ID)

	}
	return models.CarList{
		Cursor: nextCursor,
		Cars:   cars,
		Total:  len(cars),
	}, nil
}

func (c *CarRepository) UpdateCar(ctx context.Context, carID int, input domain.UpdateCarsRequest) error {
	ctx, span := c.tracer.Start(ctx, "carRepository.UpdateCar")
	defer span.End()

	tx, err := c.db.Begin(ctx)

	if err != nil {
		return fmt.Errorf("error begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

	var exists bool
	findCar := `SELECT EXISTS(SELECT 1 FROM cars WHERE id = $1)`

	err = tx.QueryRow(ctx, findCar, carID).Scan(&exists)

	if err != nil {
		return err
	}

	if !exists {
		return pgx.ErrNoRows
	}

	q := `UPDATE cars SET reg_num = COALESCE(NULLIF($1, ''), reg_num),
                mark = COALESCE(NULLIF($2, ''), mark),
				model = COALESCE(NULLIF($3, ''), model), 
				year = COALESCE(NULLIF($4, 0), year) WHERE id = $5`

	_, err = tx.Exec(ctx, q, input.RegNum, input.Mark, input.Model, input.Year, carID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)

}

func (c *CarRepository) DeleteCar(ctx context.Context, carID int) error {
	ctx, span := c.tracer.Start(ctx, "carRepository.DeleteCar")
	defer span.End()

	tx, err := c.db.Begin(ctx)
	defer tx.Rollback(ctx)

	if err != nil {
		return err
	}

	q := "DELETE FROM cars WHERE id = $1"

	tag, err := tx.Exec(ctx, q, carID)

	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return tx.Commit(ctx)
}
