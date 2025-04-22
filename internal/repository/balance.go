package repository

import (
	"context"

	"github.com/Te8va/Gofermarch/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BalanceRepository struct {
	db *pgxpool.Pool
}

func NewFinanceRepository(db *pgxpool.Pool) *BalanceRepository {
	return &BalanceRepository{db: db}
}

func (r *BalanceRepository) GetBalance(ctx context.Context, login string) (domain.Balance, error) {
	var accrued, withdrawn float64
	err := r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(accrual), 0)
		FROM orders
		WHERE login = $1`, login).Scan(&accrued)
	if err != nil {
		return domain.Balance{}, err
	}

	err = r.db.QueryRow(ctx, `
		SELECT COALESCE(SUM(withdrawn), 0)
		FROM withdrawals
		WHERE login = $1`, login).Scan(&withdrawn)
	if err != nil {
		return domain.Balance{}, err
	}

	return domain.Balance{
		Current:   accrued - withdrawn,
		Withdrawn: withdrawn,
	}, nil
}

func (r *BalanceRepository) SaveWithdrawal(ctx context.Context, w domain.Withdrawal) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO withdrawals (login, order_number, withdrawn, processed_at)
		 VALUES ($1, $2, $3, $4)`,
		w.Login, w.Order, w.Sum, w.ProcessedAt)
	return err
}

func (r *BalanceRepository) GetWithdrawalsByUser(ctx context.Context, login string) ([]domain.Withdrawal, error) {
	rows, err := r.db.Query(ctx,
		`SELECT order_number, withdrawn, processed_at
		 FROM withdrawals
		 WHERE login = $1
		 ORDER BY processed_at DESC`,
		login)
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
