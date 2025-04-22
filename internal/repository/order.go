package repository

import (
	"context"

	"github.com/Te8va/Gofermarch/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetOrder(ctx context.Context, number string) (domain.Order, error) {
	var order domain.Order
	err := r.db.QueryRow(ctx,
		`SELECT number, login, uploaded_at, status, accrual
		 FROM orders
		 WHERE number = $1`,
		number).Scan(&order.Number, &order.Login, &order.UploadedAt, &order.Status, &order.Accrual)

	if err != nil {
		return domain.Order{}, err
	}

	return order, nil
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order domain.Order) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO orders (number, login, uploaded_at, status)
		 VALUES ($1, $2, $3, $4)`,
		order.Number, order.Login, order.UploadedAt, order.Status)
	return err
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, number, status string, accrual float64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE orders
		 SET status = $2,
		     accrual = $3
		 WHERE number = $1`,
		number, status, accrual)
	return err
}

func (r *OrderRepository) GetOrdersByUser(ctx context.Context, login string) ([]domain.Order, error) {
	rows, err := r.db.Query(ctx,
		`SELECT number, status, accrual, uploaded_at
		 FROM orders
		 WHERE login = $1
		 ORDER BY uploaded_at DESC`,
		login)
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