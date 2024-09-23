package handlers

import (
	"net/http"
	"svipp-server/internal/httputil"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	ok := h.jwtService.IsAuthenticated(*r)
	if !ok {
		httputil.HtmxResponse(w, http.StatusOK, "home.gohtml", nil)
		return
	}
	h.FrontPage(w, r)
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	httputil.HtmxResponse(w, http.StatusOK, "login.gohtml", nil)
}

func (h *Handler) FrontPage(w http.ResponseWriter, r *http.Request) {
	httputil.HtmxResponse(w, http.StatusOK, "frontpage.gohtml", nil)
}
