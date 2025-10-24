package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/okoye-dev/oss-archive/internal/transport/rest"
)

func HealthHandler(c *gin.Context) {
	response := rest.HealthResponse{
		Status:    "Server is healthy!",
		Timestamp: time.Now().Unix(),
		Service:   "OSS Archive",
	}

	rest.Success(c, response)
}