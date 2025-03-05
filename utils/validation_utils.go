package utils

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validate(data interface{}) *string {
	errs := validate.Struct(data)
	if errs != nil {
		validationMsg := ""
		validationError := []string{}

		for _, err := range errs.(validator.ValidationErrors) {
			errorMsg := ""
			errorTag := err.Tag()
			errorParam := err.Param()

			if errorTag == "required" {
				errorMsg = "wajib diisi"
			} else if errorTag == "number" {
				errorMsg = "isi dengan angka"
			} else if errorTag == "len" {
				errorMsg = "harus berjumlah " + errorParam + " digit"
			} else if errorTag == "min" {
				errorMsg = "harus minimal " + errorParam + " digit"
			} else if errorTag == "max" {
				errorMsg = "harus maksimal " + errorParam + " digit"
			} else {
				errorMsg = err.Tag()

				if errorParam != "" {
					errorMsg += " " + errorParam + "."
				}
			}

			validationError = append(validationError, err.Field()+" "+errorMsg)
		}

		validationMsg += strings.Join(validationError, ", ")

		return &validationMsg
	}

	return nil
}
