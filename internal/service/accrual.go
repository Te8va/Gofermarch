package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	appErrors "github.com/Te8va/Gofermarch/internal/errors"
)

//go:generate mockgen -source=accrual.go -destination=mocks/mock_accrual.go -package=mocks

type AccrualClient interface {
	GetOrderStatus(orderNumber string) (status string, accrual float64, err error)
}

type accrualJob struct {
	OrderNumber string
}

type accrualResult struct {
	Number  string
	Status  string
	Accrual float64
	Err     error
}

func (s *OrderService) StartAccrualProcessing(ctx context.Context, orders []string) {
	const numWorkers = 5

	jobs := make(chan accrualJob, len(orders))
	results := make(chan accrualResult, len(orders))

	for i := 0; i < numWorkers; i++ {
		go s.accrualWorker(ctx, jobs, results)
	}

	for _, number := range orders {
		jobs <- accrualJob{OrderNumber: number}
	}
	close(jobs)

	go func() {
		for i := 0; i < len(orders); i++ {
			res := <-results
			if res.Err != nil {
				log.Printf("Failed to process order %s: %v", res.Number, res.Err)
				continue
			}

			internalStatus := map[string]string{
				"REGISTERED": "PROCESSING",
				"PROCESSING": "PROCESSING",
				"INVALID":    "INVALID",
				"PROCESSED":  "PROCESSED",
			}[res.Status]
			if internalStatus == "" {
				internalStatus = "NEW"
			}

			if err := s.srv.UpdateOrder(ctx, res.Number, internalStatus, res.Accrual); err != nil {
				log.Printf("Failed to update order %s: %v", res.Number, err)
			}
		}
	}()
}

func (s *OrderService) accrualWorker(ctx context.Context, jobs <-chan accrualJob, results chan<- accrualResult) {
	for job := range jobs {
		url := fmt.Sprintf("%s/api/orders/%s", strings.TrimRight(s.accrualSystemURL, "/"), job.OrderNumber)
		resp, err := http.Get(url)
		if err != nil {
			results <- accrualResult{Number: job.OrderNumber, Err: err}
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			results <- accrualResult{Number: job.OrderNumber, Err: appErrors.ErrTooManyRequests}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			results <- accrualResult{Number: job.OrderNumber, Err: fmt.Errorf("unexpected status: %d", resp.StatusCode)}
			continue
		}

		var ext struct {
			Order   string  `json:"order"`
			Status  string  `json:"status"`
			Accrual float64 `json:"accrual"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&ext); err != nil {
			results <- accrualResult{Number: job.OrderNumber, Err: err}
			continue
		}

		results <- accrualResult{
			Number:  ext.Order,
			Status:  ext.Status,
			Accrual: ext.Accrual,
		}
	}
}
