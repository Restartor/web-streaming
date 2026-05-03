package config

import (
	"backend/internal/domain"
	"backend/pkg/logger"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseConnection() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Log.Fatal().Err(err).Msg("gagal koneksi database!")
	}

	err = db.AutoMigrate(
		&domain.User{},
		&domain.Filem{},
	)

	if err != nil {
		logger.Log.Fatal().Err(err).Msg("gagal migrasi domain!")
	}

	DB = db

	logger.Log.Info().Msg("koneksi database success!, silahkan lanjut")

}
