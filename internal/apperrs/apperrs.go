package apperrs

import "errors"

var (
	//при выборке из базы данных
	ErrRecordNotFound   = errors.New("record not found")
	ErrInvalidIndicator = errors.New("invalid indicator")

	//при списании бонусов
	ErrInsufficientFunds = errors.New("insufficient funds")

	//при логине или проверке токена
	ErrUserNotFound = errors.New("user not found")

	//приш регистрации пользователя
	ErrUserAlreadyExists = errors.New("user already exists")

	ErrStatusBadRequest = errors.New("bad request")

	//при обработке акрувалов сингнализтирует о том,
	//что остались необработанные заказы
	ErrBatchIncomplete = errors.New("batch incomplete")

	//при работе с системой бонусов
	ErrTooManyRequests = errors.New("too many requests")

	ErrReadBody    = errors.New("error reading body")
	ErrConvertBody = errors.New("error unmarshalling body")

	ErrAccSystemProblem = errors.New("accrual system problem") //при работе с системой бонусов
)
