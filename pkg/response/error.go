package response

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ErrorResponse struct {
	Code             int    `json:"code"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewError(code int, error string, description string) ErrorResponse {
	return ErrorResponse{
		Code:             code,
		Error:            error,
		ErrorDescription: description,
	}
}

func JSON(c *gin.Context, code int, error string, description string) {
	c.JSON(code, NewError(code, error, description))
}

func FormatValidationError(err error) string {
	var messages []string

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := toSnakeCase(e.Field())

			switch e.Tag() {
			case "required":
				messages = append(messages, fmt.Sprintf("%s is required", field))
			case "uuid":
				messages = append(
					messages,
					fmt.Sprintf("%s must be a valid UUID", field),
				)
			case "url":
				messages = append(
					messages,
					fmt.Sprintf("%s must be a valid URL", field),
				)
			case "oneof":
				messages = append(
					messages,
					fmt.Sprintf("%s must be one of: %s", field, e.Param()),
				)
			case "min":
				messages = append(
					messages,
					fmt.Sprintf("%s must be at least %s", field, e.Param()),
				)
			case "max":
				messages = append(
					messages,
					fmt.Sprintf("%s must be at most %s", field, e.Param()),
				)
			default:
				messages = append(
					messages,
					fmt.Sprintf("%s validation failed on '%s' tag", field, e.Tag()),
				)
			}
		}
	}

	if len(messages) == 0 {
		return err.Error()
	}

	return strings.Join(messages, "; ")
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			prev := rune(s[i-1])
			if prev >= 'a' && prev <= 'z' {
				result.WriteRune('_')
			}
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
