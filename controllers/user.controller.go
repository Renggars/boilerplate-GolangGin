package controllers

import (
	"net/http"
	"restApi-GoGin/services"
	"strconv"

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

// GetUserByEmail godoc
// @Summary Get user by email
// @Tags users
// @Produce json
// @Param email query string true "User Email"
// @Success 200 {object} models.User
// @Failure 404 {object} errorhandler.NotFoundError
// @Failure 500 {object} errorhandler.InternalServerError
// @Security BearerAuth
// @Router /user/searchByEmail [get]
func (ctrl *UserController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Query("email")
	user, err := ctrl.service.GetUserByEmail(email)
	if err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// GetUserByID godoc
// @Summary Get user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} models.User
// @Failure 404 {object} errorhandler.NotFoundError
// @Failure 500 {object} errorhandler.InternalServerError
// @Security BearerAuth
// @Router /user/{id} [get]
func (ctrl *UserController) GetUserByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}
	user, err := ctrl.service.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to get user"})
		return
	}
	if user == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}
