package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Perlishnov/gotrainingproject/internal/dao"
	"github.com/Perlishnov/gotrainingproject/internal/models"
	"github.com/Perlishnov/gotrainingproject/internal/utils"
	"github.com/sirupsen/logrus"
)

type UserService interface {
    CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error)
    GetUserByID(ctx context.Context, id string) (*models.User, error)
    GetCurrentUser(ctx context.Context, userID string) (*models.User, error)
    GetAllUsers(ctx context.Context, page, pageSize int) ([]models.User, error)
    ValidateUserAccess(userID, targetID string, role string) bool
}

type userService struct{
	userDAO dao.UserDAO
	logger *logrus.Logger
}

func NewUserService(userDAO dao.UserDAO, logger *logrus.Logger) UserService {
    return &userService{
        userDAO: userDAO,
        logger:  logger,
    }
}

func (s *userService) CreateUser(ctx context.Context, req *models.CreateUserRequest) (*models.User, error) {
    // 1. Check existing user
    existingUser, err := s.userDAO.GetByEmail(ctx, req.Email)
    if err != nil {
        return nil, fmt.Errorf("check existing user: %w", err)
    }
    if existingUser != nil {
        return nil, fmt.Errorf("user with email %s already exists", req.Email)
    }

    // 2. Validate role
    if req.Role != "" && req.Role != "user" && req.Role != "admin" {
        return nil, fmt.Errorf("invalid role: must be 'user' or 'admin', got %q", req.Role)
    }
    role := req.Role
    if role == "" {
        role = "user"
    }

    // 3. Hash password
    hashed, err := utils.HashPassword(req.Password)
    if err != nil {
        // Wrap the error – the controller will decide whether to log it.
        return nil, fmt.Errorf("hash password: %w", err)
    }

    // 4. Create user object
    user := &models.User{
        Name:     req.Name,
        Email:    req.Email,
        Password: hashed,
        Role:     role,
    }

    // 5. Persist
    if err := s.userDAO.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("create user in DB: %w", err)
    }

    // 6. Return without password
    user.Password = ""
    return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
    user, err := s.userDAO.GetByID(ctx, id)
    if err != nil {
        s.logger.WithError(err).WithField("user_id", id).Error("failed to get user")
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    if user == nil {
        return nil, errors.New("user not found")
    }
    user.Password = ""
    return user, nil
}

func (s *userService) GetCurrentUser(ctx context.Context, userID string) (*models.User, error) {
    return s.GetUserByID(ctx, userID)
}

func (s *userService) GetAllUsers(ctx context.Context, page, pageSize int) ([]models.User, error)  {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	users, err := s.userDAO.GetAll(ctx, pageSize,offset)

	if err != nil {
		s.logger.WithError(err).Error("failed to list users")
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	for i := range users{
		users[i].Password = ""
	}

	return users, nil
}

func (s *userService) ValidateUserAccess(userID, targetID string, role string) bool {
    return role == "admin" || userID == targetID
}