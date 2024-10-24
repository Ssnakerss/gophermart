package accrualsrv

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ssnakerss/gophermart/internal/apperrs"
	"github.com/Ssnakerss/gophermart/internal/models"
	"github.com/Ssnakerss/gophermart/internal/types"
)

const APIURI = "/api/orders"

type HTTPAccrualsystem struct {
	endPoint string
}

func NewHTTPAccrualsystem(endPoint string) *HTTPAccrualsystem {
	return &HTTPAccrualsystem{endPoint: endPoint}
}

func (ha *HTTPAccrualsystem) GetAccrual(order types.OrderNum) (*models.AccrualResponse, error) {
	url := ha.endPoint + APIURI + "/" + order.String()

	response, err := http.Get(url)
	if err != nil {
		slog.Error("HTTPAccrualsystem.GetAccrual: " + err.Error())
		return nil, err
	}
	defer response.Body.Close()

	//провери на статус ответа 429
	//надо скорректировать размер батча, чтобы не было ошибок
	if response.StatusCode == http.StatusTooManyRequests {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return nil, apperrs.ErrReadBody
		}
		//пытаемся из текста ошибки достать максимальное количество запросов
		p := strings.Fields(string(body))
		if len(p) < 4 {
			return nil, apperrs.ErrConvertBody
		}
		maxReqNum, err := strconv.ParseInt(p[3], 10, 32)
		if err != nil {
			return nil, apperrs.ErrConvertBody
		}
		//передаем максимальное количество запросов в ответе со статусом "E"
		ar := models.AccrualResponse{
			Order:   0,
			Status:  "E",
			Accrual: types.Accrual(maxReqNum),
		}
		return &ar, apperrs.ErrTooManyRequests
	}
	//все прочие ошибки не обрабатываем
	if response.StatusCode != http.StatusOK {
		return nil, errors.New(response.Status) // TODO: change this to your own error
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, apperrs.ErrReadBody
	}

	ar := models.AccrualResponse{}

	err = json.Unmarshal(body, &ar)

	if err != nil {
		return nil, apperrs.ErrConvertBody
	}

	return &ar, nil

}
