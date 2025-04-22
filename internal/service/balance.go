package service

import (
	"context"
	"time"

	"github.com/Te8va/Gofermarch/internal/domain"
	appErrors "github.com/Te8va/Gofermarch/internal/errors"
)

type BalanceServ interface {
	GetBalance(ctx context.Context, login string) (domain.Balance, error)
	SaveWithdrawal(ctx context.Context, w domain.Withdrawal) error
	GetWithdrawalsByUser(ctx context.Context, login string) ([]domain.Withdrawal, error)
}

type BalanceService struct {
	srv BalanceServ
}

func NewBalanceService(srv BalanceServ) *BalanceService {
	return &BalanceService{srv: srv}
}

func (s *BalanceService) GetUserBalance(ctx context.Context, login string) (domain.Balance, error) {
	return s.srv.GetBalance(ctx, login)
}

func (s *BalanceService) WithdrawBalance(ctx context.Context, login, order string, sum float64) error {
	balance, err := s.srv.GetBalance(ctx, login)
	if err != nil {
		return err
	}

	if sum > balance.Current {
		return appErrors.ErrInsufficientFunds
	}

	return s.srv.SaveWithdrawal(ctx, domain.Withdrawal{
		Login:       login,
		Order:       order,
		Sum:         sum,
		ProcessedAt: time.Now(),
	})
}

func (s *BalanceService) GetUserWithdrawals(ctx context.Context, login string) ([]domain.Withdrawal, error) {
	return s.srv.GetWithdrawalsByUser(ctx, login)
}
