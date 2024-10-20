package handlers

import (
	"encoding/hex"
	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"strings"
	"svipp-server/assets/templates/pages"
	"svipp-server/internal/auth"
	"svipp-server/internal/httputil"
)

func (h *Handler) FrontPage(w http.ResponseWriter, r *http.Request) {
	err := pages.FrontPage(templ.SafeURL(h.domain)).Render(r.Context(), w)
	if err != nil {
		httputil.InternalServerError(w, err)
	}
}

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	isHxReq := httputil.IsHxRequest(r)
	err := pages.HomePage(isHxReq).Render(r.Context(), w)
	if err != nil {
		httputil.InternalServerError(w, err)
	}
}

func (h *Handler) CreateShopifyApiConfigPage(w http.ResponseWriter, r *http.Request) {
	isHxReq := httputil.IsHxRequest(r)
	err := pages.CreateApiConfigPage(isHxReq).Render(r.Context(), w)
	if err != nil {
		httputil.InternalServerError(w, err)
	}
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := pages.Login().Render(r.Context(), w)
	if err != nil {
		httputil.InternalServerError(w, err)
		return
	}
}

func (h *Handler) SignupPage(w http.ResponseWriter, r *http.Request) {
	// TODO user signup not enabled yet
	//var page templ.Component
	//if r.Host != h.businessSubDomain {
	//	page = pages.UserSignup()
	//} else {
	//	page = pages.BusinessSignup()
	//}

	err := pages.BusinessSignup().Render(r.Context(), w)
	if err != nil {
		httputil.InternalServerError(w, err)
	}
}

func (h *Handler) SingleOrderPage(w http.ResponseWriter, r *http.Request) {
	uuidStr := chi.URLParam(r, "uuid")
	uuidStr = strings.ReplaceAll(uuidStr, "-", "")
	uuidBytes, err := hex.DecodeString(uuidStr)
	if err != nil {
		log.Printf("Error decoding uuid %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if len(uuidBytes) != 16 {
		// Handle the error: UUID should be exactly 16 bytes
		log.Printf("UUID not of length 16")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	pgUuid := pgtype.UUID{
		Bytes: [16]byte(uuidBytes),
		Valid: true,
	}

	order, err := h.db.GetOrderInfoByPublicId(r.Context(), pgUuid)
	if err != nil {
		log.Printf("Error fetching order by public_id %v", err)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	httputil.JSONResponse(w, http.StatusOK, order)
	//httputil.HtmxResponse(w, http.StatusOK, "signup.gohtml", nil)
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {

	if err := h.authService.DeleteSession(r); err != nil {
		log.Printf("Failed to delete session %v", err)
	}

	// Clear the cookie
	http.SetCookie(w, auth.ClearCookie())
	// Redirect to the login page or home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
