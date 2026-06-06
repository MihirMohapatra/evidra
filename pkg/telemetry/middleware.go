package telemetry

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		routePattern := getRoutePattern(r)
		tracer := otel.Tracer("evidra")
		ctx, span := tracer.Start(r.Context(), r.Method+" "+routePattern,
			trace.WithAttributes(
				attribute.String("http.method", r.Method),
				attribute.String("http.route", routePattern),
				attribute.String("http.target", r.URL.Path),
			),
		)
		defer span.End()

		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r.WithContext(ctx))

		status := sw.status
		duration := time.Since(start).Seconds()

		span.SetAttributes(attribute.Int("http.status_code", status))

		HTTPRequestCount.WithLabelValues(r.Method, routePattern, strconv.Itoa(status)).Inc()
		HTTPRequestDuration.WithLabelValues(r.Method, routePattern).Observe(duration)

		if status >= 500 {
			HTTPRequestErrors.WithLabelValues(r.Method, routePattern, strconv.Itoa(status)).Inc()
			span.SetAttributes(attribute.Bool("error", true))
		}
	})
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func getRoutePattern(r *http.Request) string {
	routeCtx := chi.RouteContext(r.Context())
	if routeCtx != nil {
		pattern := routeCtx.RoutePattern()
		if pattern != "" {
			return pattern
		}
	}
	return r.URL.Path
}
