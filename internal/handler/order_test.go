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
	"github.com/Te8va/Gofermarch/internal/handler"
	"github.com/Te8va/Gofermarch/internal/handler/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestOrderHandler_UploadOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProcessor := mocks.NewMockOrderProcessor(ctrl)
	h := handler.NewOrderHandler(mockProcessor)

	userLogin := "testuser"
	ctxWithUser := context.WithValue(context.Background(), domain.UserCtxKey, userLogin)
	validOrderNumber := "79927398713"

	tests := []struct {
		name           string
		body           string
		ctx            context.Context
		mock           func()
		expectedStatus int
	}{
		{
			name: "success",
			body: validOrderNumber,
			ctx:  ctxWithUser,
			mock: func() {
				mockProcessor.EXPECT().
					ProcessOrder(gomock.Any(), validOrderNumber, userLogin).
					Return(domain.StatusNew, nil).
					Times(1)
			},
			expectedStatus: http.StatusAccepted,
		},
		{
			name:           "unauthorized",
			body:           validOrderNumber,
			ctx:            context.Background(),
			mock:           func() {},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid Luhn",
			body:           "1234567890",
			ctx:            ctxWithUser,
			mock:           func() {},
			expectedStatus: http.StatusUnprocessableEntity,
		},
		{
			name: "already uploaded by same user",
			body: validOrderNumber,
			ctx:  ctxWithUser,
			mock: func() {
				mockProcessor.EXPECT().
					ProcessOrder(gomock.Any(), validOrderNumber, userLogin).
					Return(domain.StatusAlreadyUploaded, errors.New("already uploaded")).
					Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "already uploaded by another user",
			body: validOrderNumber,
			ctx:  ctxWithUser,
			mock: func() {
				mockProcessor.EXPECT().
					ProcessOrder(gomock.Any(), validOrderNumber, userLogin).
					Return(domain.StatusConflict, errors.New("conflict")).
					Times(1)
			},
			expectedStatus: http.StatusConflict,
		},
		{
			name: "internal server error",
			body: validOrderNumber,
			ctx:  ctxWithUser,
			mock: func() {
				mockProcessor.EXPECT().
					ProcessOrder(gomock.Any(), validOrderNumber, userLogin).
					Return(domain.StatusNew, errors.New("unexpected")).
					Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "empty order number",
			body:           "   ",
			ctx:            ctxWithUser,
			mock:           func() {},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			req := httptest.NewRequest(http.MethodPost, "/api/user/orders", bytes.NewBufferString(tt.body))
			req = req.WithContext(tt.ctx)
			w := httptest.NewRecorder()

			h.UploadOrder(w, req)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.expectedStatus, res.StatusCode)
		})
	}
}

func TestOrderHandler_GetOrders(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProcessor := mocks.NewMockOrderProcessor(ctrl)
	h := handler.NewOrderHandler(mockProcessor)

	userLogin := "testuser"
	ctxWithUser := context.WithValue(context.Background(), domain.UserCtxKey, userLogin)

	now := time.Now().Truncate(time.Second)
	accrual := 100.0
	orders := []domain.Order{
		{
			Number:     "12345678903",
			Status:     "PROCESSED",
			Accrual:    &accrual,
			UploadedAt: now,
		},
	}

	tests := []struct {
		name           string
		ctx            context.Context
		mock           func()
		expectedStatus int
		expectBody     bool
	}{
		{
			name: "success with orders",
			ctx:  ctxWithUser,
			mock: func() {
				mockProcessor.EXPECT().
					GetOrdersByUser(gomock.Any(), userLogin).
					Return(orders, nil).
					Times(1)
			},
			expectedStatus: http.StatusOK,
			expectBody:     true,
		},
		{
			name: "no orders",
			ctx:  ctxWithUser,
			mock: func() {
				mockProcessor.EXPECT().
					GetOrdersByUser(gomock.Any(), userLogin).
					Return([]domain.Order{}, nil).
					Times(1)
			},
			expectedStatus: http.StatusNoContent,
			expectBody:     false,
		},
		{
			name: "internal error",
			ctx:  ctxWithUser,
			mock: func() {
				mockProcessor.EXPECT().
					GetOrdersByUser(gomock.Any(), userLogin).
					Return(nil, errors.New("db error")).
					Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectBody:     false,
		},
		{
			name:           "unauthorized",
			ctx:            context.Background(),
			mock:           func() {},
			expectedStatus: http.StatusUnauthorized,
			expectBody:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			req = req.WithContext(tt.ctx)
			w := httptest.NewRecorder()

			h.GetOrders(w, req)
			res := w.Result()
			defer res.Body.Close()

			require.Equal(t, tt.expectedStatus, res.StatusCode)

			if tt.expectBody {
				var result []map[string]interface{}
				err := json.NewDecoder(res.Body).Decode(&result)
				require.NoError(t, err)
				require.Len(t, result, len(orders))
			}
		})
	}
}
