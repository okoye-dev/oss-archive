package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/okoye-dev/oss-archive/internal/models"
	"github.com/okoye-dev/oss-archive/internal/transport/rest"
)

func CreateUser(c *gin.Context) {
	rest.Success(c, "User created successfully")
}

func GetUser(c *gin.Context) {
	mockUser := models.User{
		ID: "1",
		FirstName: "John",
		LastName: "Doe",
		Email: "john.doe@example.com",
		Password: "password",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	
	rest.Success(c, mockUser)
}

func UpdateUser(c *gin.Context) {
	rest.Success(c, "User updated successfully")
}