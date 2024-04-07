package handler

import (
	"bytes"
	"encoding/json"
	"github.com/Verce11o/effective-mobile-test/internal/cars/handler/mocks"
	"github.com/Verce11o/effective-mobile-test/internal/lib/logger"
	"github.com/Verce11o/effective-mobile-test/internal/lib/tracer"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateCar(t *testing.T) {
	type args struct {
		input map[string]any
	}
	tests := []struct {
		name       string
		args       args
		statusCode int
		wantErr    error
	}{
		{
			name: "create car",
			args: args{
				input: map[string]any{
					"regNums": []string{"E387IK307"},
				},
			},
			statusCode: http.StatusOK,
		},
	}

	log := logger.NewMockLogger()
	for _, tt := range tests {
		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		ctx.Request = &http.Request{
			Method: http.MethodPost,
			Header: make(http.Header),
		}

		MockJsonPost(ctx, tt.args.input)

		serviceMock := mocks.NewService(t)

		t.Run(tt.name, func(t *testing.T) {
			h := &Handler{
				log:     log,
				service: serviceMock,
				tracer:  tracer.InitTracer(ctx.Request.Context(), "", ""),
			}

			serviceMock.On("CreateCar", mock.Anything, mock.AnythingOfType("domain.CreateCarsRequest")).Return(nil)
			h.CreateCar(ctx)

			//
			assert.EqualValues(t, tt.statusCode, w.Code)
		})
	}
}

func MockJsonPost(c *gin.Context, body interface{}) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}
