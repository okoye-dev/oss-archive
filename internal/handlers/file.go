package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/okoye-dev/oss-archive/internal/services"
	"github.com/okoye-dev/oss-archive/internal/transport/rest"
)

func GetFiles(c *gin.Context) {
	files, err := services.GetFiles()
	if err != nil {
		rest.InternalError(c, err)
		return
	}
	
	rest.Success(c, files)
}