package controllers

import (
	"net/http"
	"restApi-GoGin/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	service services.UserService
}

func NewUserController(service services.UserService) *UserController {
	return &UserController{service: service}
}

// GetAllUsers godoc
// @Summary Get all users
// @Tags users
// @Produce json
// @Success 200 {array} models.User
// @Failure 500 {object} errorhandler.InternalServerError
// @Router /users [get]
func (ctrl *UserController) GetAllUsers(ctx *gin.Context) {
	users, err := ctrl.service.GetAllUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get users"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}
