package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var ErrBodyParser = Errorf("error while parsing body")

func BindAndValidate(c *fiber.Ctx, v interface{}) error {
	if err := c.BodyParser(v); err != nil {
		return ErrBodyParser
	}

	if err := validator.New().Struct(v); err != nil {
		return err
	}

	return nil
}
