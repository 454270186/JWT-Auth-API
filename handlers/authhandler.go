package handlers

import (
	"encoding/json"
	"jwtAuth/dto"
	"jwtAuth/service"
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

	tokens, err := ah.service.Login(loginReq)
	if err != nil {
		c.String(http.StatusUnauthorized, err.Error())
		return 
	}

	c.String(http.StatusOK, *tokens)
}

func (ah AuthHandler) Verify(c *gin.Context) {
	urlParams := make(map[string]string)

	for k := range c.Request.URL.Query() {
		urlParams[k] = c.Request.URL.Query().Get(k)
	}

	if urlParams["token"] != "" {
		isAuthorized, err := ah.service.Verify(urlParams)
		if err != nil || !isAuthorized {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "not authorize response",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"Verified": "Authorized",
		})
	} else {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "missing token",
		})
	}
}