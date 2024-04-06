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
	"time"
)

const (
	paginationLimit = 10
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
		_, err = tx.Exec(ctx, "INSERT INTO cars (reg_num, mark, model, year, ownerid) VALUES ($1, $2, $3, $4, $5)", car.RegNum, car.Mark, car.Model, car.Year, ownerID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (c *CarRepository) GetCars(ctx context.Context, input domain.GetCarsRequest) (domain.CarList, error) {
	var createdAt time.Time
	var id int
	var err error

	if input.Cursor != "" {
		_, id, err = pagination.DecodeCursor(input.Cursor)

		if err != nil {
			return domain.CarList{}, err
		}
	}

	query := sq.Select("cars.id, cars.reg_num, cars.mark, cars.model, cars.year, cars.created_at, o.id, o.name, o.surname, o.patronymic").
		From("cars").InnerJoin("owners o on o.id = cars.ownerid").
		OrderBy("cars.created_at, cars.id").Limit(paginationLimit)

	if id != 0 {
		query = query.Where(sq.Gt{"cars.id": id})
	}

	if !createdAt.IsZero() {
		query = query.Where(sq.GtOrEq{"cars.created_at": createdAt})
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
		return domain.CarList{}, err
	}

	rows, err := c.db.Query(ctx, sql, args...)
	if err != nil {
		return domain.CarList{}, err
	}

	defer rows.Close()

	cars := make([]models.Car, 0)

	for rows.Next() {
		var car models.Car

		err = rows.Scan(&car.ID, &car.RegNum, &car.Mark, &car.Model, &car.Year, &car.CreatedAt, &car.Owner.ID, &car.Owner.Name, &car.Owner.Surname, &car.Owner.Patronymic)
		if err != nil {
			return domain.CarList{}, err
		}

		cars = append(cars, car)

	}

	var nextCursor string

	if len(cars) > 0 {
		last := cars[len(cars)-1]
		nextCursor = pagination.EncodeCursor(last.CreatedAt, last.ID)

	}
	return domain.CarList{
		Cursor: nextCursor,
		Cars:   cars,
		Total:  len(cars),
	}, nil
}

func (c *CarRepository) UpdateCar(ctx context.Context, carID int, input domain.UpdateCarsRequest) error {

	tx, err := c.db.Begin(ctx)

	if err != nil {
		return fmt.Errorf("error begin transaction: %w", err)
	}

	defer tx.Rollback(ctx)

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
