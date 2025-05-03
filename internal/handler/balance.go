package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Te8va/Gofermarch/internal/domain"
	appErrors "github.com/Te8va/Gofermarch/internal/errors"
	"github.com/Te8va/Gofermarch/internal/utils"
)

//go:generate mockgen -source=balance.go -destination=mocks/mock_authorization.go -package=mocks

type Balance interface {
	GetUserBalance(ctx context.Context, login string) (domain.Balance, error)
	WithdrawBalance(ctx context.Context, login, order string, sum float64) error
	GetUserWithdrawals(ctx context.Context, login string) ([]domain.Withdrawal, error)
}

type BalanceHandler struct {
	bal Balance
}

func NewBalanceHandler(bal Balance) *BalanceHandler {
	return &BalanceHandler{bal: bal}
}

func (h *BalanceHandler) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(domain.UserCtxKey).(string)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	balance, err := h.bal.GetUserBalance(r.Context(), login)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := struct {
		Current   float64 `json:"current"`
		Withdrawn float64 `json:"withdrawn"`
	}{
		Current:   balance.Current,
		Withdrawn: balance.Withdrawn,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *BalanceHandler) WithdrawBalance(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(domain.UserCtxKey).(string)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if !utils.IsValidLuhn(req.Order) {
		http.Error(w, "Invalid order number format", http.StatusUnprocessableEntity)
		return
	}

	err := h.bal.WithdrawBalance(r.Context(), login, req.Order, req.Sum)
	if err != nil {
		switch {
		case errors.Is(err, appErrors.ErrInsufficientFunds):
			http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
		case errors.Is(err, appErrors.ErrInvalidOrderNumber):
			http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *BalanceHandler) GetUserWithdrawals(w http.ResponseWriter, r *http.Request) {
	login, ok := r.Context().Value(domain.UserCtxKey).(string)
	if !ok || login == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	withdrawals, err := h.bal.GetUserWithdrawals(r.Context(), login)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	type response struct {
		Order       string  `json:"order"`
		Sum         float64 `json:"sum"`
		ProcessedAt string  `json:"processed_at"`
	}

	resp := make([]response, 0, len(withdrawals))
	for _, w := range withdrawals {
		resp = append(resp, response{
			Order:       w.Order,
			Sum:         w.Sum,
			ProcessedAt: w.ProcessedAt.Format(time.RFC3339),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
