package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"svipp-server/internal/httputil"
	"svipp-server/internal/password"
)

type LoginRequest struct {
	Email       string  `json:"email" validate:"required,email"`
	Password    string  `json:"password" validate:"required"`
	DeviceToken *string `json:"deviceToken" validate:"omitempty"`
}

func (h *Handler) Authenticate(writer http.ResponseWriter, request *http.Request) {
	var params LoginRequest
	isHtmx := request.Header.Get("HX-Request") == "true"

	// Parse the request body
	if !isHtmx {
		if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
			httputil.BadRequestResponse(writer, err, false)
			return
		}
	} else {
		if err := request.ParseForm(); err != nil {
			httputil.BadRequestResponse(writer, err, true)
			return
		}
		params.Email = request.FormValue("email")
		params.Password = request.FormValue("password")
	}

	// Validate the request
	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.ValidationFailedResponse(writer, validationErrors, isHtmx)
		return
	}

	user, err := h.db.GetUserByEmail(request.Context(), &params.Email)
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusUnauthorized, fmt.Sprintf("Authentication failed %v", err), "Authentication Failed", isHtmx)
		return
	}

	err = password.CompareWithHash(params.Password, *user.Password)
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusUnauthorized, fmt.Sprintf("Authentication failed %v", err), "Authentication Failed", isHtmx)
		return
	}

	token, err := h.jwtService.GenerateJWT(user.ID, user.Role)
	if err != nil {
		httputil.InternalServerErrorResponse(writer, "Error generating token", err, isHtmx)
		return
	}

	if !isHtmx {
		httputil.JSONResponse(writer, 200, map[string]string{"token": token})
		return
	}

	cookie := h.jwtService.GenerateJwtCookie(token)
	http.SetCookie(writer, &cookie)
	writer.Header().Set("HX-Redirect", "/home")
}
