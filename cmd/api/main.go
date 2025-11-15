package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"question-service/internal/app"
	"question-service/internal/config"
	"question-service/internal/db"
	httptransport "question-service/internal/http"
	"question-service/internal/logger"
	"question-service/internal/repository"
	"question-service/internal/service"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	log := logger.New()
	defer log.Sync()

	cfg := config.Load()

	conn, err := db.New(cfg, log)
	if err != nil {
		return err
	}

	qRepo := repository.NewQuestionRepository(conn)
	aRepo := repository.NewAnswerRepository(conn)

	qSvc := service.NewQuestionService(qRepo)
	aSvc := service.NewAnswerService(aRepo, qRepo)

	router := httptransport.NewRouter(qSvc, aSvc)

	application := app.NewApp(log, app.Config{
		Address: cfg.HTTPPort,
	}, router, conn)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		log.Error("application stopped with error", zap.Error(err))
		return err
	}

	return nil
}
