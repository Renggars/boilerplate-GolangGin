package controllers

import (
	"net/http"
	"restApi-GoGin/dto"
	"restApi-GoGin/errorhandler"
	"restApi-GoGin/helpers"
	"restApi-GoGin/services"

	"github.com/gin-gonic/gin"
)

type authController struct {
	services services.AuthService
}

func NewAuthController(authService services.AuthService) *authController {
	return &authController{
		services: authService,
	}
}

func (ctrl *authController) Register(ctx *gin.Context) {
	var register dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&register); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := ctrl.services.Register(&register); err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	response := helpers.Response(dto.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "success register user",
	})

	ctx.JSON(http.StatusCreated, response)
}
