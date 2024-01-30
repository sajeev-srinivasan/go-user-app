package routes

import (
	"github.com/gin-gonic/gin"
	"user-app/internal/app/handlers"
	"user-app/internal/app/services"
)

func RegisterRoutes(engine *gin.Engine) {
	usersService := services.NewUsers()
	userHandler := handlers.NewUsers(usersService)
	userGroup := engine.Group("/api/v1")
	{
		userGroup.GET("/users", userHandler.GetUsers)
	}
}
