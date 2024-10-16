package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"svipp-server/internal/auth"
	"svipp-server/internal/database"
	"svipp-server/internal/httputil"
	"svipp-server/internal/models"
	"svipp-server/internal/password"
)

type createBusinessAccountRequest struct {
	FirstName       *string `json:"firstName" validate:"required,min=2,max=100"`
	LastName        *string `json:"lastName" validate:"required,min=2,max=100"`
	Phone           string  `json:"phone" validate:"required,e164"`
	Email           string  `json:"email" validate:"required,email"`
	Password        string  `json:"password" validate:"required,min=8,max=64"`
	ConfirmPassword string  `json:"confirmPassword" validate:"required,min=8,max=64,eqfield=Password"`
	CompanyName     string  `json:"companyName" validate:"required,min=1,max=255"`
	OrgNumber       string  `json:"orgNumber" validate:"required,len=9,numeric"`
	BusinessAddress string  `json:"businessAddress" validate:"required,min=1,max=255"`
	City            string  `json:"city" validate:"required,min=1,max=100"`
	ZipCode         string  `json:"zipCode" validate:"required,min=4,max=10,numeric"`
}

func (h *Handler) CreateBusiness(writer http.ResponseWriter, request *http.Request) {
	var params createBusinessAccountRequest
	isHtmx := httputil.IsNotJson(request)

	// Parse the signup request body
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

		params.OrgNumber = request.FormValue("orgNumber")
		firstName := request.FormValue("firstName")
		lastName := request.FormValue("lastName")
		params.FirstName = &firstName
		params.LastName = &lastName
		params.Phone = fmt.Sprintf("+%s%s", request.FormValue("countryCode"), request.FormValue("phone"))
		params.Email = request.FormValue("email")
		params.Password = request.FormValue("password")
		params.ConfirmPassword = request.FormValue("confirmPassword")
		params.BusinessAddress = request.FormValue("businessAddress")
		params.CompanyName = request.FormValue("companyName")
		params.ZipCode = request.FormValue("zipCode")
		params.City = request.FormValue("city")
	}

	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.ValidationFailedResponse(writer, validationErrors, isHtmx)
		return
	}
	// orgnr valid, parse to int64:
	orgNr, err := strconv.ParseInt(params.OrgNumber, 10, 64)
	if err != nil {
		httputil.BadRequestResponse(writer, err, true)
		return
	}

	pswdHash, err := password.Hash(params.Password)
	if err != nil {
		httputil.InternalServerErrorResponse(writer, "Failed hashing password", err, isHtmx)
		return
	}

	ctx := request.Context()
	user, err := h.db.CreateUser(ctx, database.CreateUserParams{
		FirstName: params.FirstName,
		LastName:  params.LastName,
		Phone:     &params.Phone,
		Email:     &params.Email,
		Password:  &pswdHash,
		Temporary: new(bool), // Use new(bool) to create a pointer to false
		Role:      models.RoleBusiness.String(),
	})
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusConflict, fmt.Sprintf("Failed to create user: %v", err), "Oppretting av konto mislyktes, har du en fra før av?", isHtmx)
		return
	}
	// TODO: Double check this does not overwrite exising non temp users

	addr := fmt.Sprintf("%s, %s %s", params.BusinessAddress, params.ZipCode, params.City)
	err = h.db.CreateBusiness(ctx, database.CreateBusinessParams{
		ID:      user.ID,
		Name:    params.CompanyName,
		OrgID:   orgNr,
		Address: addr,
	})
	if err != nil {
		httputil.ErrorResponse(writer, http.StatusConflict, fmt.Sprintf("Failed to create user: %v", err), "Oppretting av konto mislyktes, har du en fra før av?", isHtmx)
		return
	}

	sessionId, err := h.authService.CreateSession(ctx, user.ID, models.Role(user.Role))
	if err != nil {
		httputil.InternalServerErrorResponse(writer, "Error creating session", err, isHtmx)
		return
	}

	cookie := auth.CreateCookie(sessionId)
	http.SetCookie(writer, &cookie)
	writer.Header().Set("HX-Redirect", "/")
}
