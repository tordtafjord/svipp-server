package handlers

import (
	"net/http"
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
