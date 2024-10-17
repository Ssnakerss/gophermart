package models

import "errors"

//при выборке из базы данных
var ErrRecordNotFound = errors.New("record not found")

var ErrInvalidIndicator = errors.New("invalid indicator")

//при списании бонусов
var ErrInsufficientFunds = errors.New("insufficient funds")

//при логине или проверке токена
var ErrUserNotFound = errors.New("user not found")

//приш регистрации пользователя
var ErrUserAlreadyExists = errors.New("user already exists")

var ErrStatusBadRequest = errors.New("bad request")

//при обработке акрувалов сингнализтирует о том,
//что остались необработанные заказы
var ErrBatchIncomplete = errors.New("batch incomplete")

//при работе с системой бонусов
var ErrTooManyRequests = errors.New("too many requests")

var ErrReadBody = errors.New("error reading body")
var ErrConvertBody = errors.New("error unmarshalling body")

var ErrAccSystemProblem = errors.New("accrual system problem") //при работе с системой бонусов
