package handlers

import (
	"encoding/json"
	"jwtAuth/dto"
	"jwtAuth/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.AuthService
}

func NewAuthHandler(service service.AuthService) AuthHandler {
	return AuthHandler{
		service: service,
	}
}

// Login parse Req body for username and password, return jwt token if valid
func (ah AuthHandler) Login(c *gin.Context) {
	var loginReq dto.LoginRequest
	if err := json.NewDecoder(c.Request.Body).Decode(&loginReq); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return 
	}
	log.Println(loginReq)
	tokens, err := ah.service.Login(loginReq)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return 
	}

	c.String(http.StatusOK, *tokens)
}