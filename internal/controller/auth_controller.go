package controller

import (
	"encoding/json"
	"net/http"

	"github.com/Perlishnov/gotrainingproject/internal/models"
	"github.com/Perlishnov/gotrainingproject/internal/service"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
    authService service.AuthService
    logger      *logrus.Logger
}

func NewAuthController(authService service.AuthService, logger *logrus.Logger) *AuthController {
    return &AuthController{authService: authService, logger: logger}
}


// Login godoc
// @Summary      User login
// @Description  Authenticates a user and returns a JWT token.
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body models.LoginRequest true "Login credentials"
// @Success      200 {object} map[string]string "token"
// @Failure      400 {object} map[string]string "invalid request"
// @Failure      401 {object} map[string]string "invalid credentials"
// @Router       /auth/login [post]
func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
    var req models.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        c.logger.WithError(err).Warn("invalid login request")
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }
    token, err := c.authService.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// Logout godoc
// @Summary      User logout
// @Description  Logs out the user (client must discard token).
// @Tags         auth
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} map[string]string "message"
// @Failure      500 {object} map[string]string "internal error"
// @Router       /auth/logout [post]
func (c *AuthController) Logout(w http.ResponseWriter, r *http.Request) {
    token := r.Header.Get("Authorization")
    if err := c.authService.Logout(token); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{"message": "logged out"})
}