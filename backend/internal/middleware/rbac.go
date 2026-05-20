package middleware

import (
	"context"
	"net/http"
	"strconv"

	"gorm.io/gorm"

	"github.com/yourorg/maintenance/internal/domain"
	"github.com/yourorg/maintenance/internal/response"
)

// contextKey avoids key collisions in request context.
type contextKey string

const (
	ContextKeyUser contextKey = "authenticated_user"
)

// Authenticate reads the X-User-ID header, queries the users table,
// and stores the *domain.User in the request context.
// Returns 401 if the header is missing / invalid, 404 if user not found.
func Authenticate(db *gorm.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawID := r.Header.Get("X-User-ID")
			if rawID == "" {
				response.Unauthorized(w, "X-User-ID header is required")
				return
			}

			userID, err := strconv.ParseUint(rawID, 10, 64)
			if err != nil {
				response.Unauthorized(w, "X-User-ID must be a valid integer")
				return
			}

			var user domain.User
			if err := db.First(&user, userID).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					response.Unauthorized(w, "user not found")
					return
				}
				response.InternalError(w, "failed to authenticate user")
				return
			}

			// Inject user into context for downstream handlers
			ctx := context.WithValue(r.Context(), ContextKeyUser, &user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireRole returns a middleware that allows only users with one of the given roles.
// Must be chained AFTER Authenticate.
func RequireRole(roles ...domain.UserRole) func(http.Handler) http.Handler {
	allowed := make(map[domain.UserRole]struct{}, len(roles))
	for _, r := range roles {
		allowed[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := UserFromContext(r.Context())
			if !ok {
				response.Unauthorized(w, "unauthenticated")
				return
			}

			if _, permitted := allowed[user.Role]; !permitted {
				response.Forbidden(w, "you do not have permission to access this resource")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// UserFromContext is a helper to retrieve the authenticated user from context.
func UserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(ContextKeyUser).(*domain.User)
	return user, ok
}
