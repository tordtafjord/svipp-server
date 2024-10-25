package httputil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"svipp-server/assets/templates/components"
)

type errorResponse struct {
	Error string `json:"error"`
}

type multipleErrorsResponse struct {
	Error []string `json:"error"`
}

func BadRequestResponse(writer http.ResponseWriter, err error, isHtmx bool) {
	log.Printf("Bad Reguest: %s", err.Error())
	if !isHtmx {
		JSONResponse(writer, http.StatusBadRequest, errorResponse{Error: "Bad Request"})
		return
	}
	HtmxResponse(writer, http.StatusBadRequest, "error.gohtml", multipleErrorsResponse{Error: []string{"Bad Request"}})
}

func InternalServerError(w http.ResponseWriter, err error) {
	// Get the caller's file and line number
	_, file, line, ok := runtime.Caller(1)
	if ok {
		// Extract just the filename from the full path
		parts := strings.Split(file, "/")
		file = parts[len(parts)-1]
	}
	detailedLogMessage := fmt.Sprintf(
		"Server Error [%d] - %s\nLocation: %s:%d\n",
		http.StatusInternalServerError, err, file, line,
	)
	// Log the detailed message
	log.Printf(detailedLogMessage)

	http.Error(w, "Internal Server Error", http.StatusInternalServerError)
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

	if !isHtmx {
		JSONResponse(writer, http.StatusInternalServerError, errorResponse{Error: "Internal Server Error"})
		return
	}
	HtmxResponse(writer, http.StatusOK, "error.gohtml", multipleErrorsResponse{Error: []string{"Internal Server Error"}})
}

func UnauthorizedResponse(writer http.ResponseWriter) {
	JSONResponse(writer, http.StatusUnauthorized, errorResponse{Error: "Unauthorized"})
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

func YellowToastResponse(w http.ResponseWriter, r *http.Request, errors []string) {
	w.Header().Add("HX-Retarget", "#toasts")
	err := components.YellowToasts(errors).Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render yellow toast template.")
		InternalServerError(w, err)
	}
}

func RedToastResponse(w http.ResponseWriter, r *http.Request, error string) {
	err := components.RedToast(error).Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render red toast template.")
		InternalServerError(w, err)
	}
}

func ErrorResponse(writer http.ResponseWriter, code int, logMessage, userMessage string, isHtmx bool) {
	log.Println(logMessage)
	if !isHtmx {
		JSONResponse(writer, code, errorResponse{Error: userMessage})
		return
	}
	HtmxResponse(writer, http.StatusOK, "error.gohtml", multipleErrorsResponse{Error: []string{userMessage}})
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
	// FIXME
}
