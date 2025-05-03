package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Te8va/Gofermarch/internal/domain"
	appErrors "github.com/Te8va/Gofermarch/internal/errors"
	"github.com/Te8va/Gofermarch/internal/handler"
	"github.com/Te8va/Gofermarch/internal/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestBalanceHandler_GetUserBalance(t *testing.T) {
	type args struct {
		login  string
		ctxSet bool
	}
	type mockResp struct {
		balance domain.Balance
		err     error
	}
	tests := []struct {
		name       string
		args       args
		mockResp   mockResp
		wantStatus int
	}{
		{
			name:       "unauthorized",
			args:       args{ctxSet: false},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "success",
			args:       args{ctxSet: true, login: "user"},
			mockResp:   mockResp{balance: domain.Balance{Current: 100, Withdrawn: 20}},
			wantStatus: http.StatusOK,
		},
		{
			name:       "internal error",
			args:       args{ctxSet: true, login: "user"},
			mockResp:   mockResp{err: errors.New("db")},
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockBalance(ctrl)
			h := handler.NewBalanceHandler(mock)

			ctx := context.Background()
			if tt.args.ctxSet {
				ctx = context.WithValue(ctx, domain.UserCtxKey, tt.args.login)
				mock.EXPECT().
					GetUserBalance(gomock.Any(), tt.args.login).
					Return(tt.mockResp.balance, tt.mockResp.err).
					Times(1)
			}

			req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil).WithContext(ctx)
			w := httptest.NewRecorder()

			h.GetUserBalance(w, req)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestBalanceHandler_WithdrawBalance(t *testing.T) {
	type requestBody struct {
		Order string  `json:"order"`
		Sum   float64 `json:"sum"`
	}
	tests := []struct {
		name       string
		login      string
		body       any
		mockError  error
		wantStatus int
	}{
		{
			name:       "unauthorized",
			body:       requestBody{"79927398713", 10},
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "invalid json",
			login:      "user",
			body:       "invalid",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "invalid Luhn",
			login:      "user",
			body:       requestBody{"123456", 10},
			wantStatus: http.StatusUnprocessableEntity,
		},
		{
			name:       "insufficient funds",
			login:      "user",
			body:       requestBody{"79927398713", 10},
			mockError:  appErrors.ErrInsufficientFunds,
			wantStatus: http.StatusPaymentRequired,
		},
		{
			name:       "internal error",
			login:      "user",
			body:       requestBody{"79927398713", 10},
			mockError:  errors.New("unexpected"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "success",
			login:      "user",
			body:       requestBody{"79927398713", 10},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockBalance(ctrl)
			h := handler.NewBalanceHandler(mock)

			var buf bytes.Buffer
			switch b := tt.body.(type) {
			case string:
				buf.WriteString(b)
			default:
				_ = json.NewEncoder(&buf).Encode(b)
			}

			ctx := context.Background()
			if tt.login != "" {
				ctx = context.WithValue(ctx, domain.UserCtxKey, tt.login)
				if b, ok := tt.body.(requestBody); ok && b.Order == "79927398713" {
					mock.EXPECT().
						WithdrawBalance(gomock.Any(), tt.login, b.Order, b.Sum).
						Return(tt.mockError).
						Times(1)
				}
			}

			req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", &buf).WithContext(ctx)
			w := httptest.NewRecorder()

			h.WithdrawBalance(w, req)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}

func TestBalanceHandler_GetUserWithdrawals(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name         string
		login        string
		mockResult   []domain.Withdrawal
		mockError    error
		wantStatus   int
	}{
		{
			name:       "unauthorized",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "internal error",
			login:      "user",
			mockError:  errors.New("db error"),
			wantStatus: http.StatusInternalServerError,
		},
		{
			name:       "no withdrawals",
			login:      "user",
			wantStatus: http.StatusNoContent,
		},
		{
			name:       "success",
			login:      "user",
			mockResult: []domain.Withdrawal{{Order: "123", Sum: 10.0, ProcessedAt: now}},
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mock := mocks.NewMockBalance(ctrl)
			h := handler.NewBalanceHandler(mock)

			ctx := context.Background()
			if tt.login != "" {
				ctx = context.WithValue(ctx, domain.UserCtxKey, tt.login)
				mock.EXPECT().
					GetUserWithdrawals(gomock.Any(), tt.login).
					Return(tt.mockResult, tt.mockError).
					Times(1)
			}

			req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil).WithContext(ctx)
			w := httptest.NewRecorder()

			h.GetUserWithdrawals(w, req)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.wantStatus, res.StatusCode)
		})
	}
}
