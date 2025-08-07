package routes

import (
	"restApi-GoGin/src/config"
	"restApi-GoGin/src/controllers"
	"restApi-GoGin/src/middleware"
	"restApi-GoGin/src/repository"
	"restApi-GoGin/src/services"

	"github.com/gin-gonic/gin"
)

func UserRouter(api *gin.RouterGroup) {
	authRepository := repository.NewAuthRepository(config.DB)
	userRepository := repository.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	api.POST(
		"/user",
		middleware.Auth(authRepository),
		middleware.AuthAccess(authRepository),
		userController.CreateUser,
	)
	api.GET(
		"/users",
		middleware.Auth(authRepository),
		middleware.AuthAccess(authRepository),
		userController.GetAllUsers,
	)
	api.GET("/user/searchByEmail",
		middleware.Auth(authRepository),
		middleware.AuthAccess(authRepository),
		userController.GetUserByEmail,
	)
	api.GET(
		"/user/:id",
		middleware.Auth(authRepository),
		middleware.AuthAccess(authRepository),
		userController.GetUserByID,
	)
	api.PUT(
		"/user/:id",
		middleware.Auth(authRepository),
		middleware.AuthAccess(authRepository),
		userController.UpdateUser,
	)
	api.PUT(
		"/user/profile",
		middleware.Auth(authRepository),
		userController.UpdateProfile,
	)
	api.DELETE(
		"/user/:id",
		middleware.Auth(authRepository),
		userController.DeleteUser,
	)
}
