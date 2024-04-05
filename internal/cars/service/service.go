package service

import (
	"context"
	"fmt"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"go.uber.org/zap"
)

const (
	carInfoUri = ""
)

type Repository interface {
	CreateCar(ctx context.Context, input domain.CreateCarsRequest) error
}

type Service struct {
	log  *zap.SugaredLogger
	repo Repository
}

func NewService(log *zap.SugaredLogger, repo Repository) *Service {
	return &Service{log: log, repo: repo}
}

func (s *Service) CreateCar(ctx context.Context, input domain.CreateCarsRequest) error {

	err := s.repo.CreateCar(ctx, input)
	if err != nil {
		return fmt.Errorf("create car: %w", err)
	}

	return err
}
