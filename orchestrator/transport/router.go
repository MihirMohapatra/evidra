package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/evidra/evidra/orchestrator/service"
	"github.com/evidra/evidra/pkg/telemetry"
)

func NewRouter(svc *service.OrchestratorService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(telemetry.Middleware)

	reg := telemetry.InitMetrics()
	r.Get("/metrics", telemetry.MetricsHandler(reg).ServeHTTP)

	h := NewHandler(svc)
	r.Get("/health", h.Health)

	r.Route("/api/v1/orchestrator", func(r chi.Router) {
		r.Post("/answer", h.Answer)
		r.Get("/drafts", h.ListDrafts)
		r.Get("/drafts/{id}", h.GetDraft)
		r.Post("/drafts/{id}/approve", h.ApproveDraft)
		r.Post("/drafts/{id}/reject", h.RejectDraft)
	})

	return r
}
