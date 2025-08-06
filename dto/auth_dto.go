package dto

// RegisterRequest represents the request body for user registration
type RegisterRequest struct {
	Name            string `json:"name" validate:"required" example:"John Doe"`
	Email           string `json:"email" validate:"required,email" example:"john@example.com"`
	Password        string `json:"password" validate:"required,min=6" example:"password123"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password" example:"password123"`
}

// LoginRequest represents the request body for user login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email" example:"john@example.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// LoginResponse represents the response body for successful login
type LoginResponse struct {
	ID    int    `json:"id" example:"1"`
	Name  string `json:"name" example:"John Doe"`
	Email string `json:"email" example:"john@example.com"`
	Role  string `json:"role" example:"user"`
}

// ForgotPasswordRequest represents the request body for forgot password
type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
}

// VerifyOTPRequest represents the request body for OTP verification
type VerifyOTPRequest struct {
	Email string `json:"email" validate:"required,email" example:"john@example.com"`
	OTP   string `json:"otp" validate:"required" example:"123456"`
}

// VerifyOTPResponse represents the response body for successful OTP verification
type VerifyOTPResponse struct {
	ResetToken string `json:"reset_token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ResetPasswordRequest represents the request body for password reset
type ResetPasswordRequest struct {
	Email           string `json:"email" validate:"required,email" example:"john@example.com"`
	ResetToken      string `json:"reset_token" validate:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	Password        string `json:"password" validate:"required,min=6" example:"newpassword123"`
	PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password" example:"newpassword123"`
}

// UpdateUserRequest represents the request body for updating a user
// swagger:model
// @Description Update user fields. Password and role are optional.
type UpdateUserRequest struct {
	Name     string `json:"name" validate:"omitempty" example:"John Doe"`
	Email    string `json:"email" validate:"omitempty,email" example:"john@example.com"`
	Password string `json:"password" validate:"omitempty,min=6" example:"newpassword123"`
	Role     string `json:"role" validate:"omitempty" example:"user"`
}

// UpdateProfileRequest represents the request body for updating user profile (name and email only)
// swagger:model
// @Description Update user profile fields (name and email only)
type UpdateProfileRequest struct {
	Name  string `json:"name" validate:"omitempty" example:"John Doe"`
	Email string `json:"email" validate:"omitempty,email" example:"john@example.com"`
}
