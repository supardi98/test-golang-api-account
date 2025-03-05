package controllers

import (
	"fmt"
	"supardi98/service-account-api/database"
	"supardi98/service-account-api/dto"
	"supardi98/service-account-api/models"
	"supardi98/service-account-api/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type BalanceController struct {
	DB *gorm.DB
}

func NewBalanceController(db *gorm.DB) *BalanceController {
	return &BalanceController{
		DB: db,
	}
}

// PostTabung godoc
// @Summary Menambah tabungan ke Rekening
// @Description Menambah tabungan ke Rekening
// @Tags Menambah Tabungan
// @Accept json
// @Produce json
// @Param body body dto.InputPostTabung true "Masukkan no_rekening, nominal (minimal 1)"
// @Success 200 {object} dto.ResponsePostTabung
// @Success 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tabung [post]
func (cc *BalanceController) PostTabung(c *fiber.Ctx) error {
	inputTabung := new(dto.InputPostTabung)

	if err := c.BodyParser(inputTabung); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	// Validasi input
	if err := utils.Validate(inputTabung); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, *err)
	}

	// Cek apakah No Rekening Benar
	var nasabah models.Nasabah
	if err := database.DB.Where(&models.Nasabah{NoRekening: inputTabung.NoRekening}).First(&nasabah).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusBadRequest, "No Rekening tidak dikenali.")
		}
		log.Error(fmt.Sprintf("Database Error [POST /tabung | Mengecek No Rekening %s]:", inputTabung.NoRekening), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Hitung saldo akhir
	saldoAkhir := nasabah.Saldo + inputTabung.Nominal

	// Buat mutasi tabung
	mutasi := models.Mutasi{
		NasabahID:     nasabah.ID,
		NoRekening:    nasabah.NoRekening,
		SaldoAwal:     nasabah.Saldo,
		SaldoAkhir:    saldoAkhir,
		NominalTabung: &inputTabung.Nominal,
		TipeTransaksi: models.Tabung,
		CreatedAt:     time.Now(),
	}

	// Simpan mutasi dan perbarui saldo nasabah dalam transaksi
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&mutasi).Error; err != nil {
			return err
		}

		if err := tx.Model(&nasabah).Update("saldo", saldoAkhir).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Error(fmt.Sprintf("Database Error [POST /tabung | Gagal menyimpan transaksi %s]:", inputTabung.NoRekening), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan transaksi")
	}

	return c.Status(fiber.StatusOK).JSON(dto.ResponsePostTabung{
		Saldo: nasabah.Saldo,
	})
}

// PostTarik godoc
// @Summary Menarik saldo dari Rekening
// @Description Menarik saldo dari Rekening
// @Tags Menarik Saldo
// @Accept json
// @Produce json
// @Param body body dto.InputPostTarik true "Masukkan no_rekening, nominal tarik (minimal 1)"
// @Success 200 {object} dto.ResponsePostTarik
// @Success 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tarik [post]
func (r *BalanceController) PostTarik(c *fiber.Ctx) error {
	inputTarik := new(dto.InputPostTarik)

	if err := c.BodyParser(inputTarik); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	// Validasi input
	validationErrors := utils.Validate(inputTarik)
	if validationErrors != nil {
		return fiber.NewError(fiber.StatusBadRequest, *validationErrors)
	}

	// Cek apakah No Rekening Benar
	var nasabah models.Nasabah
	if err := database.DB.Where(&models.Nasabah{NoRekening: inputTarik.NoRekening}).First(&nasabah).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusBadRequest, "No Rekening tidak dikenali.")
		}
		log.Error(fmt.Sprintf("Database Error [POST /tarik | Mengecek No Rekening %s]:", inputTarik.NoRekening), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Validasi apakah saldo cukup
	if nasabah.Saldo < inputTarik.Nominal {
		return fiber.NewError(fiber.StatusBadRequest, "Saldo tidak cukup.")
	}

	// Hitung saldo akhir
	saldoAkhir := nasabah.Saldo - inputTarik.Nominal

	// Buat mutasi tarik
	mutasi := models.Mutasi{
		NasabahID:     nasabah.ID,
		NoRekening:    nasabah.NoRekening,
		SaldoAwal:     nasabah.Saldo,
		SaldoAkhir:    saldoAkhir,
		NominalTarik:  &inputTarik.Nominal,
		TipeTransaksi: models.Tarik,
		CreatedAt:     time.Now(),
	}

	// Simpan mutasi dan perbarui saldo nasabah dalam transaksi
	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&mutasi).Error; err != nil {
			return err
		}

		if err := tx.Model(&nasabah).Update("saldo", saldoAkhir).Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		log.Error(fmt.Sprintf("Database Error [POST /tarik | Gagal menyimpan transaksi %s]:", inputTarik.NoRekening), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal menyimpan transaksi")
	}

	return c.Status(fiber.StatusOK).JSON(dto.ResponsePostTarik{
		Saldo: nasabah.Saldo,
	})
}

// GetSaldoByNoRekening godoc
// @Summary Melihat saldo dari Rekening
// @Description Melihat saldo dari Rekening
// @Tags Melihat Saldo
// @Produce json
// @Param no_rekening path string true "Nomor Rekening"
// @Success 200 {object} dto.ResponseGetSaldoByNoRekening
// @Success 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /saldo/{no_rekening} [get]
func (r *BalanceController) GetSaldoByNoRekening(c *fiber.Ctx) error {
	noRekening := c.Params("no_rekening")

	// Cek apakah No Rekening Benar
	var nasabah models.Nasabah
	if err := database.DB.Where(&models.Nasabah{NoRekening: noRekening}).First(&nasabah).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusBadRequest, "No Rekening tidak dikenali.")
		}
		log.Error(fmt.Sprintf("Database Error [POST /saldo/%s | Mengecek No Rekening %s]:", noRekening, noRekening), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Kembalikan saldo nasabah
	return c.Status(fiber.StatusOK).JSON(dto.ResponseGetSaldoByNoRekening{
		Saldo: nasabah.Saldo,
	})
}

// GetMutasiByNoRekening godoc
// @Summary Melihat mutasi dari Rekening
// @Description Melihat mutasi dari Rekening
// @Tags Melihat Mutasi
// @Produce json
// @Param no_rekening path string true "Nomor Rekening"
// @Success 200 {object} dto.ResponseGetMutasiByNoRekening
// @Success 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /mutasi/{no_rekening} [get]
func (r *BalanceController) GetMutasiByNoRekening(c *fiber.Ctx) error {
	noRekening := c.Params("no_rekening")

	// Cek apakah No Rekening Benar
	var nasabah models.Nasabah
	if err := database.DB.Where(&models.Nasabah{NoRekening: noRekening}).First(&nasabah).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fiber.NewError(fiber.StatusBadRequest, "No Rekening tidak dikenali.")
		}
		log.Error(fmt.Sprintf("Database Error [POST /mutasi/%s | Mengecek No Rekening %s]:", noRekening, noRekening), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	var mutasi []models.Mutasi
	if err := database.DB.Where(&models.Mutasi{NasabahID: nasabah.ID}).Find(&mutasi).Error; err != nil {
		log.Error(fmt.Sprintf("Database Error [POST /mutasi/%s | Gagal Mengambil Data Mutasi dari No Rekening %s]:", noRekening, noRekening), err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal mengambil data mutasi")
	}

	// Kembalikan saldo nasabah
	return c.Status(fiber.StatusOK).JSON(dto.ResponseGetMutasiByNoRekening(mutasi))
}
