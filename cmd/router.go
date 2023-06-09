package main

import (
	"database/sql"
	"jwtAuth/domain"
	"jwtAuth/handlers"
	"jwtAuth/service"

	"github.com/gin-gonic/gin"
)

func router(DB *sql.DB) *gin.Engine {
	router := gin.Default()

	ah := handlers.NewAuthHandler(service.NewAuthService(domain.NewAuthRepo(DB)))

	// TODO: 
	// 1. /auth/login -- DONE
	// 2. /auth/register
	// 3. /auth/verify -- DONE
	// 4. /refresh -- generate a access token using refresh token
	router.POST("/auth/login", ah.Login)
	router.GET("/auth/verify", ValidateToken(), ah.Verify)
	router.POST("/refresh", ah.Refresh)

	return router
}