package app

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
)

type HTTPError struct {
	Status           int                 `json:"status,omitempty"`
	Error            string              `json:"error,omitempty"`
	Code             string              `json:"code,omitempty"`
	ValidationErrors []map[string]string `json:"validation_errors,omitempty"`
	RequestID        string              `json:"request_id,omitempty"`
}

func ErrorResponse(c *fiber.Ctx, err error) error {
	resp := HTTPError{
		Status:    fiber.StatusInternalServerError,
		Error:     err.Error(),
		Code:      internal,
		RequestID: c.Locals("requestid").(string),
	}

	var e *fiber.Error
	if errors.As(err, &e) {
		resp.Status = e.Code

		if e.Code == fiber.StatusNotFound {
			resp.Code = notFound
		}
	}

	if e, ok := err.(interface {
		GetStatus() int
	}); ok {
		resp.Status = e.GetStatus()
	}

	if e, ok := err.(interface {
		GetCode() string
	}); ok {
		resp.Code = e.GetCode()
	}

	var valErrors validator.ValidationErrors
	if errors.As(err, &valErrors) {
		resp.Code = "validation"
		resp.Status = fiber.StatusBadRequest
		resp.ValidationErrors = lo.Map(valErrors, func(item validator.FieldError, _ int) map[string]string {
			return map[string]string{
				"field": item.Field(),
				"error": item.Error(),
			}
		})
		resp.Error = "Validation Error"
	}

	return c.Status(resp.Status).JSON(resp)
}
