package shared

import (
	"encoding/json"
	"log"
	"net/http"
)

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ErrorResponseWithStructure struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
}

// SendSuccessResponse sends success response as JSON
func SendSuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := SuccessResponse{
		Status:  http.StatusOK,
		Message: message,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode success response: %v", err)
		SendErrorResponse(w, NewInternalServerError("Failed to encode response"))
	}
}

// SendErrorResponse sends error response as JSON
func SendErrorResponse(w http.ResponseWriter, errResp ErrorResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errResp.Status)

	response := ErrorResponseWithStructure(errResp)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode error response: %v", err)
		http.Error(w, errResp.Message, errResp.Status)
	}
}

// SendSuccessResponseWithStatus sends success response with custom status code
func SendSuccessResponseWithStatus(w http.ResponseWriter, data interface{}, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]interface{}{
		"status":  statusCode,
		"message": message,
		"data":    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode success response: %v", err)
		SendErrorResponse(w, NewInternalServerError("Failed to encode response"))
	}
}

// SendMessage sends a simple message response
func SendMessage(w http.ResponseWriter, message string, statusCode int) {
	response := map[string]interface{}{
		"message": message,
		"status":  statusCode,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode message response: %v", err)
		http.Error(w, message, statusCode)
	}
}
