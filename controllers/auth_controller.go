package controllers

import (
	"supardi98/service-account-api/dto"
	"supardi98/service-account-api/models"
	"supardi98/service-account-api/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type AuthController struct {
	DB *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{
		DB: db,
	}
}

// PostDaftar godoc
// @Summary Membuat akun baru nasabah
// @Description Membuat akun baru nasabah
// @Tags Daftar Nasabah
// @Accept json
// @Produce json
// @Param body body dto.InputPostDaftar true "Masukkan nama, nik (16 Digit, Angka) dan no_hp (9-14 Digit, Angka)"
// @Success 200 {object} dto.ResponsePostDaftar
// @Success 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /daftar [post]
func (cc *AuthController) PostDaftar(c *fiber.Ctx) error {
	inputDaftar := new(dto.InputPostDaftar)

	if err := c.BodyParser(inputDaftar); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Cannot parse JSON")
	}

	// Validasi input
	if err := utils.Validate(inputDaftar); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, *err)
	}

	// Cek apakah NIK sudah digunakan
	var existingNIK models.Nasabah
	if err := cc.DB.Where(&models.Nasabah{NIK: inputDaftar.NIK}).First(&existingNIK).Error; err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "NIK sudah digunakan.")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("Database Error [POST /daftar | Mengecek NIK]:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	// Cek apakah NIK sudah digunakan
	var existingNoHp models.Nasabah
	if err := cc.DB.Where(&models.Nasabah{NoHP: inputDaftar.NoHP}).First(&existingNoHp).Error; err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "No HP sudah digunakan.")
	} else if err != gorm.ErrRecordNotFound {
		log.Error("Database Error [POST /daftar | Mengecek No HP]:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	noRekening, err := utils.GenerateRandomRekening()
	if err != nil {
		log.Error("Database Error [POST /daftar | Generate Random No Rekening]:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Database error")
	}

	nasabah := models.Nasabah{
		Nama:       inputDaftar.Nama,
		NIK:        inputDaftar.NIK,
		NoHP:       inputDaftar.NoHP,
		NoRekening: *noRekening,
	}

	if err := cc.DB.Create(&nasabah).Error; err != nil {
		log.Error("Database Error [POST /daftar | Membuat Nasabah baru pada table]:", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Gagal Membuat Nasabah")
	}

	return c.Status(fiber.StatusCreated).JSON(dto.ResponsePostDaftar{
		NoRekening: nasabah.NoRekening,
	})
}
