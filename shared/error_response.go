package shared

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (e ErrorResponse) GetStatus() int {
	return e.Status
}

func (e ErrorResponse) Error() string {
	return e.Message
}

func NewBadRequestError(message string) ErrorResponse {
	return ErrorResponse{
		Status:  400,
		Message: message,
	}
}

func NewInternalServerError(message string) ErrorResponse {
	return ErrorResponse{
		Status:  500,
		Message: message,
	}
}

func NewNotFoundError(message string) ErrorResponse {
	return ErrorResponse{
		Status:  404,
		Message: message,
	}
}
