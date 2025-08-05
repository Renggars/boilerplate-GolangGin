package unit

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"restApi-GoGin/controllers"
	"restApi-GoGin/models"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockUserService struct {
	getAllUsersFunc func() ([]models.User, error)
}

func (m *MockUserService) GetAllUsers() ([]models.User, error) {
	if m.getAllUsersFunc != nil {
		return m.getAllUsersFunc()
	}
	return nil, nil
}

func TestGetAllUsers_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		getAllUsersFunc: func() ([]models.User, error) {
			return []models.User{{Id: 1, Name: "User1", Email: "user1@example.com"}}, nil
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	controller.GetAllUsers(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetAllUsers_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		getAllUsersFunc: func() ([]models.User, error) {
			return nil, errors.New("mock error")
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	controller.GetAllUsers(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}
