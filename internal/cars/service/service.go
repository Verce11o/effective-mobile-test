package service

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"github.com/Verce11o/effective-mobile-test/internal/models"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=Repository
type Repository interface {
	CreateCars(ctx context.Context, cars []domain.Car) error
	GetCars(ctx context.Context, input domain.GetCarsRequest) (models.CarList, error)
	UpdateCar(ctx context.Context, carID int, input domain.UpdateCarsRequest) error
	DeleteCar(ctx context.Context, carID int) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=CacheRepository
type CacheRepository interface {
	GetCarList(ctx context.Context, hash string) (*models.CarList, error)
	SetByIDCtx(ctx context.Context, cursor string, cars models.CarList) error
	DeleteCarList(ctx context.Context) error
}

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=ApiCommunicator
type ApiCommunicator interface {
	GetCarInfo(regNum string) (domain.Car, error)
}

type Service struct {
	log          *zap.SugaredLogger
	repo         Repository
	cache        CacheRepository
	tracer       trace.Tracer
	communicator ApiCommunicator
}

func NewService(log *zap.SugaredLogger, repo Repository, cache CacheRepository, tracer trace.Tracer, communicator ApiCommunicator) *Service {
	return &Service{log: log, repo: repo, cache: cache, tracer: tracer, communicator: communicator}
}

func (s *Service) CreateCar(ctx context.Context, input domain.CreateCarsRequest) error {
	ctx, span := s.tracer.Start(ctx, "carService.CreateCar")
	defer span.End()

	cars := make([]domain.Car, 0, len(input.RegNums))

	span.AddEvent("iterate over reg nums")

	for _, regNum := range input.RegNums {

		span.AddEvent("call external api")

		car, err := s.communicator.GetCarInfo(regNum)
		if err != nil {
			return err
		}

		cars = append(cars, car)

	}

	span.AddEvent("call postgres repo")
	err := s.repo.CreateCars(ctx, cars)

	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Infof("err create cars: %v", err)
		return err
	}

	span.AddEvent("clear redis cache")
	if err = s.cache.DeleteCarList(ctx); err != nil {
		s.log.Infof("cannot clear cache: %v", err)
	}

	s.log.Debugf("cleared cars cache")
	return err
}

func (s *Service) GetCars(ctx context.Context, input domain.GetCarsRequest) (models.CarList, error) {
	ctx, span := s.tracer.Start(ctx, "carService.GetCars")
	defer span.End()

	paramsBytes, err := json.Marshal(input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Infof("err marshalling params: %v", err)
		return models.CarList{}, err
	}

	hash := sha256.Sum256([]byte(string(paramsBytes) + input.Cursor))
	hashStr := fmt.Sprintf("%x", hash)

	cachedCars, err := s.cache.GetCarList(ctx, hashStr)

	if err != nil {
		s.log.Debugf("cannot get car list in redis: %v", err)
	}

	if cachedCars != nil {
		s.log.Debugf("get cached cars")
		return *cachedCars, nil
	}

	s.log.Debugf("cache not found")

	carList, err := s.repo.GetCars(ctx, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Infof("cannot get cars: %v", err)
		return carList, fmt.Errorf("get cars: %w", err)
	}

	if err = s.cache.SetByIDCtx(ctx, hashStr, carList); err != nil {
		s.log.Infof("cannot set product list in redis: %v", err)
	}

	return carList, nil

}

func (s *Service) UpdateCar(ctx context.Context, carID int, input domain.UpdateCarsRequest) error {
	ctx, span := s.tracer.Start(ctx, "carService.UpdateCar")
	defer span.End()

	err := s.repo.UpdateCar(ctx, carID, input)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Infof("cannot update car: %v", err)
		return fmt.Errorf("update car: %w", err)
	}

	if err = s.cache.DeleteCarList(ctx); err != nil {
		s.log.Infof("cannot clear cache: %v", err)
	}

	return nil
}

func (s *Service) DeleteCar(ctx context.Context, carID int) error {
	ctx, span := s.tracer.Start(ctx, "carService.DeleteCar")
	defer span.End()

	err := s.repo.DeleteCar(ctx, carID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		s.log.Infof("cannot delete car: %v", err)
		return fmt.Errorf("delete car: %w", err)
	}

	if err = s.cache.DeleteCarList(ctx); err != nil {
		s.log.Infof("cannot clear cache: %v", err)
	}

	return nil
}
