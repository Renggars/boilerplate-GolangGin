package repository

import (
	"restApi-GoGin/models"

	"gorm.io/gorm"
)

type AuthRepository interface {
	EmailExists(email string) bool
	Register(user *models.User) error
	GetUserById(id int) (*models.User, error)
}

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) EmailExists(email string) bool {
	var user models.User
	err := r.db.First(&user, "email = ?", email).Error

	return err == nil
}

func (r *authRepository) Register(user *models.User) error {
	err := r.db.Create(&user).Error

	return err
}

func (r *authRepository) GetUserById(id int) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error

	return &user, err
}
