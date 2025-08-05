package services

import (
	"restApi-GoGin/dto"
	"restApi-GoGin/errorhandler"
	"restApi-GoGin/models"
	"restApi-GoGin/repository"
	"restApi-GoGin/utils"
	"time"

	"github.com/google/uuid"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*dto.LoginResponse, string, string, error)
	RefreshToken(refreshToken string) (string, error)
	ForgotPassword(req *dto.ForgotPasswordRequest) error
	VerifyOTP(req *dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error)
	ResetPassword(req *dto.ResetPasswordRequest) error
}

type authService struct {
	authRepository repository.AuthRepository
	userRepository repository.UserRepository
}

func NewAuthService(authRepository repository.AuthRepository, userRepository repository.UserRepository) *authService {
	return &authService{
		authRepository: authRepository,
		userRepository: userRepository,
	}
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	if emailExist := s.authRepository.EmailExists(req.Email); emailExist {
		return &errorhandler.BadRequestError{Message: "email already exists"}
	}

	if req.Password != req.PasswordConfirm {
		return &errorhandler.BadRequestError{Message: "password not match"}
	}

	passwordHash, err := utils.HashBcrypt(req.Password)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: passwordHash,
		Role:     "user",
	}

	if err := s.authRepository.Register(&user); err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil
}

func (s *authService) Login(req *dto.LoginRequest) (*dto.LoginResponse, string, string, error) {
	var data dto.LoginResponse

	user, err := s.userRepository.GetUserByEmail(req.Email)
	if err != nil {
		return nil, "", "", &errorhandler.NotFoundError{Message: "invalid email or password"}
	}

	if err := utils.CompareBcrypt(user.Password, req.Password); err != nil {
		return nil, "", "", &errorhandler.NotFoundError{Message: "invalid email or password"}
	}

	accessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return nil, "", "", &errorhandler.InternalServerError{Message: err.Error()}
	}

	refressToken, err := utils.GenerateRefreshToken(user)
	if err != nil {
		return nil, "", "", &errorhandler.InternalServerError{Message: err.Error()}
	}

	data = dto.LoginResponse{
		ID:    user.Id,
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}

	return &data, accessToken, refressToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, error) {
	claims, err := utils.VerifyRefreshToken(refreshToken)
	if err != nil {
		return "", &errorhandler.UnauthorizedError{Message: err.Error()}
	}

	user, err := s.authRepository.GetUserById(claims.UserId)
	if err != nil {
		return "", &errorhandler.UnauthorizedError{Message: err.Error()}
	}

	newAccessToken, err := utils.GenerateAccessToken(user)
	if err != nil {
		return "", &errorhandler.InternalServerError{Message: err.Error()}
	}

	return newAccessToken, nil
}

func (s *authService) ForgotPassword(req *dto.ForgotPasswordRequest) error {
	user, err := s.userRepository.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return &errorhandler.NotFoundError{Message: "user not found"}
	}

	otp := utils.GenerateOTP()
	hashedOTP, err := utils.HashBcrypt(otp)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	exp := time.Now().Add(time.Minute * 10)
	user.OTPCode = &hashedOTP
	user.OTPCodeExp = &exp

	err = s.userRepository.UpdateUser(user)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	err = utils.SendEmail(user.Email, "OTP Reset Password", otp)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil
}

func (s *authService) VerifyOTP(req *dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error) {
	user, err := s.userRepository.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return nil, &errorhandler.NotFoundError{Message: "user not found"}
	}

	if user.OTPCode == nil || user.OTPCodeExp == nil {
		return nil, &errorhandler.BadRequestError{Message: "no OTP request found"}
	}

	if time.Now().After(*user.OTPCodeExp) {
		return nil, &errorhandler.BadRequestError{Message: "OTP expired"}
	}

	if err := utils.CompareBcrypt(*user.OTPCode, req.OTP); err != nil {
		return nil, &errorhandler.BadRequestError{Message: "invalid otp"}
	}

	rawResetToken := uuid.New().String()
	hashedResetTokenBytes, err := utils.HashBcrypt(rawResetToken)
	if err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	hashedResetToken := string(hashedResetTokenBytes)

	user.ResetToken = &hashedResetToken
	exp := time.Now().Add(time.Minute * 10)
	user.ResetTokenExp = &exp

	user.OTPCode = nil
	user.OTPCodeExp = nil

	if err := s.userRepository.UpdateUser(user); err != nil {
		return nil, &errorhandler.InternalServerError{Message: err.Error()}
	}

	return &dto.VerifyOTPResponse{
		ResetToken: rawResetToken,
	}, nil
}

func (s *authService) ResetPassword(req *dto.ResetPasswordRequest) error {
	user, err := s.userRepository.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return &errorhandler.NotFoundError{Message: "user not found"}
	}

	if user.ResetToken == nil || user.ResetTokenExp == nil {
		return &errorhandler.BadRequestError{Message: "no reset token found"}
	}

	if time.Now().After(*user.ResetTokenExp) {
		return &errorhandler.BadRequestError{Message: "reset token expired"}
	}

	if err := utils.CompareBcrypt(*user.ResetToken, req.ResetToken); err != nil {
		return &errorhandler.BadRequestError{Message: "invalid reset token"}
	}

	if req.Password != req.PasswordConfirm {
		return &errorhandler.BadRequestError{Message: "password not match"}
	}

	passwordHash, err := utils.HashBcrypt(req.Password)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	user.Password = passwordHash
	user.ResetToken = nil
	user.ResetTokenExp = nil

	if err := s.userRepository.UpdateUser(user); err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	return nil
}
