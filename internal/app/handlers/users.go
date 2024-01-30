package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"user-app/internal/app/services"
)

type Users interface {
	GetUsers(ctx *gin.Context)
}

type users struct {
	usersService services.Users
}

func NewUsers(usersService services.Users) Users {
	return &users{usersService: usersService}
}

func (u users) GetUsers(ctx *gin.Context) {
	response := u.usersService.GetUsers()
	fmt.Println(response)
	ctx.AbortWithStatusJSON(http.StatusOK, response)
}
