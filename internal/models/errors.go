package models

import "errors"

var ErrRecordNotFound = errors.New("record not found")

var ErrInvalidIndicator = errors.New("invalid indicator")

var ErrInsufficientFunds = errors.New("insufficient funds")

var ErrUserNotFound = errors.New("user not found")

var ErrUserExists = errors.New("user already exists")

var ErrStatusBadRequest = errors.New("bad request")
