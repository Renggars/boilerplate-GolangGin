package controllers

import (
	"net/http"
	"restApi-GoGin/dto"
	"restApi-GoGin/errorhandler"
	"restApi-GoGin/services"
	"restApi-GoGin/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type authController struct {
	services services.AuthService
}

func NewAuthController(authService services.AuthService) *authController {
	return &authController{
		services: authService,
	}
}

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object true "Register Request"
// @Success 201 {object} object
// @Failure 400 {object} object
// @Failure 500 {object} object
// @Router /register [post]
func (ctrl *authController) Register(ctx *gin.Context) {
	var register dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&register); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := validate.Struct(register); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := ctrl.services.Register(&register); err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	response := utils.Response(dto.ResponseParams{
		StatusCode: http.StatusCreated,
		Message:    "success register user",
	})

	ctx.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary Login user
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object true "Login Request"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /login [post]
func (ctrl *authController) Login(ctx *gin.Context) {
	var login dto.LoginRequest
	if err := ctx.ShouldBindJSON(&login); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := validate.Struct(login); err != nil {
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

	res := utils.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "success login user",
		Data:       responseData,
	})

	ctx.JSON(http.StatusOK, res)
}

// Logout godoc
// @Summary Logout user
// @Description Logout user by clearing cookies
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} object
// @Router /logout [post]
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

	res := utils.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "success logout user",
	})

	ctx.JSON(http.StatusOK, res)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Refresh access token using refresh token from cookie
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} object
// @Failure 401 {object} object
// @Failure 500 {object} object
// @Router /refresh-token [post]
func (ctrl *authController) RefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.UnauthorizedError{Message: "refresh token not found"})
		return
	}

	accessToken, err := ctrl.services.RefreshToken(refreshToken)
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

	res := utils.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "success refresh token",
	})

	ctx.JSON(http.StatusOK, res)
}

// ForgotPassword godoc
// @Summary Forgot password
// @Description Send OTP to user's email for password reset
// @Tags auth
// @Accept json
// @Produce json
// @Param request body object true "Forgot Password Request"
// @Success 200 {object} object
// @Failure 400 {object} object
// @Failure 404 {object} object
// @Failure 500 {object} object
// @Router /forgot-password [post]
func (ctrl *authController) ForgotPassword(ctx *gin.Context) {
	var forgotPassword dto.ForgotPasswordRequest
	if err := ctx.ShouldBindJSON(&forgotPassword); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := validate.Struct(forgotPassword); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := ctrl.services.ForgotPassword(&forgotPassword); err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	res := utils.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "OTP has been sent to your email",
	})

	ctx.JSON(http.StatusOK, res)
}

// VerifyOTP godoc
// @Summary Verify OTP
// @Description Verify OTP and get reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.VerifyOTPRequest true "Verify OTP Request"
// @Success 200 {object} utils.ResponseWithData{data=dto.VerifyOTPResponse}
// @Failure 400 {object} errorhandler.ErrorResponse
// @Failure 401 {object} errorhandler.ErrorResponse
// @Failure 500 {object} errorhandler.ErrorResponse
// @Router /verify-otp [post]
func (ctrl *authController) VerifyOTP(ctx *gin.Context) {
	var verifyOTP dto.VerifyOTPRequest
	if err := ctx.ShouldBindJSON(&verifyOTP); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := validate.Struct(verifyOTP); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	resetToken, err := ctrl.services.VerifyOTP(&verifyOTP)
	if err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	res := utils.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "OTP verified successfully",
		Data:       resetToken,
	})

	ctx.JSON(http.StatusOK, res)
}

// ResetPassword godoc
// @Summary Reset password
// @Description Reset password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.ResetPasswordRequest true "Reset Password Request"
// @Success 200 {object} utils.ResponseWithoutData
// @Failure 400 {object} errorhandler.ErrorResponse
// @Failure 401 {object} errorhandler.ErrorResponse
// @Failure 500 {object} errorhandler.ErrorResponse
// @Router /reset-password [post]
func (ctrl *authController) ResetPassword(ctx *gin.Context) {
	var resetPassword dto.ResetPasswordRequest
	if err := ctx.ShouldBindJSON(&resetPassword); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := validate.Struct(resetPassword); err != nil {
		errorhandler.ErrorHandler(ctx, &errorhandler.BadRequestError{Message: err.Error()})
		return
	}

	if err := ctrl.services.ResetPassword(&resetPassword); err != nil {
		errorhandler.ErrorHandler(ctx, err)
		return
	}

	res := utils.Response(dto.ResponseParams{
		StatusCode: http.StatusOK,
		Message:    "password reset successfully",
	})

	ctx.JSON(http.StatusOK, res)
}
