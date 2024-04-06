package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"go.uber.org/zap"
	"net/http"
	"net/url"
)

type Repository interface {
	CreateCars(ctx context.Context, cars []domain.Car) error
	GetCars(ctx context.Context, input domain.GetCarsRequest) (domain.CarList, error)
}

type Service struct {
	log                     *zap.SugaredLogger
	repo                    Repository
	externalCarsApiEndpoint string
}

func NewService(log *zap.SugaredLogger, repo Repository, externalCarsApiEndpoint string) *Service {
	return &Service{log: log, repo: repo, externalCarsApiEndpoint: externalCarsApiEndpoint}
}

func (s *Service) CreateCar(ctx context.Context, input domain.CreateCarsRequest) error {

	client := http.Client{}

	cars := make([]domain.Car, 0, len(input.RegNums))

	for _, regNum := range input.RegNums {
		req, err := http.NewRequest(http.MethodGet, s.externalCarsApiEndpoint+"/info", nil)

		if err != nil {
			return fmt.Errorf("err creating new request: %w", err)
		}

		params := url.Values{}
		params.Add("regNum", regNum)

		req.URL.RawQuery = params.Encode()

		resp, err := client.Do(req)

		if err != nil {
			return fmt.Errorf("send get info request: %w", err)
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			continue
		}

		var car domain.Car

		err = json.NewDecoder(resp.Body).Decode(&car)
		if err != nil {
			return fmt.Errorf("decode response: %w", err)
		}

		cars = append(cars, car)

	}

	err := s.repo.CreateCars(ctx, cars)
	if err != nil {
		return fmt.Errorf("create cars: %w", err)
	}

	return err
}

func (s *Service) GetCars(ctx context.Context, input domain.GetCarsRequest) (domain.CarList, error) {

	carList, err := s.repo.GetCars(ctx, input)
	if err != nil {
		return carList, fmt.Errorf("get cars: %w", err)
	}

	return carList, nil

}
