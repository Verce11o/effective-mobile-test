package handler

import (
	"context"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"github.com/Verce11o/effective-mobile-test/internal/lib/request"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

type Service interface {
	CreateCar(ctx context.Context, input domain.CreateCarsRequest) error
	GetCars(ctx context.Context, input domain.GetCarsRequest) (domain.CarList, error)
}

type Handler struct {
	log     *zap.SugaredLogger
	service Service
}

func NewHandler(log *zap.SugaredLogger, service Service) *Handler {
	return &Handler{log: log, service: service}
}

func (h *Handler) CreateCar(c *gin.Context) {

	var input domain.CreateCarsRequest

	if err := request.Read(c, &input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}
	err := h.service.CreateCar(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *Handler) GetCars(c *gin.Context) {
	var input domain.GetCarsRequest

	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request",
		})
		return
	}

	cars, err := h.service.GetCars(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, cars)

}
