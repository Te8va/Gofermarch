package router

import (
	"net/http"

	"github.com/Te8va/Gofermarch/internal/handler"
	"github.com/Te8va/Gofermarch/internal/middleware"
)

func NewRouter(authHandler *handler.AuthorizationHandler, orderHandler *handler.OrderHandler, balanceHandler *handler.BalanceHandler, authMiddleware func(http.Handler) http.Handler) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /api/user/register", middleware.Log(http.HandlerFunc(authHandler.RegisterHandler)))
	mux.Handle("POST /api/user/login", middleware.Log(http.HandlerFunc(authHandler.LoginHandler)))

	mux.Handle("POST /api/user/orders", middleware.Log(authMiddleware(http.HandlerFunc(orderHandler.UploadOrder))))
	mux.Handle("GET /api/user/orders", middleware.Log(authMiddleware(http.HandlerFunc(orderHandler.GetOrders))))

	mux.Handle("GET /api/user/balance", middleware.Log(authMiddleware(http.HandlerFunc(balanceHandler.GetUserBalance))))
	mux.Handle("POST /api/user/balance/withdraw", middleware.Log(authMiddleware(http.HandlerFunc(balanceHandler.WithdrawBalance))))
	mux.Handle("GET /api/user/withdrawals", middleware.Log(authMiddleware(http.HandlerFunc(balanceHandler.GetUserWithdrawals))))

	return mux
}
