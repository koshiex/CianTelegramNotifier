package database

import (
	"telegram_bot_service/internal/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Initialize(databasePath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto migrate schemas
	err = db.AutoMigrate(
		&models.User{},
		&models.Favorite{},
		&models.Subscription{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
