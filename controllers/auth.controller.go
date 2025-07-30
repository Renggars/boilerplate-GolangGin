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

func (ctrl *authController) Login(ctx *gin.Context) {
	var login dto.LoginRequest
	if err := ctx.ShouldBindJSON(&login); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	responseData, refreshToken, err := ctrl.services.Login(&login)
	if err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	ctx.SetCookie(
		"refreshToken",
		refreshToken,
		1*24*60*60,
		"/",
		"localhost",
		false,
		true,
	)

	res := helpers.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "success login user",
		Data:       responseData,
	})

	ctx.JSON(http.StatusOK, res)
}
