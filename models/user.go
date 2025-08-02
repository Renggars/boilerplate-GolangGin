package models

import "time"

type User struct {
	Id            int        `gorm:"primaryKey" json:"id"`
	Name          string     `gorm:"not null" json:"name"`
	Email         string     `gorm:"unique; not null" json:"email"`
	Password      string     `gorm:"not null" json:"password"`
	Role          string     `gorm:"default:user" json:"role"`
	OTPCode       *string    `gorm:"column:otp_code" json:"-"`
	OTPCodeExp    *time.Time `gorm:"column:otp_code_exp" json:"-"`
	ResetToken    *string    `gorm:"column:reset_token" json:"-"`
	ResetTokenExp *time.Time `gorm:"column:reset_token_exp" json:"-"`
	CreatedAt     time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}
