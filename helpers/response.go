package helpers

import "restApi-GoGin/dto"

// ResponseWithData represents a successful API response with data
type ResponseWithData struct {
	Code      int           `json:"code" example:"200"`
	Status    string        `json:"status" example:"success"`
	Message   string        `json:"message" example:"Operation completed successfully"`
	Paginaate *dto.Paginate `json:"paginate,omitempty"`
	Data      any           `json:"data"`
}

// ResponseWithoutData represents a successful API response without data
type ResponseWithoutData struct {
	Code    int    `json:"code" example:"200"`
	Status  string `json:"status" example:"success"`
	Message string `json:"message" example:"Operation completed successfully"`
}

// Response creates a standardized API response based on the provided parameters
func Response(params dto.ResponseParams) any {
	var response any
	var status string

	if params.StatusCode >= 200 && params.StatusCode <= 299 {
		status = "success"
	} else {
		status = "error"
	}

	if params.Data != nil {
		response = &ResponseWithData{
			Code:      params.StatusCode,
			Status:    status,
			Message:   params.Message,
			Paginaate: params.Paginate,
			Data:      params.Data,
		}
	} else {
		response = ResponseWithoutData{
			Code:    params.StatusCode,
			Status:  status,
			Message: params.Message,
		}
	}

	return response
}
