package unit

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"restApi-GoGin/controllers"
	"restApi-GoGin/dto"
	"testing"

	"github.com/gin-gonic/gin"
)

// Mock untuk AuthService
type MockAuthService struct {
	registerFunc       func(*dto.RegisterRequest) error
	loginFunc          func(*dto.LoginRequest) (*dto.LoginResponse, string, string, error)
	refreshTokenFunc   func(string) (string, error)
	forgotPasswordFunc func(*dto.ForgotPasswordRequest) error
	verifyOTPFunc      func(*dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error)
	resetPasswordFunc  func(*dto.ResetPasswordRequest) error
}

func (m *MockAuthService) Register(register *dto.RegisterRequest) error {
	if m.registerFunc != nil {
		return m.registerFunc(register)
	}
	return nil
}

func (m *MockAuthService) Login(login *dto.LoginRequest) (*dto.LoginResponse, string, string, error) {
	if m.loginFunc != nil {
		return m.loginFunc(login)
	}
	return nil, "", "", nil
}

func (m *MockAuthService) RefreshToken(refreshToken string) (string, error) {
	if m.refreshTokenFunc != nil {
		return m.refreshTokenFunc(refreshToken)
	}
	return "", nil
}

func (m *MockAuthService) ForgotPassword(forgotPassword *dto.ForgotPasswordRequest) error {
	if m.forgotPasswordFunc != nil {
		return m.forgotPasswordFunc(forgotPassword)
	}
	return nil
}

func (m *MockAuthService) VerifyOTP(verifyOTP *dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error) {
	if m.verifyOTPFunc != nil {
		return m.verifyOTPFunc(verifyOTP)
	}
	return nil, nil
}

func (m *MockAuthService) ResetPassword(resetPassword *dto.ResetPasswordRequest) error {
	if m.resetPasswordFunc != nil {
		return m.resetPasswordFunc(resetPassword)
	}
	return nil
}

func TestRegister_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := &MockAuthService{
		registerFunc: func(register *dto.RegisterRequest) error {
			return nil
		},
	}
	controller := controllers.NewAuthController(mockService)

	// Test data
	registerData := dto.RegisterRequest{
		Name:            "Test User",
		Email:           "test@example.com",
		Password:        "password123",
		PasswordConfirm: "password123",
	}

	// Create request
	jsonData, _ := json.Marshal(registerData)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Execute
	controller.Register(c)

	// Assertions
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "success register user" {
		t.Errorf("Expected message 'success register user', got '%v'", response["message"])
	}
}

func TestRegister_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockAuthService{}
	controller := controllers.NewAuthController(mockService)

	// Test data dengan email tidak valid
	registerData := map[string]interface{}{
		"name":             "Test User",
		"email":            "invalid-email",
		"password":         "password123",
		"password_confirm": "password123",
	}

	jsonData, _ := json.Marshal(registerData)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.Register(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockAuthService{
		loginFunc: func(login *dto.LoginRequest) (*dto.LoginResponse, string, string, error) {
			return &dto.LoginResponse{
				ID:    1,
				Name:  "Test User",
				Email: "test@example.com",
				Role:  "user",
			}, "access_token", "refresh_token", nil
		},
	}
	controller := controllers.NewAuthController(mockService)

	loginData := dto.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.Login(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "success login user" {
		t.Errorf("Expected message 'success login user', got '%v'", response["message"])
	}

	if response["data"] == nil {
		t.Error("Expected data in response, got nil")
	}
}

func TestLogout_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockAuthService{}
	controller := controllers.NewAuthController(mockService)

	req, _ := http.NewRequest("POST", "/logout", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.Logout(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "success logout user" {
		t.Errorf("Expected message 'success logout user', got '%v'", response["message"])
	}

	// Check cookies are cleared
	cookies := w.Result().Cookies()
	var accessTokenCleared, refreshTokenCleared bool
	for _, cookie := range cookies {
		if cookie.Name == "accessToken" && cookie.Value == "" {
			accessTokenCleared = true
		}
		if cookie.Name == "refreshToken" && cookie.Value == "" {
			refreshTokenCleared = true
		}
	}
	if !accessTokenCleared {
		t.Error("Expected accessToken cookie to be cleared")
	}
	if !refreshTokenCleared {
		t.Error("Expected refreshToken cookie to be cleared")
	}
}

func TestRefreshToken_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockAuthService{
		refreshTokenFunc: func(refreshToken string) (string, error) {
			return "new_access_token", nil
		},
	}
	controller := controllers.NewAuthController(mockService)

	req, _ := http.NewRequest("POST", "/refresh-token", nil)
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{
		Name:  "refreshToken",
		Value: "test_refresh_token",
	})

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.RefreshToken(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "success refresh token" {
		t.Errorf("Expected message 'success refresh token', got '%v'", response["message"])
	}
}

func TestRefreshToken_NoToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockAuthService{}
	controller := controllers.NewAuthController(mockService)

	// Create request without cookie
	req, _ := http.NewRequest("POST", "/refresh-token", nil)
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.RefreshToken(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestForgotPassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockAuthService{
		forgotPasswordFunc: func(forgotPassword *dto.ForgotPasswordRequest) error {
			return nil
		},
	}
	controller := controllers.NewAuthController(mockService)

	forgotPasswordData := dto.ForgotPasswordRequest{
		Email: "test@example.com",
	}

	jsonData, _ := json.Marshal(forgotPasswordData)
	req, _ := http.NewRequest("POST", "/forgot-password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.ForgotPassword(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "OTP has been sent to your email" {
		t.Errorf("Expected message 'OTP has been sent to your email', got '%v'", response["message"])
	}
}

func TestVerifyOTP_Success(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := &MockAuthService{
		verifyOTPFunc: func(verifyOTP *dto.VerifyOTPRequest) (*dto.VerifyOTPResponse, error) {
			return &dto.VerifyOTPResponse{
				ResetToken: "test_reset_token",
			}, nil
		},
	}
	controller := controllers.NewAuthController(mockService)

	verifyOTPData := dto.VerifyOTPRequest{
		Email: "test@example.com",
		OTP:   "123456",
	}

	jsonData, _ := json.Marshal(verifyOTPData)
	req, _ := http.NewRequest("POST", "/verify-otp", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.VerifyOTP(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "OTP verified successfully" {
		t.Errorf("Expected message 'OTP verified successfully', got '%v'", response["message"])
	}

	if response["data"] == nil {
		t.Error("Expected data in response, got nil")
	}
}

func TestResetPassword_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockAuthService{
		resetPasswordFunc: func(resetPassword *dto.ResetPasswordRequest) error {
			return nil
		},
	}
	controller := controllers.NewAuthController(mockService)

	resetPasswordData := dto.ResetPasswordRequest{
		Email:           "test@example.com",
		ResetToken:      "test_reset_token",
		Password:        "newpassword123",
		PasswordConfirm: "newpassword123",
	}

	jsonData, _ := json.Marshal(resetPasswordData)
	req, _ := http.NewRequest("POST", "/reset-password", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	controller.ResetPassword(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response["message"] != "password reset successfully" {
		t.Errorf("Expected message 'password reset successfully', got '%v'", response["message"])
	}
}
