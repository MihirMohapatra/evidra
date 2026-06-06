package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/evidra/evidra/export/service"
	"github.com/evidra/evidra/pkg/telemetry"
)

func NewRouter(svc *service.ExportService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(telemetry.Middleware)

	reg := telemetry.InitMetrics()
	r.Get("/metrics", telemetry.MetricsHandler(reg).ServeHTTP)
	r.Get("/health", health)

	h := NewHandler(svc)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/exports", h.Export)
		r.Get("/exports/{id}", h.GetExport)
		r.Get("/exports", h.ListExports)
	})

	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
