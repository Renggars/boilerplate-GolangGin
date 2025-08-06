package controllers

import (
	"net/http"
	"restApi-GoGin/services"
	"strconv"

	"restApi-GoGin/dto"
	"restApi-GoGin/models"
	"restApi-GoGin/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validateUser = validator.New()

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

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user (admin only)
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "User Data"
// @Success 201 {object} utils.ResponseWithoutData "Created"
// @Failure 400 {object} errorhandler.BadRequestError
// @Failure 401 {object} errorhandler.UnauthorizedError
// @Failure 403 {object} errorhandler.ForbiddenError
// @Failure 500 {object} errorhandler.InternalServerError
// @Security BearerAuth
// @Router /user [post]
func (ctrl *UserController) CreateUser(ctx *gin.Context) {
	var req dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validateUser.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	passwordHash, err := utils.HashBcrypt(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}
	role := "user"
	if ctx.PostForm("role") != "" {
		role = ctx.PostForm("role")
	}
	if err := ctrl.service.CreateUser(req.Name, req.Email, passwordHash, role); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	response := utils.Response(dto.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "success create user",
	})
	ctx.JSON(http.StatusCreated, response)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user data (name, email, password, role). Hanya field yang diisi yang akan diupdate.
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body dto.UpdateUserRequest true "User Data"
// @Success 200 {object} utils.ResponseWithoutData "OK"
// @Failure 400 {object} errorhandler.BadRequestError
// @Failure 401 {object} errorhandler.UnauthorizedError
// @Failure 404 {object} errorhandler.NotFoundError
// @Failure 500 {object} errorhandler.InternalServerError
// @Security BearerAuth
// @Router /user/{id} [put]
func (ctrl *UserController) UpdateUser(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	var req dto.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validateUser.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var namePtr, emailPtr, passwordPtr, rolePtr *string
	if req.Name != "" {
		namePtr = &req.Name
	}
	if req.Email != "" {
		emailPtr = &req.Email
	}
	if req.Password != "" {
		hash, err := utils.HashBcrypt(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
			return
		}
		passwordPtr = &hash
	}
	if req.Role != "" {
		rolePtr = &req.Role
	}
	if err := ctrl.service.UpdateUser(id, namePtr, emailPtr, passwordPtr, rolePtr); err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := utils.Response(dto.ResponseParams{
		StatusCode: 200,
		Message:    "success update user",
	})
	ctx.JSON(http.StatusOK, response)
}

// UpdateProfile godoc
// @Summary Update user profile (name and email only)
// @Description Update name and email for the authenticated user. Admin can update any user by id.
// @Tags users
// @Accept json
// @Produce json
// @Param request body dto.UpdateProfileRequest true "Profile Data"
// @Success 200 {object} utils.ResponseWithoutData "OK"
// @Failure 400 {object} errorhandler.BadRequestError
// @Failure 401 {object} errorhandler.UnauthorizedError
// @Failure 403 {object} errorhandler.ForbiddenError
// @Failure 404 {object} errorhandler.NotFoundError
// @Failure 500 {object} errorhandler.InternalServerError
// @Security BearerAuth
// @Router /user/profile [put]
func (ctrl *UserController) UpdateProfile(ctx *gin.Context) {
	userObj, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	user, ok := userObj.(*models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user context"})
		return
	}

	var req dto.UpdateProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := validateUser.Struct(req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var namePtr, emailPtr *string
	if req.Name != "" {
		namePtr = &req.Name
	}
	if req.Email != "" {
		emailPtr = &req.Email
	}

	id := user.Id
	if user.Role == "admin" && ctx.Query("id") != "" {
		if idParam, err := strconv.Atoi(ctx.Query("id")); err == nil {
			id = idParam
		}
	}
	if err := ctrl.service.UpdateUser(id, namePtr, emailPtr, nil, nil); err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	response := utils.Response(dto.ResponseParams{
		StatusCode: 200,
		Message:    "success update profile",
	})
	ctx.JSON(http.StatusOK, response)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user. User can only delete themselves, admin can delete any user by ID.
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} utils.ResponseWithoutData "OK"
// @Failure 400 {object} errorhandler.BadRequestError
// @Failure 401 {object} errorhandler.UnauthorizedError
// @Failure 403 {object} errorhandler.ForbiddenError
// @Failure 404 {object} errorhandler.NotFoundError
// @Failure 500 {object} errorhandler.InternalServerError
// @Security BearerAuth
// @Router /user/{id} [delete]
func (ctrl *UserController) DeleteUser(ctx *gin.Context) {
	// Get authenticated user from context
	userObj, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	user, ok := userObj.(*models.User)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Invalid user context"})
		return
	}

	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	if user.Role != "admin" && user.Id != id {
		ctx.JSON(http.StatusForbidden, gin.H{"message": "Access denied: you can only delete your own account"})
		return
	}

	if err := ctrl.service.DeleteUser(id); err != nil {
		if err.Error() == "record not found" {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user"})
		return
	}

	// If user is deleting their own account (including admin), clear cookies to logout
	// If admin is deleting other user, the deleted user will be automatically logged out
	// when they try to access any protected endpoint due to DeletedAt check in middleware
	if user.Id == id {
		ctx.SetCookie("accessToken", "", -1, "/", "", false, true)
		ctx.SetCookie("refreshToken", "", -1, "/", "", false, true)
	}

	response := utils.Response(dto.ResponseParams{
		StatusCode: 200,
		Message:    "success delete user",
	})
	ctx.JSON(http.StatusOK, response)
}
