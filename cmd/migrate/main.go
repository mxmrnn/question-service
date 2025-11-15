package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"

	"question-service/internal/config"
)

const migrationsDir = "migrations"

func main() {
	cfg := config.Load()

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBSSL,
	)

	log.Println("connecting to database for migrations:", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping db: %v", err)
	}

	goose.SetLogger(log.New(os.Stdout, "[goose] ", log.LstdFlags))

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatalf("failed to set goose dialect: %v", err)
	}

	log.Println("running migrations from", migrationsDir)

	if err := goose.Up(db, migrationsDir); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	current, err := goose.EnsureDBVersion(db)
	if err != nil {
		log.Printf("migrations applied, current version: unknown")
	} else {
		log.Printf("migrations already up to date, current version: %d", current)
	}

	log.Println("migrations applied successfully")
}
