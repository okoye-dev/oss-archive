package server

import (
	"github.com/gin-gonic/gin"
	"github.com/okoye-dev/oss-archive/internal/handlers"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api/v1")
	
	setupHealthRoutes(api)
	setupUserRoutes(api)
	setupFileRoutes(api)
}

func setupHealthRoutes(rg *gin.RouterGroup) {
	health := rg.Group("/health")
	health.GET("", handlers.HealthHandler)
}

func setupUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.POST("", handlers.CreateUser)
	// users.GET("/:id", handlers.GetUser)
	// users.PUT("/:id", handlers.UpdateUser)
	// users.DELETE("/:id", handlers.DeleteUser)
}

func setupFileRoutes(rg *gin.RouterGroup) {
	files := rg.Group("/files")
	files.GET("", handlers.GetFiles)
	// files.POST("", handlers.UploadFile)
	// files.GET("/:id", handlers.GetFile)
	// files.DELETE("/:id", handlers.DeleteFile)
}