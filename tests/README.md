# Unit Testing Documentation

This folder contains unit tests for all endpoints in the Go Gin boilerplate API.

## Folder Structure

```
tests/
├── README.md                    # This file
└── unit/
    ├── auth_controller_test.go  # Unit tests for auth controller
    └── user_controller_test.go  # Unit tests for user controller
```

## Running Tests

To run all unit tests:
```bash
go test ./tests/unit/... -v
```

## Types of Unit Tests

### Auth Controller Tests
- `TestRegister_Success` - Register success
- `TestRegister_InvalidRequest` - Register with invalid request
- `TestLogin_Success` - Login success
- `TestLogout_Success` - Logout success
- `TestRefreshToken_Success` - Refresh token success
- `TestRefreshToken_NoToken` - Refresh token with no token
- `TestForgotPassword_Success` - Forgot password success
- `TestVerifyOTP_Success` - Verify OTP success
- `TestResetPassword_Success` - Reset password success

### User Controller Tests
- `TestGetAllUsers_Success` - Get all users success
- `TestGetAllUsers_Error` - Get all users error
- `TestGetUserByEmail_Success` - Get user by email success
- `TestGetUserByEmail_NotFound` - Get user by email not found
- `TestGetUserByEmail_Error` - Get user by email error
- `TestGetUserByID_Success` - Get user by ID success
- `TestGetUserByID_NotFound` - Get user by ID not found
- `TestGetUserByID_Error` - Get user by ID error
- `TestCreateUser_Success` - Create user success
- `TestCreateUser_ValidationError` - Create user validation error
- `TestCreateUser_ServiceError` - Create user service error
- `TestUpdateUser_Success` - Update user success
- `TestUpdateUser_ValidationError` - Update user validation error
- `TestUpdateUser_NotFound` - Update user not found
- `TestUpdateUser_ServiceError` - Update user service error
- `TestUpdateProfile_Success` - Update profile success
- `TestUpdateProfile_ValidationError` - Update profile validation error
- `TestUpdateProfile_ServiceError` - Update profile service error
- `TestUpdateProfile_InvalidUserContext` - Update profile invalid user context
- `TestDeleteUser_Success_UserDeletingSelf` - User deletes self success
- `TestDeleteUser_Success_AdminDeletingOtherUser` - Admin deletes other user success
- `TestDeleteUser_Success_AdminDeletingSelf` - Admin deletes self success
- `TestDeleteUser_Forbidden_UserDeletingOtherUser` - User forbidden to delete other user
- `TestDeleteUser_InvalidUserID` - Delete user with invalid user ID
- `TestDeleteUser_UserNotFound` - Delete user not found
- `TestDeleteUser_ServiceError` - Delete user service error
- `TestDeleteUser_InvalidUserContext` - Delete user invalid user context


## How to Add a New Test

1. Create a test function with the `Test` prefix
2. Use mock services for dependencies
3. Test both success and error scenarios
4. Use descriptive names

Example:
```go
func TestNewFunction_Success(t *testing.T) {
    // Setup
    mockService := &MockAuthService{
        someFunc: func() error {
            return nil
        },
    }
    controller := controllers.NewAuthController(mockService)
    // Test logic here
    // Assertions here
}
```
