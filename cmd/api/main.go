package main

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"

	"question-service/internal/app"
	"question-service/internal/config"
	httptransport "question-service/internal/http"
	"question-service/internal/logger"
)

func main() {
	log := logger.New()
	defer log.Sync()
	cfg := config.Load()

	router := httptransport.NewRouter()

	application := app.NewApp(log, app.Config{
		Address: cfg.HTTPPort,
	}, router)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		log.Fatal("application stopped with error", zap.Error(err))
	}
}
