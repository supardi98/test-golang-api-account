package routes

import (
	"supardi98/service-account-api/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

func SetupBalanceRoutes(app *fiber.App, db *gorm.DB) {
	log.Info("Setup BalanceRoutes...")
	balanceController := controllers.NewBalanceController(db)

	app.Post("/tabung", balanceController.PostTabung)
	app.Post("/tarik", balanceController.PostTarik)
	app.Get("/saldo/:no_rekening", balanceController.GetSaldoByNoRekening)
	app.Get("/mutasi/:no_rekening", balanceController.GetMutasiByNoRekening)

	log.Info("Berhasil BalanceRoutes")
}
