package services

import (
	"restApi-GoGin/models"
	"restApi-GoGin/repository"

	"gorm.io/gorm"
)

// UserService interface
type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
	CreateUser(name, email, password, role string) error
	UpdateUser(id int, name, email, password, role *string) error
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

func (s *userService) GetUserByID(id int) (*models.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *userService) CreateUser(name, email, password, role string) error {
	user := &models.User{
		Name:     name,
		Email:    email,
		Password: password,
		Role:     role,
	}
	return s.repo.CreateUser(user)
}

func (s *userService) UpdateUser(id int, name, email, password, role *string) error {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return gorm.ErrRecordNotFound
	}
	if name != nil {
		user.Name = *name
	}
	if email != nil {
		user.Email = *email
	}
	if password != nil {
		user.Password = *password
	}
	if role != nil {
		user.Role = *role
	}
	return s.repo.UpdateUser(user)
}
