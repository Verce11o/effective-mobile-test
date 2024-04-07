package service

import (
	"context"
	repoMock "github.com/Verce11o/effective-mobile-test/internal/cars/service/mocks"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"github.com/Verce11o/effective-mobile-test/internal/lib/logger"
	"github.com/Verce11o/effective-mobile-test/internal/lib/response"
	"github.com/Verce11o/effective-mobile-test/internal/lib/tracer"
	"github.com/Verce11o/effective-mobile-test/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestService_CreateCar(t *testing.T) {

	type args struct {
		ctx   context.Context
		input domain.CreateCarsRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				input: domain.CreateCarsRequest{
					RegNums: []string{"J623FP555", "Z407GI541"},
				},
			},
		},
		{
			name: "invalid regnum",
			args: args{
				ctx: context.Background(),
				input: domain.CreateCarsRequest{
					RegNums: []string{"hello world"},
				},
			},
			wantErr: response.ErrGettingCarInfo,
		},
	}

	log := logger.NewMockLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			repo := repoMock.NewRepository(t)
			cache := repoMock.NewCacheRepository(t)
			communicator := repoMock.NewApiCommunicator(t)

			repo.On("CreateCars", mock.Anything, mock.AnythingOfType("[]domain.Car")).Maybe().Return(nil)
			cache.On("DeleteCarList", mock.Anything).Maybe().Return(nil)

			communicator.On("GetCarInfo", mock.AnythingOfType("string")).Return(domain.Car{}, nil)

			s := &Service{
				log:          log,
				repo:         repo,
				cache:        cache,
				tracer:       tracer.InitTracer(tt.args.ctx, "", ""),
				communicator: communicator,
			}

			err := s.CreateCar(tt.args.ctx, tt.args.input)

			if err != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func TestService_GetCars(t *testing.T) {

	type args struct {
		ctx   context.Context
		input domain.GetCarsRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				input: domain.GetCarsRequest{
					Cursor: "",
					RegNum: "E387IK307",
					Mark:   "Chevrolet",
					Model:  "Camaro",
				},
			},
			wantErr: nil,
		},
	}

	log := logger.NewMockLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repoMock.NewRepository(t)
			cache := repoMock.NewCacheRepository(t)
			communicator := repoMock.NewApiCommunicator(t)

			cache.On("GetCarList", mock.Anything, mock.AnythingOfType("string")).Return(nil, nil)

			repo.On("GetCars", mock.Anything, mock.AnythingOfType("domain.GetCarsRequest")).Maybe().Return(models.CarList{}, nil)
			cache.On("SetByIDCtx", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("models.CarList")).Maybe().Return(nil)

			s := &Service{
				log:          log,
				repo:         repo,
				cache:        cache,
				tracer:       tracer.InitTracer(tt.args.ctx, "", ""),
				communicator: communicator,
			}

			_, err := s.GetCars(tt.args.ctx, tt.args.input)
			if err != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}

func TestService_UpdateCar(t *testing.T) {
	type args struct {
		ctx   context.Context
		carID int
		input domain.UpdateCarsRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				carID: 1,
				input: domain.UpdateCarsRequest{},
			},

			wantErr: nil,
		},
	}

	log := logger.NewMockLogger()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := repoMock.NewRepository(t)
			cache := repoMock.NewCacheRepository(t)
			communicator := repoMock.NewApiCommunicator(t)

			repo.On("UpdateCar", mock.Anything, mock.AnythingOfType("int"), mock.AnythingOfType("domain.UpdateCarsRequest")).Return(nil).Once()
			cache.On("DeleteCarList", mock.Anything).Maybe().Return(nil)

			s := &Service{
				log:          log,
				repo:         repo,
				cache:        cache,
				tracer:       tracer.InitTracer(tt.args.ctx, "", ""),
				communicator: communicator,
			}
			err := s.UpdateCar(tt.args.ctx, tt.args.carID, tt.args.input)

			if err != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			}

		})
	}
}

func TestService_DeleteCar(t *testing.T) {
	type args struct {
		ctx   context.Context
		carID int
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				carID: 1,
			},
		},
	}

	log := logger.NewMockLogger()

	for _, tt := range tests {
		repo := repoMock.NewRepository(t)
		cache := repoMock.NewCacheRepository(t)
		communicator := repoMock.NewApiCommunicator(t)

		repo.On("DeleteCar", mock.Anything, mock.AnythingOfType("int")).Return(nil).Once()
		cache.On("DeleteCarList", mock.Anything).Maybe().Return(nil)

		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				log:          log,
				repo:         repo,
				cache:        cache,
				tracer:       tracer.InitTracer(tt.args.ctx, "", ""),
				communicator: communicator,
			}

			err := s.DeleteCar(tt.args.ctx, tt.args.carID)
			if err != nil {
				require.EqualError(t, err, tt.wantErr.Error())
			}
		})
	}
}
