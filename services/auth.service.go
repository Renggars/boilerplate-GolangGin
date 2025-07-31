package services

import (
	"restApi-GoGin/dto"
	"restApi-GoGin/errorhandler"
	"restApi-GoGin/models"
	"restApi-GoGin/repository"
	"restApi-GoGin/utils"
)

type AuthService interface {
	Register(req *dto.RegisterRequest) error
	Login(req *dto.LoginRequest) (*dto.LoginResponse, string, string, error)
}

type authService struct {
	authRepository repository.AuthRepository
}

func NewAuthService(authRepository repository.AuthRepository) *authService {
	return &authService{authRepository: authRepository}
}

func (s *authService) Register(req *dto.RegisterRequest) error {
	if emailExist := s.authRepository.EmailExists(req.Email); emailExist {
		return &errorhandler.BadRequestError{Message: "email already exists"}
	}

	if req.Password != req.PasswordConfirm {
		return &errorhandler.BadRequestError{Message: "password not match"}
	}

	passwordHash, err := utils.HashPassword(req.Password)
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

	if err := utils.ComparePassword(user.Password, req.Password); err != nil {
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
