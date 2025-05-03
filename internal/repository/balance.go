package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Te8va/Gofermarch/internal/domain"
)

type BalanceRepository struct {
	db *pgxpool.Pool
}

func NewFinanceRepository(db *pgxpool.Pool) *BalanceRepository {
	return &BalanceRepository{db: db}
}

func (r *BalanceRepository) GetBalance(ctx context.Context, login string) (domain.Balance, error) {
	var accrued, withdrawn float64
	err := r.db.QueryRow(ctx,queryGetAccrued, login).Scan(&accrued)
	if err != nil {
		return domain.Balance{}, err
	}

	err = r.db.QueryRow(ctx,queryGetWithdrawn, login).Scan(&withdrawn)
	if err != nil {
		return domain.Balance{}, err
	}

	return domain.Balance{
		Current:   accrued - withdrawn,
		Withdrawn: withdrawn,
	}, nil
}

func (r *BalanceRepository) SaveWithdrawal(ctx context.Context, w domain.Withdrawal) error {
	_, err := r.db.Exec(ctx,queryInsertWithdrawal, w.Login, w.Order, w.Sum, w.ProcessedAt)
	return err
}

func (r *BalanceRepository) GetWithdrawalsByUser(ctx context.Context, login string) ([]domain.Withdrawal, error) {
	rows, err := r.db.Query(ctx, queryGetWithdrawalsByUser,login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []domain.Withdrawal
	for rows.Next() {
		var w domain.Withdrawal
		err := rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt)
		if err != nil {
			return nil, err
		}
		w.Login = login
		withdrawals = append(withdrawals, w)
	}
	return withdrawals, nil
}
