package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/Te8va/Gofermarch/internal/domain"
	"github.com/Te8va/Gofermarch/internal/utils"
)

//go:generate mockgen -source=order.go -destination=mocks/mock_order.go -package=mocks

type OrderProcessor interface {
	ProcessOrder(ctx context.Context, number, login string) (domain.OrderStatus, error)
	GetOrdersByUser(ctx context.Context, login string) ([]domain.Order, error)
}

type OrderHandler struct {
	order OrderProcessor
}

func NewOrderHandler(order OrderProcessor) *OrderHandler {
	return &OrderHandler{order: order}
}

func (h *OrderHandler) UploadOrder(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(domain.UserCtxKey).(string)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	orderNumber := strings.TrimSpace(string(body))
	if orderNumber == "" {
		http.Error(w, "Empty order number", http.StatusBadRequest)
		return
	}

	if !utils.IsValidLuhn(orderNumber) {
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	status, err := h.order.ProcessOrder(r.Context(), orderNumber, login)
	if err != nil {
		switch status {
		case domain.StatusAlreadyUploaded:
			w.WriteHeader(http.StatusOK)
		case domain.StatusConflict:
			http.Error(w, "Order already uploaded by another user", http.StatusConflict)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(domain.UserCtxKey).(string)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.order.GetOrdersByUser(r.Context(), login)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type response struct {
		Number     string   `json:"number"`
		Status     string   `json:"status"`
		Accrual    *float64 `json:"accrual,omitempty"`
		UploadedAt string   `json:"uploaded_at"`
	}

	resp := make([]response, 0, len(orders))
	for _, o := range orders {
		resp = append(resp, response{
			Number:     o.Number,
			Status:     o.Status,
			Accrual:    o.Accrual,
			UploadedAt: o.UploadedAt.Local().Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
