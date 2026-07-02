package database

import (
	"fmt"
	"log"

	"github.com/refda/backend/internal/domain"
	"github.com/refda/backend/internal/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.OTPChallenge{},
		&domain.Event{},
		&domain.Gift{},
		&domain.Contribution{},
		&domain.Media{},
		&domain.WithdrawalRequest{},
	); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	log.Println("database connected and migrated")
	return db, nil
}
