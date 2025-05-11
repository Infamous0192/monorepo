package middleware

import (
	"app/pkg/exception"
	"app/pkg/types/http"
	"app/pkg/validation"

	"github.com/gofiber/fiber/v2"
)

// ErrorMiddleware provides custom error handling for the application
type ErrorMiddleware struct{}

// NewErrorMiddleware creates a new instance of ErrorMiddleware
func NewErrorMiddleware() *ErrorMiddleware {
	return &ErrorMiddleware{}
}

// Handler returns a Fiber error handler function that processes different types of errors
func (h *ErrorMiddleware) Handler() fiber.ErrorHandler {
	return func(ctx *fiber.Ctx, err error) error {
		switch err := err.(type) {
		case exception.HttpError:
			return ctx.Status(err.Code).JSON(http.ErrorResponse{
				Status:  err.Code,
				Message: err.Message,
				Errors:  err.Errors,
			})
		case validation.ValidationError:
			return ctx.Status(fiber.StatusUnprocessableEntity).JSON(http.ErrorResponse{
				Status:  fiber.StatusUnprocessableEntity,
				Message: "Invalid Payload",
				Errors:  err.Errors,
			})
		default:
			return ctx.Status(fiber.StatusInternalServerError).JSON(http.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: err.Error(),
				Errors:  nil,
			})
		}
	}
}
