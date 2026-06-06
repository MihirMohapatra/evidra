package transport

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/evidra/evidra/identity/domain"
	"github.com/evidra/evidra/identity/service"
)

type contextKey string

const (
	userKey    contextKey = "user"
	sessionKey contextKey = "session"
)

type Authenticator interface {
	ValidateSession(ctx context.Context, token string) (*domain.User, error)
	ValidateAPIKey(ctx context.Context, key string) (*domain.APIKey, error)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
		)
		next.ServeHTTP(w, r)
		slog.Info("response",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", time.Since(start),
		)
	})
}

func Authentication(auth Authenticator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			if header == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			if strings.HasPrefix(header, "Bearer ") {
				token := strings.TrimPrefix(header, "Bearer ")
				user, err := auth.ValidateSession(r.Context(), token)
				if err != nil {
					http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), userKey, user)
				ctx = service.CtxWithActor(ctx, user)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			if strings.HasPrefix(header, "ApiKey ") {
				key := strings.TrimPrefix(header, "ApiKey ")
				apiKey, err := auth.ValidateAPIKey(r.Context(), key)
				if err != nil {
					http.Error(w, `{"error":"invalid api key"}`, http.StatusUnauthorized)
					return
				}
				ctx := context.WithValue(r.Context(), sessionKey, apiKey)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			http.Error(w, `{"error":"unsupported authorization type"}`, http.StatusUnauthorized)
		})
	}
}

func RequirePermission(perm domain.Permission) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := r.Context().Value(userKey).(*domain.User)
			if !ok {
				http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
				return
			}
			if !user.Role.HasPermission(perm) {
				http.Error(w, `{"error":"forbidden"}`, http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func UserFromContext(ctx context.Context) *domain.User {
	user, _ := ctx.Value(userKey).(*domain.User)
	return user
}
