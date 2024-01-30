package main

import (
	"github.com/gin-gonic/gin"
	"user-app/internal/app/routes"
)

func main() {
	engine := gin.Default()
	routes.RegisterRoutes(engine)
	err := engine.Run("localhost:8080")
	if err != nil {
		return
	}
	println("Listening and serving on 8080....")
}
