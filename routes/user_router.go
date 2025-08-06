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

	api.POST(
		"/user",
		middlewares.Auth(authRepository),
		middlewares.AuthAccess(authRepository),
		userController.CreateUser,
	)
	api.GET(
		"/users",
		middlewares.Auth(authRepository),
		middlewares.AuthAccess(authRepository),
		userController.GetAllUsers,
	)
	api.GET("/user/searchByEmail",
		middlewares.Auth(authRepository),
		middlewares.AuthAccess(authRepository),
		userController.GetUserByEmail,
	)
	api.GET(
		"/user/:id",
		middlewares.Auth(authRepository),
		middlewares.AuthAccess(authRepository),
		userController.GetUserByID,
	)
	api.PUT(
		"/user/:id",
		middlewares.Auth(authRepository),
		middlewares.AuthAccess(authRepository),
		userController.UpdateUser,
	)
}
