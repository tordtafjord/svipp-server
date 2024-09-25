package handlers

import (
	"encoding/hex"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"strings"
	"svipp-server/internal/httputil"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	if h.jwtService.IsAuthenticated(*r) {
		h.FrontPage(w, r)
		return
	}
	httputil.HtmxResponse(w, http.StatusOK, "home.gohtml", nil)
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	if h.jwtService.IsAuthenticated(*r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	httputil.HtmxResponse(w, http.StatusOK, "login.gohtml", nil)
}

func (h *Handler) FrontPage(w http.ResponseWriter, r *http.Request) {
	httputil.HtmxResponse(w, http.StatusOK, "frontpage.gohtml", nil)
}

func (h *Handler) SignupPage(w http.ResponseWriter, r *http.Request) {
	httputil.HtmxResponse(w, http.StatusOK, "signup.gohtml", nil)
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
	// Clear the JWT cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	// Redirect to the login page or home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
