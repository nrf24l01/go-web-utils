package pg_kit

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nrf24l01/go-web-utils/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func RegisterPostgres(cfg *config.PGConfig, noSlowSQLLogs bool, models ...interface{}) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.PGHost, cfg.PGUser, cfg.PGPassword, cfg.PGDatabase, cfg.PGPort, cfg.PGSSLMode, cfg.PGTimeZone,
	)

	// Настройка логгера GORM: отключаем лог медленных запросов, если передан флаг
	slowThreshold := 200 * time.Millisecond
	if noSlowSQLLogs {
		slowThreshold = 0
	}
	gormLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: slowThreshold,
			LogLevel:      logger.Warn,
			Colorful:      false,
		},
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get db instance: %w", err)
	}

	// Создание расширения pgcrypto
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto";`).Error; err != nil {
		return nil, fmt.Errorf("failed to create extension: %w", err)
	}

	// Автоматическая миграция переданных моделей
	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			return nil, fmt.Errorf("failed to migrate database: %w", err)
		}
		log.Printf("Database migrated successfully for models: %v", models)
	}

	return db, nil
}
