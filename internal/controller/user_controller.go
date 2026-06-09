package controller

import (
	"net/http"
	"strconv"
	"strings"
    "github.com/gin-gonic/gin"
	"github.com/Perlishnov/gotrainingproject/internal/models"
	"github.com/Perlishnov/gotrainingproject/internal/service"
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
func (c *UserController) CreateUser(ctx *gin.Context) {
    var req models.CreateUserRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        c.logger.WithError(err).Warn("invalid create user request")
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
        return
    }

    user, err := c.userService.CreateUser(ctx.Request.Context(), &req)
    if err != nil {
        c.logger.WithError(err).WithField("email", req.Email).Error("create user failed")
        if strings.Contains(err.Error(), "already exists") {
            ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
            return
        }
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, user)
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
func (c *UserController) GetCurrentUser(ctx *gin.Context) {
    userID, exists := ctx.Get("userID")
    if !exists {
        ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
        return
    }
    user, err := c.userService.GetCurrentUser(ctx.Request.Context(), userID.(string))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, user)
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
func (c *UserController) GetUserByID(ctx *gin.Context) {
    targetID := ctx.Param("id")
    if targetID == "" {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
        return
    }
    userID, _ := ctx.Get("userID")
    userRole, _ := ctx.Get("userRole")

    if !c.userService.ValidateUserAccess(userID.(string), targetID, userRole.(string)) {
        ctx.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
        return
    }
    user, err := c.userService.GetUserByID(ctx.Request.Context(), targetID)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, user)
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
func (c *UserController) GetAllUsers(ctx *gin.Context) {
    userRole, exists := ctx.Get("userRole")
    if !exists || userRole.(string) != "admin" {
        ctx.JSON(http.StatusForbidden, gin.H{"error": "admin access required"})
        return
    }
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
    users, err := c.userService.GetAllUsers(ctx.Request.Context(), page, pageSize)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusOK, gin.H{
        "users":     users,
        "page":      page,
        "page_size": pageSize,
    })
}