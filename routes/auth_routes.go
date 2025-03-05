package routes

import (
	"supardi98/service-account-api/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func SetupAuthRoutes(app *fiber.App, db *gorm.DB) {
	log.Info("Setup AuthRoutes...")
	authController := controllers.NewAuthController(db)

	app.Post("/daftar", authController.PostDaftar)

	log.Info("Berhasil Setup AuthRoutes")
}
