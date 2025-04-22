package service

import (
	"context"
	"time"

	"github.com/Te8va/Gofermarch/internal/domain"
	appErrors "github.com/Te8va/Gofermarch/internal/errors"
)

type BalanceRepository interface {
	GetBalance(ctx context.Context, login string) (domain.Balance, error)
	SaveWithdrawal(ctx context.Context, w domain.Withdrawal) error
	GetWithdrawalsByUser(ctx context.Context, login string) ([]domain.Withdrawal, error)
}

type BalanceService struct {
	repo BalanceRepository
}

func NewBalanceService(repo BalanceRepository) *BalanceService {
	return &BalanceService{repo: repo}
}

func (s *BalanceService) GetUserBalance(ctx context.Context, login string) (domain.Balance, error) {
	return s.repo.GetBalance(ctx, login)
}

func (s *BalanceService) WithdrawBalance(ctx context.Context, login, order string, sum float64) error {
	balance, err := s.repo.GetBalance(ctx, login)
	if err != nil {
		return err
	}

	if sum > balance.Current {
		return appErrors.ErrInsufficientFunds
	}

	return s.repo.SaveWithdrawal(ctx, domain.Withdrawal{
		Login:       login,
		Order:       order,
		Sum:         sum,
		ProcessedAt: time.Now(),
	})
}

func (s *BalanceService) GetUserWithdrawals(ctx context.Context, login string) ([]domain.Withdrawal, error) {
	return s.repo.GetWithdrawalsByUser(ctx, login)
}
