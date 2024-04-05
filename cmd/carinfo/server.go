package carinfo

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	middleware "github.com/oapi-codegen/nethttp-middleware"
	"go.uber.org/zap"
	"net/http"
)

type ExternalApiHandler struct {
	cars []Car
}

func NewExternalApiHandler() *ExternalApiHandler {
	cars := make([]Car, 10)
	// fill with random data
	for i, _ := range cars {
		patronymic := patronymics[i]
		year := years[i]

		cars[i] = Car{
			Mark:  marks[i],
			Model: models[i],
			Owner: People{
				Name:       names[i],
				Patronymic: &patronymic,
				Surname:    surnames[i],
			},
			RegNum: regNums[i],
			Year:   &year,
		}
	}
	return &ExternalApiHandler{cars: cars}
}

func (h *ExternalApiHandler) GetInfo(w http.ResponseWriter, r *http.Request, params GetInfoParams) {

	// simulate some kind of validation
	if len(params.RegNum) < 9 {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	findCar := func(regNum string) *Car {
		for _, car := range h.cars {
			if car.RegNum == regNum {
				return &car
			}
		}
		return nil
	}

	car := findCar(params.RegNum)

	if car == nil {
		http.Error(w, "Bad request", http.StatusBadRequest) // should be 404, but it's not described in the spec
		return
	}

	render.JSON(w, r, &car)
}

func Run(log *zap.SugaredLogger, addr string) error {
	swagger, err := GetSwagger()
	if err != nil {
		return err
	}

	r := chi.NewRouter()
	r.Use(middleware.OapiRequestValidator(swagger))

	handler := NewExternalApiHandler()

	HandlerFromMux(handler, r)

	log.Infof("External CarsApi running on: %v", addr)
	return http.ListenAndServe(addr, r)

}
