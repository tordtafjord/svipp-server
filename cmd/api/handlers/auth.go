package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"svipp-server/internal/httputil"
	"svipp-server/internal/password"

	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required,min=8,max=50"`
	DeviceToken *string `json:"device_token" validate:"omitempty"`
}

func (h *Handler) Authenticate(writer http.ResponseWriter, request *http.Request) {
	var params LoginRequest

	// Parse the JSON request body
	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		httputil.BadRequestResponse(writer, err)
		return
	}

	// Validate the request
	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.UnvalidResponse(writer, validationErrors)
		return
	}

	user, err := h.db.GetUserByEmail(request.Context(), &params.Email)
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusUnauthorized, fmt.Sprintf("Authentication failed %v", err), "Authentication Failed")
		return
	}

	err = password.CompareWithHash(params.Password, *user.Password)
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusUnauthorized, fmt.Sprintf("Authentication failed %v", err), "Authentication Failed")
		return
	}

	token, err := h.jwtService.GenerateJWT(user.ID, "user")
	if err != nil {
		httputil.InternalServerErrorResponse(writer, "Error generating token", err)
		return
	}

	httputil.JSONResponse(writer, 200, token)
}

func getReadableValidationErrors(err error) []string {
	var errorMessages []string

	for _, err := range err.(validator.ValidationErrors) {
		switch err.Tag() {
		case "required":
			errorMessages = append(errorMessages, fmt.Sprintf("%s is required", err.Field()))
		case "email":
			errorMessages = append(errorMessages, fmt.Sprintf("%s must be a valid email address", err.Field()))
		case "min":
			errorMessages = append(errorMessages, fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param()))
		default:
			errorMessages = append(errorMessages, fmt.Sprintf("%s failed on the %s tag", err.Field(), err.Tag()))
		}
	}

	return errorMessages
}
