package transport

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/evidra/evidra/identity/domain"
	"github.com/evidra/evidra/identity/service"
)

func NewRouter(svc *service.IdentityService) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(Logging)
	r.Use(middleware.Recoverer)

	h := NewHandler(svc)

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Post("/auth/login", h.Login)
		r.Post("/auth/refresh", h.RefreshToken)

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(Authentication(svc))

			r.Post("/auth/logout", h.Logout)

			// Organizations
			r.Get("/organizations", h.ListOrganizations)
			r.Post("/organizations", h.CreateOrganization)
			r.Get("/organizations/{id}", h.GetOrganization)
			r.Put("/organizations/{id}", h.UpdateOrganization)
			r.Delete("/organizations/{id}", h.DeleteOrganization)

			// Users
			r.Get("/users", h.ListUsers)
			r.Post("/users", h.CreateUser)
			r.Get("/users/{id}", h.GetUser)
			r.Put("/users/{id}", h.UpdateUser)
			r.Delete("/users/{id}", h.DeleteUser)

			// API Keys
			r.Get("/api-keys", h.ListAPIKeys)
			r.Post("/api-keys", h.CreateAPIKey)
			r.Delete("/api-keys/{id}", h.RevokeAPIKey)
		})

		// Admin-only routes
		r.Group(func(r chi.Router) {
			r.Use(Authentication(svc))
			r.Use(RequirePermission(domain.PermissionDeleteOrganization))
			r.Delete("/organizations/{id}", h.DeleteOrganization)
		})
	})

	return r
}
