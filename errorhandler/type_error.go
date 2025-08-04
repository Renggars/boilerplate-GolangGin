package errorhandler

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code    int    `json:"code" example:"400"`
	Status  string `json:"status" example:"error"`
	Message string `json:"message" example:"Bad request"`
}

// NotFoundError represents a 404 Not Found error
type NotFoundError struct {
	Message string
}

// BadRequestError represents a 400 Bad Request error
type BadRequestError struct {
	Message string
}

// InternalServerError represents a 500 Internal Server Error
type InternalServerError struct {
	Message string
}

// UnauthorizedError represents a 401 Unauthorized error
type UnauthorizedError struct {
	Message string
}

func (e *NotFoundError) Error() string {
	return e.Message
}

func (e *BadRequestError) Error() string {
	return e.Message
}

func (e *InternalServerError) Error() string {
	return e.Message
}

func (e *UnauthorizedError) Error() string {
	return e.Message
}
