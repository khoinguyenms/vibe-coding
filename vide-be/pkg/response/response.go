package response

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidInput       = errors.New("Invalid input")
	ErrNotFound           = errors.New("Not found")
	ErrInternalError      = errors.New("Internal error")
	ErrEmailAlreadyExists = errors.New("Email already exists")
)

type Response struct {
	Data  any    `json:"data,omitempty"`
	Error string `json:"error,omitempty"`
}

func Success(c *gin.Context, status int, data any) {
	c.JSON(status, Response{Data: data})
}

func Error(c *gin.Context, err error) {
	var status int
	switch {
	case errors.Is(err, ErrInvalidInput):
		status = http.StatusBadRequest
	case errors.Is(err, ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, ErrEmailAlreadyExists):
		status = http.StatusConflict
	default:
		status = http.StatusInternalServerError
	}

	c.AbortWithStatusJSON(status, Response{Error: err.Error()})
}

func ValidationError(c *gin.Context, err error) {
	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		errors := make(map[string]string)

		for _, fieldErr := range validationErrors {
			switch fieldErr.Tag() {
			case "gt":
				errors[fieldErr.Field()] = fmt.Sprintf("%s must be greater than %s", fieldErr.Field(), fieldErr.Param())
			case "lt":
				errors[fieldErr.Field()] = fmt.Sprintf("%s must be less than %s", fieldErr.Field(), fieldErr.Param())
			case "gte":
				errors[fieldErr.Field()] = fmt.Sprintf("%s must be greater than or equal to %s", fieldErr.Field(), fieldErr.Param())
			case "lte":
				errors[fieldErr.Field()] = fmt.Sprintf("%s must be less than or equal to %s", fieldErr.Field(), fieldErr.Param())
			case "required":
				errors[fieldErr.Field()] = fmt.Sprintf("%s is required", fieldErr.Field())
			case "slug":
				errors[fieldErr.Field()] = fmt.Sprintf("%s is not a valid slug", fieldErr.Field())
			case "min":
				errors[fieldErr.Field()] = fmt.Sprintf("%s must be at least %s characters long", fieldErr.Field(), fieldErr.Param())
			case "max":
				errors[fieldErr.Field()] = fmt.Sprintf("%s must be at most %s characters long", fieldErr.Field(), fieldErr.Param())
			case "email":
				errors[fieldErr.Field()] = fmt.Sprintf("%s invalid email format", fieldErr.Field())
			case "email_advanced":
				errors[fieldErr.Field()] = fmt.Sprintf("%s contains a blocked email domain", fieldErr.Field())
			case "strong_password":
				errors[fieldErr.Field()] = fmt.Sprintf("%s must be at least 8 characters long and include uppercase, lowercase, number, and special character", fieldErr.Field())
			case "oneof":
				allowed := strings.Join(strings.Split(fieldErr.Param(), " "), ", ")
				errors[fieldErr.Field()] = fmt.Sprintf("%s must be one of the following: %s", fieldErr.Field(), allowed)
			case "min_int":
				errors[fieldErr.Field()] = fmt.Sprintf("%s must be at least %s", fieldErr.Field(), fieldErr.Param())
			case "file_ext":
				allowed := strings.Join(strings.Split(fieldErr.Param(), " "), ",")
				errors[fieldErr.Field()] = fmt.Sprintf("%s must have one of the following extensions: %s", fieldErr.Field(), allowed)
			default:
				errors[fieldErr.Field()] = fmt.Sprintf("%s is not valid", fieldErr.Field())
			}
		}

		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Validation Error",
			"error":   errors,
		})
	}

	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
		"message": "Validation Error",
		"error":   map[string]string{"detail": err.Error()},
	})
}
