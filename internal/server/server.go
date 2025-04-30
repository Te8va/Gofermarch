package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"

	"github.com/Te8va/Gofermarch/internal/config"
	"github.com/Te8va/Gofermarch/internal/handler"
	"github.com/Te8va/Gofermarch/internal/middleware"
	"github.com/Te8va/Gofermarch/internal/repository"
	"github.com/Te8va/Gofermarch/internal/router"
	"github.com/Te8va/Gofermarch/internal/service"
	"github.com/Te8va/Gofermarch/pkg/logger"
)

func Serve() error {
	cfg := config.NewConfig()

	m, err := migrate.New("file://migrations", cfg.DatabaseURI)
	if err != nil {
		logger.Logger().Fatalln(zap.Error(err))
	}
	if err := repository.ApplyMigrations(m); err != nil {
		logger.Logger().Fatalln(zap.Error(err))
	}
	logger.Logger().Infoln("Migrations applied successfully")

	pool, err := repository.GetPgxPool(cfg.DatabaseURI)
	if err != nil {
		logger.Logger().Fatalln(zap.Error(err))
	}
	logger.Logger().Infoln("Postgres connection pool created")

	authRepo := repository.NewAuthorizationRepository(pool)
	authService := service.NewAuthorization(authRepo, cfg.JWTKey)
	authHandler := handler.NewAuthorizationHandler(authService)

	orderRepo := repository.NewOrderRepository(pool)
	orderService := service.NewOrderService(orderRepo, cfg.AccrualSystemAddress)
	orderHandler := handler.NewOrderHandler(orderService)

	balanceRepo := repository.NewFinanceRepository(pool)
	balanceService := service.NewBalanceService(balanceRepo)
	balanceHandler := handler.NewBalanceHandler(balanceService)

	authMiddleware := middleware.Auth(cfg.JWTKey)

	mux := router.NewRouter(authHandler, orderHandler, balanceHandler, authMiddleware)

	server := &http.Server{
		Addr:     cfg.RunAddress,
		ErrorLog: log.New(logger.Logger(), "", 0),
		Handler:  mux,
	}

	var wg sync.WaitGroup

	go func() {
		logger.Logger().Infoln("Server started, listening on port", cfg.RunAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Logger().Fatalln("ListenAndServe failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	logger.Logger().Infoln("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Logger().Fatalln("Server was forced to shutdown:", zap.Error(err))
	}

	waitGroupChan := make(chan struct{})
	go func() {
		wg.Wait()
		waitGroupChan <- struct{}{}
	}()

	select {
	case <-waitGroupChan:
		logger.Logger().Infoln("All delete goroutines successfully finished")
	case <-time.After(time.Second * 3):
		logger.Logger().Infoln("Some of delete goroutines have not completed their job due to shutdown timeout")
	}

	logger.Logger().Infoln("Server was shut down")
	return nil
}
