package database

import (
	"fmt"
	"supardi98/service-account-api/config"
	"supardi98/service-account-api/models"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {

	// Connect to database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort)

	log.Info("Mencoba menghubungkan ke database...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	log.Info("Berhasil terhubung ke database")

	// Migrate the schema
	db.AutoMigrate(
		&models.Nasabah{},
		&models.Mutasi{})

	DB = db

	return db
}
