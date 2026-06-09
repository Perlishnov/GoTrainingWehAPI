package middleware

import (
    "context"
    "net/http"
    "strings"

    "github.com/sirupsen/logrus"
	"github.com/Perlishnov/gotrainingproject/internal/utils"
)

type AuthMiddleware struct {
    jwtUtil *utils.JWTUtil
    logger  *logrus.Logger
}

func NewAuthMiddleware(jwtUtil *utils.JWTUtil, logger *logrus.Logger) *AuthMiddleware {
    return &AuthMiddleware{jwtUtil: jwtUtil, logger: logger}
}

func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler  {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			m.logger.Warn("missing auth header")
			http.Error(w, "authorization header required", http.StatusUnauthorized)
            return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            m.logger.Warn("invalid auth header format")
            http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
            return
        }
        claims, err := m.jwtUtil.ValidateToken(parts[1])
        if err != nil {
            m.logger.WithError(err).Warn("invalid token")
            http.Error(w, "invalid or expired token", http.StatusUnauthorized)
            return
        }
		ctx := r.Context()
		ctx = context.WithValue(ctx, "userID", claims.UserID)
		ctx = context.WithValue(ctx, "userRole", claims.Role)
		ctx = context.WithValue(ctx, "userEmail", claims.Email)
		m.logger.WithField("user_id", claims.UserID).Debug("authenticated request")
		next.ServeHTTP(w,r.WithContext(ctx))
	})
}