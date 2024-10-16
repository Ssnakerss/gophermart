package accrual

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/types"
)

const ServerAddr = "localhost:8080" // TODO: change this to your own endpoint
const ApiUri = "/api/orders"

type HTTPAccrualsystem struct {
	endPoint string
}

func NewHTTPAccrualsystem(endPoint string) *HTTPAccrualsystem {
	return &HTTPAccrualsystem{endPoint: endPoint}
}

func (ha *HTTPAccrualsystem) GetAccrual(order types.OrderNum) (*models.AccrualResponse, error) {
	url := ha.endPoint + "/" + order.String()
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status) // TODO: change this to your own error
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.New("error reading body") // TODO: change this to your own error
	}
	ar := models.AccrualResponse{}

	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, errors.New("error unmarshalling body") // TODO: change this to your own error

	}

	return &ar, nil
	// TODO: implement this

}
