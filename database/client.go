package database

import (
	"GoShort/internal/models"
	"GoShort/settings"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var (
	instance *gorm.DB
	dbError  error
)

func Connect() {
	instance, dbError = gorm.Open(postgres.Open(settings.Settings.Database.Dsn), &gorm.Config{})
	if dbError != nil {
		log.Fatal(dbError)
	}
}

func MigrateModels() {
	err := instance.AutoMigrate(
		&models.User{},
		&models.Session{},
		&models.Shortlink{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database schema:", err)
	}
}
