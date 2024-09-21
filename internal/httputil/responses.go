package httputil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

type errorResponse struct {
	Error string `json:"error"`
}

type multipleErrorsResponse struct {
	Error []string `json:"error"`
}

func BadRequestResponse(writer http.ResponseWriter, err error, isHtmx bool) {
	payload := errorResponse{Error: "Bad Request"}
	log.Printf("Bad Reguest: %s", err.Error())
	if !isHtmx {
		JSONResponse(writer, http.StatusBadRequest, payload)
		return
	}
	HtmxResponse(writer, http.StatusBadRequest, "error.gohtml", payload)
}

func InternalServerErrorResponse(writer http.ResponseWriter, logMessage string, err error, isHtmx bool) {
	// Get the caller's file and line number
	_, file, line, ok := runtime.Caller(1)
	if ok {
		// Extract just the filename from the full path
		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]
	}
	detailedLogMessage := fmt.Sprintf(
		"Server Error [%d] - %s\nLocation: %s:%d\n",
		http.StatusInternalServerError, fmt.Sprintf(logMessage, err), file, line,
	)
	// Log the detailed message
	log.Printf(detailedLogMessage)

	// Optionally, you could also send this to an error tracking service
	// sendToErrorTrackingService(detailedLogMessage)

	payload := errorResponse{Error: "Internal Server Error"}
	if !isHtmx {
		JSONResponse(writer, http.StatusInternalServerError, payload)
		return
	}
	HtmxResponse(writer, http.StatusInternalServerError, "error.gohtml", payload)
}

func UnauthorizedResponse(writer http.ResponseWriter, isHtmx bool) {
	if !isHtmx {
		JSONResponse(writer, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
		return
	}
	// For HTMX requests, set the HX-Redirect header to redirect to the login page
	writer.Header().Set("HX-Redirect", "/login")
	writer.WriteHeader(http.StatusUnauthorized)
}

func ForbiddenResponse(writer http.ResponseWriter, isHtmx bool) {
	if !isHtmx {
		JSONResponse(writer, http.StatusForbidden, errorResponse{Error: "Unauthorized"})
		return
	}
	// For HTMX requests, set the HX-Redirect header to redirect to the home page
	writer.Header().Set("HX-Redirect", "/")
	writer.WriteHeader(http.StatusForbidden)
}

func ErrorResponse(writer http.ResponseWriter, code int, logMessage, userMessage string, isHtmx bool) {
	log.Println(logMessage)
	payload := errorResponse{Error: userMessage}
	if !isHtmx {
		JSONResponse(writer, code, payload)
		return
	}
	HtmxResponse(writer, code, "error.gohtml", payload)
}

func ValidationFailedResponse(writer http.ResponseWriter, validationMessages []string, isHtmx bool) {
	payload := multipleErrorsResponse{Error: validationMessages}
	if !isHtmx {
		JSONResponse(writer, http.StatusBadRequest, payload)
		return
	}
	HtmxResponse(writer, http.StatusOK, "error.gohtml", payload)
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

func HtmxResponse(writer http.ResponseWriter, code int, template string, data interface{}) {
	writer.Header().Set("Content-Type", "text/html")
	writer.WriteHeader(code)
	err := Tmpl.ExecuteTemplate(writer, template, data)
	if err != nil {
		log.Printf("Failed executing template: %v", err)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
