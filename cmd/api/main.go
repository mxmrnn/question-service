package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"question-service/internal/app"
	"question-service/internal/config"
	"question-service/internal/db"
	httptransport "question-service/internal/http"
	"question-service/internal/logger"
)

func main() {
	log := logger.New()
	defer log.Sync()
	cfg := config.Load()

	router := httptransport.NewRouter()

	conn, err := db.New(cfg, log)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	application := app.NewApp(log, app.Config{
		Address: cfg.HTTPPort,
	}, router, conn)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		log.Fatal("application stopped with error", zap.Error(err))
	}
}
