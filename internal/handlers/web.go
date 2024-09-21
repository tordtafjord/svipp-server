package handlers

import (
	"net/http"
	"svipp-server/internal/httputil"
)

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	httputil.HtmxResponse(w, http.StatusOK, "home.gohtml", nil)
}

func (h *Handler) LoginPage(w http.ResponseWriter, r *http.Request) {
	httputil.HtmxResponse(w, http.StatusOK, "login.gohtml", nil)
}
