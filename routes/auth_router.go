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
	userRepository := repository.NewUserRepository(config.DB)
	authService := services.NewAuthService(authRepository, userRepository)
	authController := controllers.NewAuthController(authService)

	api.POST("/register", authController.Register)
	api.POST("/login", authController.Login)
	api.POST("/logout", authController.Logout)
	api.POST("/refresh-token", authController.RefreshToken)
	api.POST("/forgot-password", authController.ForgotPassword)
	api.POST("/verify-otp", authController.VerifyOTP)
}
