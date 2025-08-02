package repository

import (
	"restApi-GoGin/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	UpdateUser(user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) UpdateUser(user *models.User) error {
	return r.db.Save(user).Error
}
