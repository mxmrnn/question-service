package app

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"net/http"
	"time"

	"question-service/internal/logger"
)

// Config описывает конфигурацию HTTP-приложения.
type Config struct {
	Address string
}

// App представляет собой HTTP-приложение.
type App struct {
	Logger     *logger.Logger
	Router     http.Handler
	HTTPServer *http.Server
}

// NewApp создаёт новый экземпляр App на основе переданных зависимостей и конфигурации.
func NewApp(logger *logger.Logger, cfg Config, router http.Handler) *App {
	if cfg.Address == "" {
		cfg.Address = ":8080"
	}

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	return &App{
		Logger:     logger,
		Router:     router,
		HTTPServer: server,
	}
}

// Run запускает HTTP-сервер и блокируется до отмены контекста или ошибки сервера.
func (a *App) Run(ctx context.Context) error {
	serverErr := make(chan error, 1)

	go func() {
		a.Logger.Info("starting HTTP server",
			zap.String("address", a.HTTPServer.Addr),
		)

		if err := a.HTTPServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	// ждём либо сигнала остановки (ctx.Done), либо ошибки сервера
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		a.Logger.Info("shutting down HTTP server")

		if err := a.HTTPServer.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
			return err
		}

		return nil

	case err := <-serverErr:
		return err
	}
}
