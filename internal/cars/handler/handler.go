package handler

import (
	"context"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"github.com/Verce11o/effective-mobile-test/internal/lib/request"
	"github.com/Verce11o/effective-mobile-test/internal/lib/response"
	"github.com/Verce11o/effective-mobile-test/internal/models"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

//go:generate go run github.com/vektra/mockery/v2@v2.42.2 --name=Service
type Service interface {
	CreateCar(ctx context.Context, input domain.CreateCarsRequest) error
	GetCars(ctx context.Context, input domain.GetCarsRequest) (models.CarList, error)
	UpdateCar(ctx context.Context, carID int, input domain.UpdateCarsRequest) error
	DeleteCar(ctx context.Context, carID int) error
}

type Handler struct {
	log     *zap.SugaredLogger
	service Service
	tracer  trace.Tracer
}

func NewHandler(log *zap.SugaredLogger, service Service, tracer trace.Tracer) *Handler {
	return &Handler{log: log, service: service, tracer: tracer}
}

// CreateCars godoc
// @Summary Create new cars
// @Description Create new cars with provided regnums
// @Tags cars
// @Accept  json
// @Produce  json
// @Param   cars body domain.CreateCarsRequest true "Create Cars Request"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /cars [post]
func (h *Handler) CreateCar(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "carHandler.CreateCar")
	defer span.End()

	var input domain.CreateCarsRequest

	if err := request.Read(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	err := h.service.CreateCar(ctx, input)
	if err != nil {

		h.log.Infof("error while creating car: %v", err)
		response.WithHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *Handler) GetCars(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "carHandler.GetCars")
	defer span.End()

	var input domain.GetCarsRequest

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	cars, err := h.service.GetCars(ctx, input)
	if err != nil {
		h.log.Infof("error while getting cars: %v", err)
		response.WithHTTPError(c, err)
		return

	}

	c.JSON(http.StatusOK, cars)

}

func (h *Handler) UpdateCar(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "carHandler.UpdateCar")
	defer span.End()

	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id param is required",
		})
		return
	}

	var input domain.UpdateCarsRequest
	if err := request.Read(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
		})
		return
	}

	carID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id must be integer",
		})
		return
	}

	err = h.service.UpdateCar(ctx, carID, input)
	if err != nil {
		h.log.Infof("error while updating car %v: %v", carID, err)
		response.WithHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})

}

func (h *Handler) DeleteCar(c *gin.Context) {
	ctx, span := h.tracer.Start(c.Request.Context(), "carHandler.DeleteCar")
	defer span.End()

	id := c.Query("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id param is required",
		})
		return
	}

	carID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "id must be integer",
		})
		return
	}

	err = h.service.DeleteCar(ctx, carID)
	if err != nil {
		h.log.Infof("error while deleting car %v: %v", carID, err)
		response.WithHTTPError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
