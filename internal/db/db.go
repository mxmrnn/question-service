package db

import (
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"question-service/internal/config"
	"question-service/internal/logger"
)

const maxAttempts = 5

// New создаёт новое подключение к PostgreSQL с помощью GORM.
func New(cfg *config.Config, log *logger.Logger) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=UTC",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBName,
		cfg.DBSSL,
	)

	var (
		db  *gorm.DB
		err error
	)

	log.Info("connecting to database",
		zap.String("host", cfg.DBHost),
		zap.String("port", cfg.DBPort),
		zap.String("db", cfg.DBName),
	)

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Info("connecting to database",
			zap.String("host", cfg.DBHost),
			zap.String("port", cfg.DBPort),
			zap.String("db", cfg.DBName),
			zap.Int("attempt", attempt),
		)

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		log.Warn("failed to connect to database, will retry",
			zap.Error(err),
			zap.Int("attempt", attempt),
		)

		time.Sleep(time.Duration(attempt) * time.Second)
	}

	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// базовые настройки пула соединений
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Info("database connection established")

	return db, nil
}
