package routes

import (
	"restApi-GoGin/config"
	"restApi-GoGin/controllers"
	"restApi-GoGin/repository"
	"restApi-GoGin/services"

	"github.com/gin-gonic/gin"
)

func AuthRouter(api *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(config.DB)
	authService := services.NewAuthService(authRepository)
	authController := controllers.NewAuthController(authService)

	api.POST("/register", authController.Register)
	api.POST("/login", authController.Login)
	api.POST("/logout", authController.Logout)
	api.POST("/refresh-token", authController.RefreshToken)
}
