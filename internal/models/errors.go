package models

import "errors"

var ErrRecordNotFound = errors.New("record not found")

var ErrInvalidIndicator = errors.New("invalid indicator")

var ErrInsufficientFunds = errors.New("insufficient funds")
