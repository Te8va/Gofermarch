package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Te8va/Gofermarch/internal/domain"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetOrder(ctx context.Context, number string) (domain.Order, error) {
	var order domain.Order
	err := r.db.QueryRow(ctx, queryGetOrder, number).Scan(&order.Number, &order.Login, &order.UploadedAt, &order.Status, &order.Accrual)

	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order domain.Order) error {
	_, err := r.db.Exec(ctx, queryInsertOrder, order.Number, order.Login, order.UploadedAt, order.Status)
	return err
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, number, status string, accrual float64) error {
	_, err := r.db.Exec(ctx, queryUpdateOrder, number, status, accrual)
	return err
}

func (r *OrderRepository) GetOrdersByUser(ctx context.Context, login string) ([]domain.Order, error) {
	rows, err := r.db.Query(ctx, queryGetOrdersByUser, login)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var o domain.Order
		err := rows.Scan(&o.Number, &o.Status, &o.Accrual, &o.UploadedAt)
		if err != nil {
			return nil, err
		}
		o.Login = login
		orders = append(orders, o)
	}
	return orders, nil
}
