package routes

import (
	"restApi-GoGin/config"
	"restApi-GoGin/controllers"
	"restApi-GoGin/middlewares"
	"restApi-GoGin/repository"
	"restApi-GoGin/services"

	"github.com/gin-gonic/gin"
)

func UserRouter(api *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(config.DB)
	userRepository := repository.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	api.GET("/users", middlewares.Auth(authRepository), middlewares.AuthAccess(authRepository), userController.GetAllUsers)
}
