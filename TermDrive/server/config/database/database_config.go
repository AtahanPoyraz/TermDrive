package database

import (
	"fmt"

	"github.com/AtahanPoyraz/TermDrive/server/config"
	"github.com/AtahanPoyraz/TermDrive/server/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// SetConfiguration establishes a connection to the PostgreSQL database using the configuration provided.
// It builds a Data Source Name (DSN) string from the provided database configuration,
// which includes the database host, user credentials, database name, port, SSL mode, and time zone.
// After opening the connection, the function performs automatic database migrations for the UserModel and FileModel.
//
// Arguments:
// - configuration: A pointer to the Configuration struct that contains database connection settings.
//
// Returns:
// - A pointer to the `gorm.DB` instance representing the database connection.
// - An error if the connection to the database or migration fails.
func SetConfiguration(configuration *config.Configuration) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		configuration.Database.Host,
		configuration.Database.User,
		configuration.Database.Password,
		configuration.Database.Name,
		configuration.Database.Port,
		configuration.Database.SSLMode,
		configuration.Database.TimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := db.AutoMigrate(&model.UserModel{}, &model.FileModel{}); err != nil {
		return nil, fmt.Errorf("failed to migrate database model: %v", err)
	}

	return db, nil
}
