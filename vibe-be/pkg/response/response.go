package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidInput = NewAppError(http.StatusBadRequest, "invalid input")
	ErrNotFound     = NewAppError(http.StatusNotFound, "resource not found")
	ErrInternal     = NewAppError(http.StatusInternalServerError, "internal server error")
)

type SuccessResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string `json:"message"`
	Error   any    `json:"error,omitempty"`
}

type AppError struct {
	Code    int
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func NewAppError(code int, msg string) error {
	return AppError{Code: code, Message: msg}
}

func Success(c *gin.Context, msg string, data any) {
	c.JSON(http.StatusOK, SuccessResponse{
		Message: msg,
		Data:    data,
	})
}

func Error(c *gin.Context, err error) {
	if appErr, ok := err.(AppError); ok {
		c.AbortWithStatusJSON(appErr.Code, ErrorResponse{Message: appErr.Message})
		return
	}

	c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorResponse{Message: "Internal Server Error", Error: err.Error()})
}

func ValidationError(c *gin.Context, err map[string]string) {
	c.AbortWithStatusJSON(http.StatusBadRequest, ErrorResponse{Message: "Validation Error", Error: err})
}
