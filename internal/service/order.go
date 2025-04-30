package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
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

	go func() {
		if err := s.updateOrderFromAccrualSystem(context.Background(), number); err != nil {
			log.Printf("updateOrderFromAccrualSystem error: %v", err)
		}
	}()

	return domain.StatusNew, nil
}

func (s *OrderService) updateOrderFromAccrualSystem(ctx context.Context, number string) error {
	url := fmt.Sprintf("%s/api/orders/%s", strings.TrimRight(s.accrualSystemURL, "/"), number)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var ext struct {
			Order   string  `json:"order"`
			Status  string  `json:"status"`
			Accrual float64 `json:"accrual"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&ext); err != nil {
			return err
		}

		var internalStatus string
		switch ext.Status {
		case "REGISTERED", "PROCESSING":
			internalStatus = "PROCESSING"
		case "INVALID":
			internalStatus = "INVALID"
		case "PROCESSED":
			internalStatus = "PROCESSED"
		default:
			internalStatus = "NEW"
		}

		return s.srv.UpdateOrder(ctx, number, internalStatus, ext.Accrual)

	case http.StatusTooManyRequests:
		return appErrors.ErrTooManyRequests
	}
	return nil
}

func (s *OrderService) GetOrdersByUser(ctx context.Context, login string) ([]domain.Order, error) {
	return s.srv.GetOrdersByUser(ctx, login)
}
