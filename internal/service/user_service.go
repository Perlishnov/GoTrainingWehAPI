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
    GetUserByID(ctx context.Context, id int64) (*models.User, error)
    GetCurrentUser(ctx context.Context, userID int64) (*models.User, error)
    GetAllUsers(ctx context.Context, page, pageSize int) ([]models.User, error)
    ValidateUserAccess(userID, targetID int64, role string) bool
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
	existingUser, _ := s.userDAO.GetByEmail(ctx, req.Email)

	if existingUser != nil {
		s.logger.WithField("email", req.Email).Warn("user already exists")
		return nil, errors.New("user with this email already exists")
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logger.WithError(err).Error("password hashing failed")
		return nil, fmt.Errorf("failed to hash password")
	}
	role := req.Role
	if role == ""{
		role = "user"
	}

	user := &models.User{
		Name: req.Name,
		Email: req.Email,
		Password: hashed,
		Role: role,
	}

	if err := s.userDAO.Create(ctx, user); err != nil {
        s.logger.WithError(err).Error("failed to create user in DB")
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    user.Password = ""
    s.logger.WithField("user_id", user.ID).Info("user created")

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, id int64) (*models.User, error) {
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

func (s *userService) GetCurrentUser(ctx context.Context, userID int64) (*models.User, error) {
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

func (s *userService) ValidateUserAccess(userID, targetID int64, role string) bool {
    return role == "admin" || userID == targetID
}