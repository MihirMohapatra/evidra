package transport

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/evidra/evidra/pkg/telemetry"
	"github.com/evidra/evidra/questionnaire/service"
)

func NewRouter(svc *service.QuestionnaireService) *chi.Mux {
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

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/questionnaires/upload", h.Upload)
		r.Get("/questionnaires", h.ListQuestionnaires)
		r.Get("/questionnaires/{id}", h.GetQuestionnaire)
		r.Delete("/questionnaires/{id}", h.Delete)
		r.Get("/questionnaires/{id}/questions", h.GetQuestions)
	})

	return r
}
