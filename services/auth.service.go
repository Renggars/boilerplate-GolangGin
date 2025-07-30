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
