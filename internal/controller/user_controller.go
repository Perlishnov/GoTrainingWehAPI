package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Perlishnov/gotrainingproject/internal/models"
	"github.com/Perlishnov/gotrainingproject/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type UserController struct {
    userService service.UserService
    logger      *logrus.Logger
}

func NewUserController(userService service.UserService, logger *logrus.Logger) *UserController {
    return &UserController{userService: userService, logger: logger}
}
// CreateUser godoc
// @Summary      Create a new user
// @Description  Registers a new user. Role can be "user" or "admin" (default "user").
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        request body models.CreateUserRequest true "User data"
// @Success      201 {object} models.User
// @Failure      400 {object} map[string]string "invalid input"
// @Failure      409 {object} map[string]string "email already exists"
// @Router       /users [post]
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req models.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        c.logger.WithError(err).Warn("invalid JSON")
        http.Error(w, "invalid request body", http.StatusBadRequest)
        return
    }

    user, err := c.userService.CreateUser(r.Context(), &req)
    if err != nil {
        // Log once with full context
        c.logger.WithError(err).WithField("email", req.Email).Error("create user failed")

        // Map error to HTTP status
        if strings.Contains(err.Error(), "already exists") {
            http.Error(w, err.Error(), http.StatusConflict)
        } else if strings.Contains(err.Error(), "invalid role") {
            http.Error(w, err.Error(), http.StatusBadRequest)
        } else {
            http.Error(w, "internal server error", http.StatusInternalServerError)
        }
        return
    }

    c.logger.WithField("user_id", user.ID).Info("user created")
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(user)
}
// GetCurrentUser godoc
// @Summary      Get current user profile
// @Description  Returns the profile of the authenticated user.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} models.User
// @Failure      401 {object} map[string]string "unauthorized"
// @Failure      404 {object} map[string]string "user not found"
// @Router       /users/me [get]
func (c *UserController) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
    userID := r.Context().Value("userID").(string)
    user, err := c.userService.GetCurrentUser(r.Context(), userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
// GetUserByID godoc
// @Summary      Get user by ID
// @Description  Returns a user's details. Requires admin role or matching user ID.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        id path string true "User ID"
// @Success      200 {object} models.User
// @Failure      401 {object} map[string]string "unauthorized"
// @Failure      403 {object} map[string]string "access denied"
// @Failure      404 {object} map[string]string "user not found"
// @Router       /users/{id} [get]
func (c *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
    targetID := chi.URLParam(r, "id")
    if targetID == "" {
        http.Error(w, "invalid user ID", http.StatusBadRequest)
        return
    }
    userID := r.Context().Value("userID").(string)
    userRole := r.Context().Value("userRole").(string)

    if !c.userService.ValidateUserAccess(userID, targetID, userRole) {
        http.Error(w, "access denied", http.StatusForbidden)
        return
    }
    user, err := c.userService.GetUserByID(r.Context(), targetID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(user)
}
// GetAllUsers godoc
// @Summary      List users
// @Description  Returns a paginated list of users. Admin only.
// @Tags         users
// @Produce      json
// @Security     BearerAuth
// @Param        page query int false "Page number (default 1)"
// @Param        page_size query int false "Items per page (default 10, max 100)"
// @Success      200 {object} map[string]interface{} "users, page, page_size"
// @Failure      401 {object} map[string]string "unauthorized"
// @Failure      403 {object} map[string]string "admin required"
// @Router       /users [get]
func (c *UserController) GetAllUsers(w http.ResponseWriter, r *http.Request) {
    // Get the role from the context
    userRole := r.Context().Value("userRole").(string)

    // Validate if user from the token users
    if userRole != "admin" {
        http.Error(w, "admin access required", http.StatusForbidden)
        return
    }

    page, _ := strconv.Atoi(r.URL.Query().Get("page"))
    pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
    users, err := c.userService.GetAllUsers(r.Context(), page, pageSize)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "users":     users,
        "page":      page,
        "page_size": pageSize,
    })
}