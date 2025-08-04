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

	responseData, accessToken, refreshToken, err := ctrl.services.Login(&login)
	if err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	ctx.SetCookie(
		"accessToken",
		accessToken,
		15*60,
		"/",
		"localhost",
		false,
		true,
	)

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

func (ctrl *authController) Logout(ctx *gin.Context) {
	ctx.SetCookie(
		"accessToken",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)

	ctx.SetCookie(
		"refreshToken",
		"",
		-1,
		"/",
		"localhost",
		false,
		true,
	)

	res := helpers.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "success logout user",
	})

	ctx.JSON(http.StatusOK, res)
}

func (ctrl *authController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.UnauthorizedError{Message: err.Error()})
		return
	}

	newAccessToken, err := ctrl.services.RefreshToken(refreshToken)
	if err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	ctx.SetCookie(
		"accessToken",
		newAccessToken,
		15*60,
		"/",
		"localhost",
		false,
		true,
	)

	res := helpers.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "success refresh token",
	})

	ctx.JSON(http.StatusOK, res)
}

func (ctrl *authController) ForgotPassword(ctx *gin.Context) {
	var request dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := ctrl.services.ForgotPassword(&request); err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	res := helpers.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "OTP sent to your email",
	})

	ctx.JSON(http.StatusOK, res)
}

func (ctrl *authController) VerifyOTP(ctx *gin.Context) {
	var request dto.VerifyOTPRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	ResetToken, err := ctrl.services.VerifyOTP(&request)
	if err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	res := helpers.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "success verify otp",
		Data:       ResetToken,
	})

	ctx.JSON(http.StatusOK, res)
}

func (ctrl *authController) ResetPassword(ctx *gin.Context) {
	var request dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := ctrl.services.ResetPassword(&request); err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	res := helpers.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "password reset successfully",
	})

	ctx.JSON(http.StatusOK, res)
}
