package services

import (
	"restApi-GoGin/dto"
	"restApi-GoGin/errorhandler"
	"restApi-GoGin/models"
	"restApi-GoGin/repository"
	"restApi-GoGin/utils"
	"time"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*dto.LoginResponse, string, string, error)
	RefreshToken(refreshToken string) (string, error)
	ForgotPassword(req *dto.ForgotPasswordRequest) error
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

	user, err := s.authRepository.GetUserByEmail(req.Email)
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
	user, err := s.authRepository.GetUserByEmail(req.Email)
	if err != nil || user == nil {
		return &errorhandler.NotFoundError{Message: "user not found"}
	}

	otp := utils.GenerateOTP()
	hashedOTP, err := utils.HashBcrypt(otp)
	if err != nil {
		return &errorhandler.InternalServerError{Message: err.Error()}
	}

	exp := time.Now().Add(time.Minute * 5)
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
