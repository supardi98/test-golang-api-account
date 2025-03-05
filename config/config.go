package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

var (
	Host string
	Port string

	APP_URL string

	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
)

func LoadConfig() {
	log.Info("Memuat environment variabel")

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Warn("File environment (.env) tidak ditemukan")
		log.Info("Menggunakan variabel lingkungan")
	} else {
		log.Info("Berhasil mengambil file environment (.env)")
	}

	flag.StringVar(&Host, "host", "0.0.0.0", "REST API Host")
	flag.StringVar(&Port, "port", "8080", "REST API Port")
	flag.Parse()

	APP_URL = os.Getenv("APP_URL")
	if APP_URL == "" {
		if Host == "0.0.0.0" {
			APP_URL = fmt.Sprintf("%s:%s", "localhost", Port)
		} else {
			APP_URL = fmt.Sprintf("%s:%s", Host, Port)

		}
	}

	DBHost = os.Getenv("DB_HOST")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")
	DBPort = os.Getenv("DB_PORT")

	// Log informasi database
	if DBHost != "" && DBUser != "" && DBName != "" && DBPort != "" {
		log.Info("Database configuration berhasil dimuat.")
	} else {
		log.Fatal("Database configuration gagal dimuat.")
	}
}
