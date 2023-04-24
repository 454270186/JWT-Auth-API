package handlers

import (
	"jwtAuth/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func (ah AuthHandler) Login(c *gin.Context) {
	
}