package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/evidra/evidra/evidence/service"
	"github.com/evidra/evidra/pkg/telemetry"
)

func NewRouter(svc *service.EvidenceService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(telemetry.Middleware)

	reg := telemetry.InitMetrics()
	r.Get("/metrics", telemetry.MetricsHandler(reg).ServeHTTP)
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	h := NewHandler(svc)

	r.Route("/api/v1/evidence", func(r chi.Router) {
		r.Post("/", h.Create)
		r.Get("/", h.List)
		r.Get("/{id}", h.Get)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)

		r.Post("/{id}/submit", h.Submit)
		r.Post("/{id}/approve", h.Approve)
		r.Post("/{id}/reject", h.Reject)
		r.Post("/{id}/export", h.Export)
		r.Get("/{id}/approvals", h.GetApprovalHistory)
	})

	return r
}
