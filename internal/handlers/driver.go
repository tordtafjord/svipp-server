package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"svipp-server/internal/database"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
)

func (h *Handler) CreateDriver(writer http.ResponseWriter, request *http.Request) {

	var params createUserRequest

	// Parse the JSON request body
	if err := json.NewDecoder(request.Body).Decode(&params); err != nil {
		httputil.BadRequestResponse(writer, err)
		return
	}

	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.UnvalidResponse(writer, validationErrors)
		return
	}

	user, err := h.db.CreateUser(request.Context(), database.CreateUserParams{
		Name:        params.Name,
		Phone:       params.Phone,
		Email:       &params.Email,
		Password:    &params.Password,
		DeviceToken: params.DeviceToken,
		Role:        models.RoleDriver.String(),
		Temporary:   new(bool),
	})
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusConflict, err.Error(), "User with email or phone already exists")
		return
	}

	err = h.db.CreateDriver(request.Context(), user.ID)
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusInternalServerError, fmt.Sprintf("Error creating driver %v", err), "Internal Server Error")
	}

	httputil.JSONResponse(writer, http.StatusCreated, user)
}
