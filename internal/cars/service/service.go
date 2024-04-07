package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"github.com/Verce11o/effective-mobile-test/internal/lib/response"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type Repository interface {
	CreateCars(ctx context.Context, cars []domain.Car) error
	GetCars(ctx context.Context, input domain.GetCarsRequest) (domain.CarList, error)
	UpdateCar(ctx context.Context, carID int, input domain.UpdateCarsRequest) error
	DeleteCar(ctx context.Context, carID int) error
}

type CacheRepository interface {
	GetCarList(ctx context.Context, cursor string) (*domain.CarList, error)
	SetByIDCtx(ctx context.Context, cursor string, cars domain.CarList) error
	DeleteCarList(ctx context.Context) error
}

type Service struct {
	log                     *zap.SugaredLogger
	repo                    Repository
	cache                   CacheRepository
	externalCarsApiEndpoint string
}

func NewService(log *zap.SugaredLogger, repo Repository, cache CacheRepository, externalCarsApiEndpoint string) *Service {
	return &Service{log: log, repo: repo, cache: cache, externalCarsApiEndpoint: externalCarsApiEndpoint}
}

func (s *Service) CreateCar(ctx context.Context, input domain.CreateCarsRequest) error {

	client := http.Client{}

	cars := make([]domain.Car, 0, len(input.RegNums))

	for _, regNum := range input.RegNums {
		req, err := http.NewRequest(http.MethodGet, s.externalCarsApiEndpoint+"/info", nil)

		if err != nil {
			s.log.Infof("err creating new request: %v", err)
			return err
		}

		params := url.Values{}
		params.Add("regNum", regNum)

		req.URL.RawQuery = params.Encode()

		resp, err := client.Do(req)

		if err != nil {
			s.log.Infof("err send get info request: %v", err)
			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			if resp.StatusCode == http.StatusNotFound {
				return response.ErrNotFound
			}
			continue
		}

		var car domain.Car

		err = json.NewDecoder(resp.Body).Decode(&car)
		if err != nil {
			s.log.Infof("err decode response: %v", err)
			return err
		}

		cars = append(cars, car)

	}

	err := s.repo.CreateCars(ctx, cars)
	if err != nil {
		s.log.Infof("err create cars: %v", err)
		return err
	}

	if err = s.cache.DeleteCarList(ctx); err != nil {
		s.log.Infof("cannot clear cache: %v", err)
	}

	s.log.Debugf("cleared cars cache")

	return err
}

func (s *Service) GetCars(ctx context.Context, input domain.GetCarsRequest) (domain.CarList, error) {

	cachedCars, err := s.cache.GetCarList(ctx, input.Cursor)

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
		s.log.Infof("cannot get cars: %v", err)
		return carList, fmt.Errorf("get cars: %w", err)
	}

	if err = s.cache.SetByIDCtx(ctx, input.Cursor, carList); err != nil {
		s.log.Infof("cannot set product list in redis: %v", err)
	}

	return carList, nil

}

func (s *Service) UpdateCar(ctx context.Context, carID int, input domain.UpdateCarsRequest) error {

	err := s.repo.UpdateCar(ctx, carID, input)
	if err != nil {
		s.log.Infof("cannot update car: %v", err)
		return fmt.Errorf("update car: %w", err)
	}

	if err = s.cache.DeleteCarList(ctx); err != nil {
		s.log.Infof("cannot clear cache: %v", err)
	}

	return nil
}

func (s *Service) DeleteCar(ctx context.Context, carID int) error {
	err := s.repo.DeleteCar(ctx, carID)
	if err != nil {
		s.log.Infof("cannot delete car: %v", err)
		return fmt.Errorf("delete car: %w", err)
	}

	if err = s.cache.DeleteCarList(ctx); err != nil {
		s.log.Infof("cannot clear cache: %v", err)
	}

	return nil
}
