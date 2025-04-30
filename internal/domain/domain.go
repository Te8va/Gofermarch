package domain

import (
	"time"
)

type ctxKey string

const UserCtxKey ctxKey = "user"

type Order struct {
	Number     string    `json:"number"`
	Login      string    `json:"-"`
	UploadedAt time.Time `json:"uploaded_at"`
	Status     string    `json:"status"`
	Accrual    *float64  `json:"accrual,omitempty"`
}

type OrderStatus string

const (
	StatusNew             OrderStatus = "NEW"
	StatusProcessing      OrderStatus = "PROCESSING"
	StatusInvalid         OrderStatus = "INVALID"
	StatusProcessed       OrderStatus = "PROCESSED"
	StatusAlreadyUploaded OrderStatus = "ALREADY_UPLOADED"
	StatusConflict        OrderStatus = "CONFLICT"
	StatusInternal        OrderStatus = "INTERNAL"
)

type Balance struct {
	Current   float64
	Withdrawn float64
}

type Withdrawal struct {
	Login       string
	Order       string
	Sum         float64
	ProcessedAt time.Time
}
