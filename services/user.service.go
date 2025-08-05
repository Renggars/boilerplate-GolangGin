package services

import (
	"restApi-GoGin/models"
	"restApi-GoGin/repository"
)

// UserService interface
type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}

// userService struct
type userService struct {
	repo repository.UserRepository
}

// NewUserService constructor
func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// GetAllUsers implementation
func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	return s.repo.GetUserByEmail(email)
}
