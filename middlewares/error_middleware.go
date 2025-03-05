package middlewares

import (
	"supardi98/service-account-api/dto"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500 (Internal Server Error)
	code := fiber.StatusInternalServerError
	if e, ok := err.(*fiber.Error); ok {
		// Override status code if it's a fiber.Error
		code = e.Code
	}

	return c.Status(code).JSON(dto.ErrorResponse{
		Remark: err.Error(),
	})
}
