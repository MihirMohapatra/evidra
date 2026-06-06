package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/evidra/evidra/audit/service"
)

func NewRouter(svc *service.AuditService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	h := NewHandler(svc)

	r.Get("/health", h.Health)

	r.Route("/api/v1/audit", func(r chi.Router) {
		r.Post("/events", h.Record)
		r.Get("/events", h.List)
	})

	return r
}
