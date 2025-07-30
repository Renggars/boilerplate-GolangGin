package errorhandler

import (
	"restApi-GoGin/dto"
	"restApi-GoGin/helpers"

	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context, err error) {
	var statusCode int

	switch err.(type) {
	case *NotFoundError:
		statusCode = 404
	case *BadRequestError:
		statusCode = 400
	case *UnauthorizedError:
		statusCode = 401
	case *InternalServerError:
		statusCode = 500
	}

	response := helpers.Response(dto.ResponseParams{
		StatusCode: statusCode,
		Message:    err.Error(),
	})

	c.JSON(statusCode, response)
}
