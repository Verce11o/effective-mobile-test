package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"net/http"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
	ErrNotFound       = errors.New("not found")
	ErrGettingCarInfo = errors.New("error getting car info")
)

func MapHTTPError(err error) (int, string) {
	switch {
	case errors.Is(err, ErrInvalidRequest):
		return http.StatusBadRequest, "bad request"
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, "not found"
	case errors.Is(err, pgx.ErrNoRows):
		return http.StatusNotFound, "not found"
	case errors.Is(err, ErrGettingCarInfo):
		return http.StatusBadRequest, "error getting car info"
	}

	return http.StatusInternalServerError, "server error"

}

func WithHTTPError(c *gin.Context, err error) {
	status, message := MapHTTPError(err)
	c.JSON(status, gin.H{
		"message": message,
	})
	return
}
