package service

import (
	"context"
	"errors"
	"time"

	"github.com/Te8va/Gofermarch/internal/domain"
	appErrors "github.com/Te8va/Gofermarch/internal/errors"
)

//go:generate mockgen -source=order.go -destination=mocks/mock_order.go -package=mocks

type OrderServ interface {
	GetOrder(ctx context.Context, number string) (domain.Order, error)
	SaveOrder(ctx context.Context, order domain.Order) error
	GetOrdersByUser(ctx context.Context, login string) ([]domain.Order, error)
	UpdateOrder(ctx context.Context, number, status string, accrual float64) error
}

type OrderService struct {
	srv              OrderServ
	accrualSystemURL string
}

func NewOrderService(srv OrderServ, accrualSystemURL string) *OrderService {
	return &OrderService{
		srv:              srv,
		accrualSystemURL: accrualSystemURL,
	}
}

func (s *OrderService) ProcessOrder(ctx context.Context, number, login string) (domain.OrderStatus, error) {
	order, err := s.srv.GetOrder(ctx, number)
	if err == nil {
		if order.Login == login {
			return domain.StatusAlreadyUploaded, errors.New("order already uploaded by the same user")
		}
		return domain.StatusConflict, appErrors.ErrOrderExists
	}

	err = s.srv.SaveOrder(ctx, domain.Order{
		Number:     number,
		Login:      login,
		UploadedAt: time.Now(),
		Status:     string(domain.StatusNew),
	})
	if err != nil {
		return domain.StatusInternal, err
	}

	go s.StartAccrualProcessing(context.Background(), []string{number})

	return domain.StatusNew, nil
}

func (s *OrderService) GetOrdersByUser(ctx context.Context, login string) ([]domain.Order, error) {
	return s.srv.GetOrdersByUser(ctx, login)
}
