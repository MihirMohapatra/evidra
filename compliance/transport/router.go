package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/evidra/evidra/compliance/service"
	"github.com/evidra/evidra/pkg/telemetry"
)

func NewRouter(svc *service.ComplianceService) *chi.Mux {
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
		// Frameworks
		r.Post("/compliance/frameworks", h.CreateFramework)
		r.Get("/compliance/frameworks", h.ListFrameworks)
		r.Get("/compliance/frameworks/{frameworkId}", h.GetFramework)
		r.Delete("/compliance/frameworks/{frameworkId}", h.DeleteFramework)

		// Controls
		r.Post("/compliance/frameworks/{frameworkId}/controls", h.CreateControl)
		r.Get("/compliance/frameworks/{frameworkId}/controls", h.ListControls)

		// Evidence mappings
		r.Post("/compliance/mappings", h.MapEvidence)
		r.Delete("/compliance/mappings/{id}", h.UnmapEvidence)
		r.Get("/compliance/mappings/by-control/{controlId}", h.ListMappingsByControl)

		// Coverage
		r.Get("/compliance/frameworks/{frameworkId}/coverage", h.GetFrameworkCoverage)
	})

	return r
}

func health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
