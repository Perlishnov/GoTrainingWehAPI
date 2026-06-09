package middleware

import (
	"net/http"
	"strings"

	"github.com/Perlishnov/gotrainingproject/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type AuthMiddleware struct {
    jwtUtil *utils.JWTUtil
    logger  *logrus.Logger
}

func NewAuthMiddleware(jwtUtil *utils.JWTUtil, logger *logrus.Logger) *AuthMiddleware {
    return &AuthMiddleware{jwtUtil: jwtUtil, logger: logger}
}

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            m.logger.Warn("missing auth header")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
            c.Abort()
            return
        }
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
            m.logger.Warn("invalid auth header format")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
            c.Abort()
            return
        }
        claims, err := m.jwtUtil.ValidateToken(parts[1])
        if err != nil {
            m.logger.WithError(err).Warn("invalid token")
            c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
            c.Abort()
            return
        }
        c.Set("userID", claims.UserID)
        c.Set("userRole", claims.Role)
        c.Set("userEmail", claims.Email)
        m.logger.WithField("user_id", claims.UserID).Debug("authenticated request")
        c.Next()
    }
}