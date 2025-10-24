package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp int64  `json:"timestamp"`
	Service   string `json:"service"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func SuccessWithMessage(c *gin.Context, message string, data interface{}) {
	response := map[string]interface{}{
		"message": message,
		"data":    data,
	}
	c.JSON(http.StatusOK, response)
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, ErrorResponse{
		Error:   message,
		Code:    code,
		Message: message,
	})
}

func InternalError(c *gin.Context, err error) {
	Error(c, http.StatusInternalServerError, err.Error())
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, message)
}