package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/Te8va/Gofermarch/internal/domain"
	appErrors "github.com/Te8va/Gofermarch/internal/errors"
	"github.com/Te8va/Gofermarch/internal/service/mocks"
)

func TestBalanceService_GetUserBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBalanceServ(ctrl)
	balanceService := NewBalanceService(mockRepo)

	ctx := context.Background()
	login := "user1"

	testCases := []struct {
		name            string
		mockSetup       func()
		expectedBalance domain.Balance
		expectedError   error
	}{
		{
			name: "successful balance fetch",
			mockSetup: func() {
				mockRepo.EXPECT().GetBalance(ctx, login).Return(domain.Balance{Current: 100.0}, nil)
			},
			expectedBalance: domain.Balance{Current: 100.0},
			expectedError:   nil,
		},
		{
			name: "error fetching balance",
			mockSetup: func() {
				mockRepo.EXPECT().GetBalance(ctx, login).Return(domain.Balance{}, errors.New("db error"))
			},
			expectedBalance: domain.Balance{},
			expectedError:   errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			balance, err := balanceService.GetUserBalance(ctx, login)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.expectedBalance, balance)
			}
		})
	}
}

func TestBalanceService_WithdrawBalance(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBalanceServ(ctrl)
	balanceService := NewBalanceService(mockRepo)

	ctx := context.Background()
	login := "user1"
	order := "order123"
	balance := domain.Balance{Current: 100.0}
	insufficientBalance := 200.0
	sum := 50.0

	testCases := []struct {
		name          string
		mockSetup     func()
		amount        float64
		expectedError error
	}{
		{
			name: "successful withdrawal",
			mockSetup: func() {
				mockRepo.EXPECT().GetBalance(ctx, login).Return(balance, nil)
				mockRepo.EXPECT().SaveWithdrawal(ctx, gomock.Any()).Return(nil)
			},
			amount:        sum,
			expectedError: nil,
		},
		{
			name: "insufficient funds",
			mockSetup: func() {
				mockRepo.EXPECT().GetBalance(ctx, login).Return(domain.Balance{Current: insufficientBalance}, nil)
			},
			amount:        500.0,
			expectedError: appErrors.ErrInsufficientFunds,
		},
		{
			name: "error fetching balance",
			mockSetup: func() {
				mockRepo.EXPECT().GetBalance(ctx, login).Return(domain.Balance{}, errors.New("db error"))
			},
			amount:        sum,
			expectedError: errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			err := balanceService.WithdrawBalance(ctx, login, order, tc.amount)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestBalanceService_GetUserWithdrawals(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockBalanceServ(ctrl)
	balanceService := NewBalanceService(mockRepo)

	ctx := context.Background()
	login := "user1"

	testCases := []struct {
		name                string
		mockSetup           func()
		expectedWithdrawals []domain.Withdrawal
		expectedError       error
	}{
		{
			name: "successful withdrawal fetch",
			mockSetup: func() {
				mockRepo.EXPECT().GetWithdrawalsByUser(ctx, login).Return([]domain.Withdrawal{
					{Login: login, Order: "order123", Sum: 50.0, ProcessedAt: time.Now()},
				}, nil)
			},
			expectedWithdrawals: []domain.Withdrawal{
				{Login: login, Order: "order123", Sum: 50.0, ProcessedAt: time.Now().Truncate(time.Second)}, // Truncate to second
			},
			expectedError: nil,
		},
		{
			name: "error fetching withdrawals",
			mockSetup: func() {
				mockRepo.EXPECT().GetWithdrawalsByUser(ctx, login).Return(nil, errors.New("db error"))
			},
			expectedWithdrawals: nil,
			expectedError:       errors.New("db error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockSetup()
			withdrawals, err := balanceService.GetUserWithdrawals(ctx, login)
			if tc.expectedError != nil {
				require.Error(t, err)
				require.Contains(t, err.Error(), tc.expectedError.Error())
			} else {
				require.NoError(t, err)

				for i := range withdrawals {
					withdrawals[i].ProcessedAt = withdrawals[i].ProcessedAt.Truncate(time.Second)
				}

				require.Equal(t, tc.expectedWithdrawals, withdrawals)
			}
		})
	}
}
