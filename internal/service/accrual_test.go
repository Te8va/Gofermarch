package service

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/Te8va/Gofermarch/internal/service/mocks"
)

func TestOrderService_StartAccrualProcessing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	orderNumber := "123456"
	orderStatus := "PROCESSED"
	accrualAmount := 42.5

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"order":   orderNumber,
			"status":  orderStatus,
			"accrual": accrualAmount,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	mockRepo := mocks.NewMockOrderServ(ctrl)

	var wg sync.WaitGroup
	wg.Add(1)
	mockRepo.EXPECT().
		UpdateOrder(gomock.Any(), orderNumber, orderStatus, accrualAmount).
		DoAndReturn(func(ctx context.Context, number, status string, accrual float64) error {
			defer wg.Done()
			require.Equal(t, orderNumber, number)
			require.Equal(t, orderStatus, status)
			require.Equal(t, accrualAmount, accrual)
			return nil
		})

	service := &OrderService{
		srv:              mockRepo,
		accrualSystemURL: server.URL,
	}

	ctx := context.Background()
	service.StartAccrualProcessing(ctx, []string{orderNumber})

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("timeout waiting for order update")
	}
}
