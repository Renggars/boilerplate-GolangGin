package models

import "time"

type User struct {
	ID       int    `gorm:"primaryKey" json:"id"`
	Name     string `gorm:"not null" json:"name"`
	Email    string `gorm:"unique; not null" json:"email"`
	Password string `gorm:"not null" json:"password"`
	Role     string `gorm:"default:user" json:"role"`
	CreateAt time.Time
	UpdateAt time.Time
}
