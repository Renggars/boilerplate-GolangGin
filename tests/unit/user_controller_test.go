package unit

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"restApi-GoGin/src/controllers"
	"restApi-GoGin/src/models"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

type MockUserService struct {
	getAllUsersFunc    func() ([]models.User, error)
	getUserByEmailFunc func(email string) (*models.User, error)
	getUserByIDFunc    func(id int) (*models.User, error)
	createUserFunc     func(name, email, password, role string) error
	updateUserFunc     func(id int, name, email, password, role *string) error // Tambahkan ini
	deleteUserFunc     func(id int) error
}

func (m *MockUserService) GetAllUsers() ([]models.User, error) {
	if m.getAllUsersFunc != nil {
		return m.getAllUsersFunc()
	}
	return nil, nil
}

func (m *MockUserService) GetUserByEmail(email string) (*models.User, error) {
	if m.getUserByEmailFunc != nil {
		return m.getUserByEmailFunc(email)
	}
	return nil, nil
}

func (m *MockUserService) GetUserByID(id int) (*models.User, error) {
	if m.getUserByIDFunc != nil {
		return m.getUserByIDFunc(id)
	}
	return nil, nil
}

func (m *MockUserService) CreateUser(name, email, password, role string) error {
	if m.createUserFunc != nil {
		return m.createUserFunc(name, email, password, role)
	}
	return nil
}

func (m *MockUserService) UpdateUser(id int, name, email, password, role *string) error {
	if m.updateUserFunc != nil {
		return m.updateUserFunc(id, name, email, password, role)
	}
	return nil
}

