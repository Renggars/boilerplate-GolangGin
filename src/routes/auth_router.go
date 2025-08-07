package routes

import (
	"restApi-GoGin/src/config"
	"restApi-GoGin/src/controllers"
	"restApi-GoGin/src/repository"
	"restApi-GoGin/src/services"

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
	api.POST("/reset-password", authController.ResetPassword)
}
