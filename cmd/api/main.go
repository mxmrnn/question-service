package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	app "question-service/internal/app"
	internalhttp "question-service/internal/http"
)

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	router := internalhttp.NewRouter()

	cfg := app.Config{Address: ":8080"}
	application := app.NewApp(logger, cfg, router)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := application.Run(ctx); err != nil {
		logger.Fatalf("application stopped with error: %v", err)
	}
}
