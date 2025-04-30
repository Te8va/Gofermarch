package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Te8va/Gofermarch/internal/domain"
	appErrors "github.com/Te8va/Gofermarch/internal/errors"
	"github.com/Te8va/Gofermarch/internal/service/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestOrderService_ProcessOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderServ(ctrl)
	orderService := NewOrderService(mockRepo, "http://localhost")

	ctx := context.Background()
	number := "order123"
	login := "user1"

	testCases := []struct {
		name           string
		mockSetup      func()
		expectedStatus domain.OrderStatus
		expectedError  error
	}{
		{
			name: "order already uploaded by the same user",
			mockSetup: func() {
				mockRepo.EXPECT().GetOrder(ctx, number).Return(domain.Order{Login: login}, nil)
			},
			expectedStatus: domain.StatusAlreadyUploaded,
			expectedError:  errors.New("order already uploaded by the same user"),
		},
		{
			name: "successful order creation",
			mockSetup: func() {
				mockRepo.EXPECT().GetOrder(ctx, number).Return(domain.Order{}, appErrors.ErrOrderExists)
				mockRepo.EXPECT().SaveOrder(ctx, gomock.Any()).Return(nil)
			},
			expectedStatus: domain.StatusNew,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			status, err := orderService.ProcessOrder(ctx, number, login)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedStatus, status)
			}
		})
	}
}

func TestOrderService_GetOrdersByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockOrderServ(ctrl)
	orderService := NewOrderService(mockRepo, "http://localhost")

	ctx := context.Background()
	login := "user1"

	testCases := []struct {
		name           string
		mockSetup      func()
		expectedOrders []domain.Order
		expectedError  error
	}{
		{
			name: "successful fetch orders",
			mockSetup: func() {
				mockRepo.EXPECT().GetOrdersByUser(ctx, login).Return([]domain.Order{
					{Number: "order123", Login: login, Status: string(domain.StatusNew)},
				}, nil)
			},
			expectedOrders: []domain.Order{
				{Number: "order123", Login: login, Status: string(domain.StatusNew)},
			},
			expectedError: nil,
		},
		{
			name: "error fetching orders",
			mockSetup: func() {
				mockRepo.EXPECT().GetOrdersByUser(ctx, login).Return(nil, errors.New("db error"))
			},
			expectedOrders: nil,
			expectedError:  errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			orders, err := orderService.GetOrdersByUser(ctx, login)

			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedOrders, orders)
			}
		})
	}
}
