package server

import (
	"github.com/gin-gonic/gin"
	"github.com/okoye-dev/oss-archive/internal/handlers"
	"github.com/okoye-dev/oss-archive/internal/storage"
)

func SetupRoutes(router *gin.Engine, storage storage.StorageInterface) {
	api := router.Group("/api/v1")
	
	setupHealthRoutes(api)
	setupUserRoutes(api)
	setupFileRoutes(api, storage)
}

func setupHealthRoutes(rg *gin.RouterGroup) {
	health := rg.Group("/health")
	health.GET("", handlers.HealthHandler)
}

func setupUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.POST("", handlers.CreateUser)
}

func setupFileRoutes(rg *gin.RouterGroup, storage storage.StorageInterface) {
	fileHandler := handlers.NewFileHandler(storage)
	
	files := rg.Group("/files")
	files.GET("", fileHandler.GetFiles)
	files.POST("", fileHandler.UploadFile)
	files.GET("/:id", fileHandler.GetFile)
	files.DELETE("/:id", fileHandler.DeleteFile)
}