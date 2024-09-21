package handlers

import (
	"net/http"
	"svipp-server/internal/httputil"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "home.gohtml", nil)
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to execute home.html", err)
	}
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	err := h.templates.ExecuteTemplate(w, "login.gohtml", nil)
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to execute login.html", err)
	}
}
