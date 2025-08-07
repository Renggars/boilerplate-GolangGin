# API Documentation

## Overview

The API documentation is generated using [swaggo/swag](https://github.com/swaggo/swag) and served using [gin-swagger](https://github.com/swaggo/gin-swagger).

## Accessing the Documentation

Once the server is running, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

## API Endpoints

### Authentication Endpoints

- `POST /api/register` - Register a new user
- `POST /api/login` - Login user
- `POST /api/logout` - Logout user
- `POST /api/refresh-token` - Refresh access token
- `POST /api/forgot-password` - Send OTP for password reset
- `POST /api/verify-otp` - Verify OTP and get reset token
- `POST /api/reset-password` - Reset password using reset token

### User Endpoints

- `POST /api/user` - Create a new user
- `PUT /api/user/profile` - Update user profile (name and email)
- `GET /api/user/searchByEmail` - Get user by email
- `GET /api/user/{id}` - Get user by ID
- `PUT /api/user/{id}` - Update user by ID
- `DELETE /api/user/{id}` - Delete user by ID
- `GET /api/users` - Get all users

### Health Check

- `GET /api/ping` - Health check endpoint

## Regenerating Documentation

To regenerate the Swagger documentation after making changes to the code:

```bash
# Install swag if not already installed
go install github.com/swaggo/swag/cmd/swag@latest

# Generate documentation
swag init -g main.go
```
