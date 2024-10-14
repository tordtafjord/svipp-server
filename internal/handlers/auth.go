package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"svipp-server/internal/auth"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
	"svipp-server/internal/password"
)

type LoginRequest struct {
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required"`
	DeviceToken *string `json:"deviceToken" validate:"omitempty"`
}

func (h *Handler) Authenticate(writer http.ResponseWriter, request *http.Request) {
	var params LoginRequest
	notJson := httputil.IsNotJson(request)

	// Parse the request body
	if notJson {
		if err := request.ParseForm(); err != nil {
			httputil.BadRequestResponse(writer, err, true)
			return
		}
		params.Email = request.FormValue("email")
		params.Password = request.FormValue("password")
	} else {
		if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
			httputil.BadRequestResponse(writer, err, false)
			return
		}
	}

	// Validate the request
	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.ValidationFailedResponse(writer, validationErrors, notJson)
		return
	}

	ctx := request.Context()
	user, err := h.db.GetUserByEmail(ctx, &params.Email)
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusUnauthorized, fmt.Sprintf("Authentication failed %v", err), "Authentication Failed", notJson)
		return
	}

	err = password.CompareWithHash(params.Password, *user.Password)
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusUnauthorized, fmt.Sprintf("Authentication failed %v", err), "Authentication Failed", notJson)
		return
	}

	token, err := h.authService.CreateSession(ctx, user.ID, models.Role(user.Role))
	if err != nil {
		httputil.InternalServerErrorResponse(writer, "Error generating token", err, notJson)
		return
	}

	cookie := auth.CreateCookie(token)
	http.SetCookie(writer, &cookie)
	writer.Header().Set("HX-Redirect", "/home")
}
