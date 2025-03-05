package utils

import (
	"math/rand"
	"strconv"
	"supardi98/service-account-api/database"
	"supardi98/service-account-api/models"
	"time"

	"gorm.io/gorm"
)

func GenerateRandomRekening() (*string, error) {
	source := rand.NewSource(time.Now().UnixNano())
	random := rand.New(source)
	randomNumber := random.Intn(9000000000) + 1000000000 // random angka 8 digit

	noRekening := strconv.Itoa(randomNumber)

	var existingNoRekening models.Nasabah
	if err := database.DB.Where(&models.Nasabah{NoRekening: noRekening}).First(&existingNoRekening).Error; err == nil {
		return GenerateRandomRekening()
	} else if err != gorm.ErrRecordNotFound {
		// Database Error
		return nil, err
	}

	return &noRekening, nil
}
