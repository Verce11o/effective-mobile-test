package communicator

import (
	"encoding/json"
	"github.com/Verce11o/effective-mobile-test/internal/domain"
	"github.com/Verce11o/effective-mobile-test/internal/lib/response"
	"net/http"
	"time"
)

type Communicator struct {
	apiEndpoint string
}

func NewCommunicator(apiEndpoint string) *Communicator {
	return &Communicator{apiEndpoint: apiEndpoint}
}

func (c *Communicator) GetCarInfo(regNum string) (domain.Car, error) {
	client := http.Client{Timeout: 3 * time.Second}

	resp, err := client.Get(c.apiEndpoint + "/info?regNum" + regNum)

	if err != nil {
		return domain.Car{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return domain.Car{}, response.ErrGettingCarInfo
	}

	var car domain.Car

	err = json.NewDecoder(resp.Body).Decode(&car)
	if err != nil {
		return domain.Car{}, err
	}

	return car, nil
}
