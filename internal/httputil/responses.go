package httputil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}
type errorsResponse struct {
	Errors []string `json:"errors"`
}

func BadRequestResponse(writer http.ResponseWriter, err error) {
	ErrorResponse(writer, http.StatusBadRequest, fmt.Sprintf("Error decoding json: %v", err), "Bad Request")
}

func InternalServerErrorResponse(writer http.ResponseWriter, logMessage string, err error) {
	ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("%s: %v", logMessage, err), "Internal Server Error")
}

func UnauthorizedResponse(writer http.ResponseWriter, logMessage string, err error) {
	ErrorResponse(writer, http.StatusUnauthorized, fmt.Sprintf("%s: %v", logMessage, err), "Unauthorized")
}

func ForbiddenResponse(writer http.ResponseWriter, logMessage string, err error) {
	ErrorResponse(writer, http.StatusForbidden, fmt.Sprintf("%s: %v", logMessage, err), "Unauthorized")
}

func UnvalidResponse(writer http.ResponseWriter, userMessages []string) {
	if len(userMessages) == 1 {
		JSONResponse(writer, http.StatusBadRequest, errorResponse{Error: userMessages[0]})
	} else {
		JSONResponse(writer, http.StatusBadRequest, errorsResponse{Errors: userMessages})
	}
}

func ErrorResponse(writer http.ResponseWriter, code int, logMessage string, userMessage string) {
	// Log the error with detailed internal message if it's a server error (500 and above)
	if code >= 500 {
		log.Printf("Responding with %v: %s", code, logMessage)
	}

	JSONResponse(writer, code, errorResponse{Error: userMessage})
}

func JSONResponse(writer http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(dat)
}

func OKResponse(writer http.ResponseWriter) {
	writer.WriteHeader(http.StatusOK)
}
