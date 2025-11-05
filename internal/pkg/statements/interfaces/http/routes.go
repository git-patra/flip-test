package httpiface

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *Handler) {
	r.Post("/statements", h.Upload)
	r.Get("/balance", h.GetBalance)
	r.Get("/transactions/issues", h.GetIssues)
}
