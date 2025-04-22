package errors

import "errors"

var (
	ErrAlreadyRegistered          = errors.New("user with this username is already registered")
	ErrInsufficientBalance        = errors.New("insufficient balance")
	ErrItemNotFound               = errors.New("item not found")
	ErrUnauthorized               = errors.New("unauthorized")
	ErrUserNotFound               = errors.New("user not found")
	ErrWrongPassword              = errors.New("wrong password provided")
	ErrNoLoginOrPassword          = errors.New("no login or password provided")
	ErrWrongAdminHeader           = errors.New("wrong admin header")
	ErrWrongJSON                  = errors.New("something is wrong in json")
	ErrOrderAlreadyUploadedByUser = errors.New("order already uploaded by this user")
	ErrOrderExists                = errors.New("order already uploaded by another user")
	ErrInsufficientFunds          = errors.New("insufficient balance")
	ErrInvalidOrderNumber         = errors.New("invalid order number")
	ErrTooManyRequests            = errors.New("too many requests to accrual system")
)
