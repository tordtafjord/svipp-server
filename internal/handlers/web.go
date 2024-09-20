package handlers

import (
	"net/http"
	"svipp-server/internal/httputil"
)

func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title string
		// Add other fields as needed
	}{
		Title: "Svipp – Levering på dine premisser",
	}

	err := h.templates.ExecuteTemplate(w, "layout.html", data)
	if err != nil {
		httputil.InternalServerErrorResponse(w, "Failed to execute layout.html", err)
	}
}
