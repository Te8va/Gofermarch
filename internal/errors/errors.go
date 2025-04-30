package errors

import "errors"

var (
	ErrAlreadyRegistered          = errors.New("user with this username is already registered")
	ErrUserNotFound               = errors.New("user not found")
	ErrWrongPassword              = errors.New("wrong password provided")
	ErrOrderExists                = errors.New("order already uploaded by another user")
	ErrInsufficientFunds          = errors.New("insufficient balance")
	ErrInvalidOrderNumber         = errors.New("invalid order number")
	ErrTooManyRequests            = errors.New("too many requests to accrual system")
)
