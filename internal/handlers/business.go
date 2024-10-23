package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"svipp-server/assets/templates/components"
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
	LocationName    string  `json:"locationName" validate:"required,min=1,max=100"`
}

type NewShopifyConfigForm struct {
	LocationName       string `json:"locationName" validate:"required"`
	UseShopifyAddress  bool   `json:"useShopifyAddress"`
	Address            string `json:"address" validate:"omitempty"`
	ZipCode            string `json:"zipCode" validate:"omitempty,numeric"`
	City               string `json:"city" validate:"omitempty"`
	PickupInstructions string `json:"pickupInstructions" validate:"omitempty"`
	PickupWindows      PickupWindowsStruct
}

type PickupWindow struct {
	Start string `json:"start" validate:"omitempty,datetime=15:04"`
	End   string `json:"end" validate:"omitempty,datetime=15:04"`
}

type PickupWindowsStruct struct {
	Monday    PickupWindow
	Tuesday   PickupWindow
	Wednesday PickupWindow
	Thursday  PickupWindow
	Friday    PickupWindow
	Saturday  PickupWindow
	Sunday    PickupWindow
}

func (h *Handler) CreateShopifyConfig(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form create shopify config %v", err)
		httputil.RedToastResponse(w, r, "Error: Bad Request")
		return
	}

	pickupWindows := PickupWindowsStruct{
		Monday:    PickupWindow{Start: r.FormValue("0Start"), End: r.FormValue("0End")},
		Tuesday:   PickupWindow{Start: r.FormValue("1Start"), End: r.FormValue("1End")},
		Wednesday: PickupWindow{Start: r.FormValue("2Start"), End: r.FormValue("2End")},
		Thursday:  PickupWindow{Start: r.FormValue("3Start"), End: r.FormValue("3End")},
		Friday:    PickupWindow{Start: r.FormValue("4Start"), End: r.FormValue("4End")},
		Saturday:  PickupWindow{Start: r.FormValue("5Start"), End: r.FormValue("5End")},
		Sunday:    PickupWindow{Start: r.FormValue("6Start"), End: r.FormValue("6End")},
	}

	params := NewShopifyConfigForm{
		LocationName:       r.FormValue("locationName"),
		Address:            r.FormValue("address"),
		ZipCode:            r.FormValue("zipCode"),
		City:               r.FormValue("city"),
		PickupInstructions: r.FormValue("pickupInstructions"),
		UseShopifyAddress:  r.FormValue("useShopifyAddress") == "true" || r.FormValue("useShopifyAddress") == "on",
		PickupWindows:      pickupWindows,
	}
	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.YellowToastResponse(w, r, validationErrors)
		return
	}

	userId, err := auth.GetUserIdFromCtx(r.Context())
	if err != nil {
		log.Printf("Failed to get user id")
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}

	var addr *string
	if !params.UseShopifyAddress {
		// Create the address string and assign its address to addr
		address := fmt.Sprintf("%s, %s %s", params.Address, params.ZipCode, params.City)
		addr = &address
	}

	shopifyConfig := database.CreateShopifyApiKeyParams{
		BusinessID:         userId,
		LocationName:       &params.LocationName,
		PickupAddress:      addr,
		PickupCoords:       nil,
		PickupInstructions: &params.PickupInstructions,
	}

	_, token, err := h.authService.CreateShopifyApiKey(r.Context(), shopifyConfig)
	if err != nil {
		log.Printf("Failed to create shopify api config")
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}

	//TODO: INSERT BUSINESS HOURS

	err = components.ApiKeyModal(token).Render(r.Context(), w)
	if err != nil {
		log.Printf("Failed to render api key template")
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}
}

func (h *Handler) CreateBusiness(w http.ResponseWriter, r *http.Request) {
	var params createBusinessAccountRequest

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form create business %v", err)
		httputil.RedToastResponse(w, r, "Error: Bad Request")
		return
	}

	params.OrgNumber = r.FormValue("orgNumber")
	firstName := r.FormValue("firstName")
	lastName := r.FormValue("lastName")
	params.FirstName = &firstName
	params.LastName = &lastName
	params.Phone = fmt.Sprintf("+%s%s", r.FormValue("countryCode"), r.FormValue("phone"))
	params.Email = r.FormValue("email")
	params.Password = r.FormValue("password")
	params.ConfirmPassword = r.FormValue("confirmPassword")
	params.BusinessAddress = r.FormValue("businessAddress")
	params.CompanyName = r.FormValue("companyName")
	params.ZipCode = r.FormValue("zipCode")
	params.City = r.FormValue("city")
	params.LocationName = r.FormValue("locationName")

	if validationErrors := validateStruct(params); validationErrors != nil {
		httputil.YellowToastResponse(w, r, validationErrors)
		return
	}
	// orgnr valid, parse to int64:
	orgNr, err := strconv.ParseInt(params.OrgNumber, 10, 64)
	if err != nil {
		log.Printf("Feilet orgnr convert %v", err)
		httputil.RedToastResponse(w, r, "Ugyldig Orgnr")
		return
	}

	pswdHash, err := password.Hash(params.Password)
	if err != nil {
		log.Printf("Feilet pswdhash %v", err)
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}

	ctx := r.Context()
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
		log.Printf("Oppretting av bedriftskonto mislyktes %v", err)
		httputil.YellowToastResponse(w, r, []string{"Oppretting av konto feilet. Har du en konto i fra før?"})
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
		log.Printf("Oppretting av bedriftskonto mislyktes %v", err)
		httputil.YellowToastResponse(w, r, []string{"Oppretting av konto feilet. Har du en konto i fra før?"})
		return
	}

	sessionId, err := h.authService.CreateSession(ctx, user.ID, models.Role(user.Role))
	if err != nil {
		log.Printf("Error creating session %v", err)
		httputil.RedToastResponse(w, r, "Error: Intern Serverfeil")
		return
	}

	cookie := auth.CreateCookie(sessionId)
	http.SetCookie(w, &cookie)
	w.Header().Set("HX-Redirect", "/")
}