func (m *MockUserService) DeleteUser(id int) error {
	if m.deleteUserFunc != nil {
		return m.deleteUserFunc(id)
	}
	return nil
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

func TestGetUserByEmail_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		getUserByEmailFunc: func(email string) (*models.User, error) {
			return &models.User{Id: 1, Name: "User1", Email: email}, nil
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/searchByEmail?email=user1@example.com", nil)

	controller.GetUserByEmail(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		getUserByEmailFunc: func(email string) (*models.User, error) {
			return nil, errors.New("record not found")
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/searchByEmail?email=notfound@example.com", nil)

	controller.GetUserByEmail(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetUserByEmail_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		getUserByEmailFunc: func(email string) (*models.User, error) {
			return nil, errors.New("db error")
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/user/searchByEmail?email=error@example.com", nil)

	controller.GetUserByEmail(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestGetUserByID_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		getUserByIDFunc: func(id int) (*models.User, error) {
			return &models.User{Id: id, Name: "User1", Email: "user1@example.com"}, nil
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	controller.GetUserByID(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestGetUserByID_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		getUserByIDFunc: func(id int) (*models.User, error) {
			return nil, nil
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "2"}}

	controller.GetUserByID(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestGetUserByID_Error(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		getUserByIDFunc: func(id int) (*models.User, error) {
			return nil, errors.New("db error")
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "3"}}

	controller.GetUserByID(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestCreateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		createUserFunc: func(name, email, password, role string) error {
			return nil
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/user", strings.NewReader(`{"name":"User1","email":"user1@example.com","password":"password123","password_confirm":"password123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateUser(c)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestCreateUser_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/user", strings.NewReader(`{"name":"","email":"invalid","password":"123","password_confirm":"321"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateUser(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCreateUser_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		createUserFunc: func(name, email, password, role string) error {
			return errors.New("service error")
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/user", strings.NewReader(`{"name":"User1","email":"user1@example.com","password":"password123","password_confirm":"password123"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.CreateUser(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		updateUserFunc: func(id int, name, email, password, role *string) error {
			return nil
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("PUT", "/user/1", strings.NewReader(`{"name":"User Updated"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.UpdateUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateUser_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "1"}}
	c.Request = httptest.NewRequest("PUT", "/user/1", strings.NewReader(`{"email":"invalid"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.UpdateUser(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		updateUserFunc: func(id int, name, email, password, role *string) error {
			return errors.New("record not found")
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "2"}}
	c.Request = httptest.NewRequest("PUT", "/user/2", strings.NewReader(`{"name":"User Updated"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.UpdateUser(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateUser_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		updateUserFunc: func(id int, name, email, password, role *string) error {
			return errors.New("service error")
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{{Key: "id", Value: "3"}}
	c.Request = httptest.NewRequest("PUT", "/user/3", strings.NewReader(`{"name":"User Updated"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.UpdateUser(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUpdateProfile_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		updateUserFunc: func(id int, name, email, password, role *string) error {
			return nil
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &models.User{Id: 1, Name: "User1", Email: "user1@example.com", Role: "user"})
	c.Request = httptest.NewRequest("PUT", "/user/profile", strings.NewReader(`{"name":"User Updated","email":"updated@example.com"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.UpdateProfile(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpdateProfile_ValidationError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &models.User{Id: 1, Name: "User1", Email: "user1@example.com", Role: "user"})
	c.Request = httptest.NewRequest("PUT", "/user/profile", strings.NewReader(`{"email":"invalid"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.UpdateProfile(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestUpdateProfile_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		updateUserFunc: func(id int, name, email, password, role *string) error {
			return errors.New("service error")
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &models.User{Id: 1, Name: "User1", Email: "user1@example.com", Role: "user"})
	c.Request = httptest.NewRequest("PUT", "/user/profile", strings.NewReader(`{"name":"User Updated"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.UpdateProfile(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestUpdateProfile_InvalidUserContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Tidak set user context
	c.Request = httptest.NewRequest("PUT", "/user/profile", strings.NewReader(`{"name":"User Updated"}`))
	c.Request.Header.Set("Content-Type", "application/json")

	controller.UpdateProfile(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}

func TestDeleteUser_Success_UserDeletingSelf(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		deleteUserFunc: func(id int) error {
			return nil
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &models.User{Id: 1, Name: "User1", Email: "user1@example.com", Role: "user"})
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	controller.DeleteUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check if cookies are cleared when user deletes their own account
	cookies := w.Result().Cookies()
	accessTokenCleared := false
	refreshTokenCleared := false

	for _, cookie := range cookies {
		if cookie.Name == "accessToken" && cookie.Value == "" && cookie.MaxAge == -1 {
			accessTokenCleared = true
		}
		if cookie.Name == "refreshToken" && cookie.Value == "" && cookie.MaxAge == -1 {
			refreshTokenCleared = true
		}
	}

	if !accessTokenCleared || !refreshTokenCleared {
		t.Error("Expected cookies to be cleared when user deletes their own account")
	}
}

func TestDeleteUser_Success_AdminDeletingOtherUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		deleteUserFunc: func(id int) error {
			return nil
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &models.User{Id: 1, Name: "Admin", Email: "admin@example.com", Role: "admin"})
	c.Params = []gin.Param{{Key: "id", Value: "2"}}

	controller.DeleteUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check that cookies are NOT cleared when admin deletes other user
	cookies := w.Result().Cookies()
	for _, cookie := range cookies {
		if (cookie.Name == "accessToken" || cookie.Name == "refreshToken") && cookie.Value == "" && cookie.MaxAge == -1 {
			t.Errorf("Expected cookies NOT to be cleared when admin deletes other user, but %s was cleared", cookie.Name)
		}
	}
}

func TestDeleteUser_Success_AdminDeletingSelf(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		deleteUserFunc: func(id int) error {
			return nil
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &models.User{Id: 1, Name: "Admin", Email: "admin@example.com", Role: "admin"})
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	controller.DeleteUser(c)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// Check if cookies are cleared when admin deletes their own account
	cookies := w.Result().Cookies()
	accessTokenCleared := false
	refreshTokenCleared := false

	for _, cookie := range cookies {
		if cookie.Name == "accessToken" && cookie.Value == "" && cookie.MaxAge == -1 {
			accessTokenCleared = true
		}
		if cookie.Name == "refreshToken" && cookie.Value == "" && cookie.MaxAge == -1 {
			refreshTokenCleared = true
		}
	}

	if !accessTokenCleared || !refreshTokenCleared {
		t.Error("Expected cookies to be cleared when admin deletes their own account")
	}
}

func TestDeleteUser_Forbidden_UserDeletingOtherUser(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &models.User{Id: 1, Name: "User1", Email: "user1@example.com", Role: "user"})
	c.Params = []gin.Param{{Key: "id", Value: "2"}}

	controller.DeleteUser(c)

	if w.Code != http.StatusForbidden {
		t.Errorf("Expected status code %d, got %d", http.StatusForbidden, w.Code)
	}
}

func TestDeleteUser_InvalidUserID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &models.User{Id: 1, Name: "User1", Email: "user1@example.com", Role: "user"})
	c.Params = []gin.Param{{Key: "id", Value: "invalid"}}

	controller.DeleteUser(c)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		deleteUserFunc: func(id int) error {
			return errors.New("record not found")
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &models.User{Id: 999, Name: "User1", Email: "user1@example.com", Role: "user"})
	c.Params = []gin.Param{{Key: "id", Value: "999"}}

	controller.DeleteUser(c)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status code %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestDeleteUser_ServiceError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{
		deleteUserFunc: func(id int) error {
			return errors.New("service error")
		},
	}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Set("user", &models.User{Id: 1, Name: "User1", Email: "user1@example.com", Role: "user"})
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	controller.DeleteUser(c)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
	}
}

func TestDeleteUser_InvalidUserContext(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := &MockUserService{}
	controller := controllers.NewUserController(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Tidak set user context
	c.Params = []gin.Param{{Key: "id", Value: "1"}}

	controller.DeleteUser(c)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
	}
}
