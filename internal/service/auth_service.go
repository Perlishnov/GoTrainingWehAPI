package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Perlishnov/gotrainingproject/internal/dao"
	"github.com/Perlishnov/gotrainingproject/internal/utils"
	"github.com/sirupsen/logrus"
)


type AuthService interface{
	Login(ctx context.Context, email, password string) (string, error)
	Logout(token string) error
}

type authService struct{
	userDAO dao.UserDAO
	jwtUtil *utils.JWTUtil
	logger *logrus.Logger
}

func  NewAuthService(userDAO dao.UserDAO, jwtUtil *utils.JWTUtil, logger *logrus.Logger) AuthService  {
	return &authService{
		userDAO: userDAO,
		jwtUtil: jwtUtil,
		logger: logger,
	}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.userDAO.GetByEmail(ctx, email)
	if err != nil {
		s.logger.WithError(err).Error("database error during login")
		return "", fmt.Errorf("database error: %w",err)
	}
	if user == nil || !utils.CheckPasswordHash(password, user.Password) {
		s.logger.WithField("email", email).Warn("failed login attempt")
		return "", errors.New("invalid credentials")
	}
	token, err := s.jwtUtil.GenerateToken(user.ID,user.Email,user.Role)
	if err != nil{
			s.logger.WithError(err).Error("token generation failed")
			return "", fmt.Errorf("failed to generate token %w",err)
		}
	s.logger.WithField("user_id", user.ID).Info("user logged in")
	return token, nil
}

func (s *authService) Logout(token string) error  {
	s.logger.Info("logout called")
	return nil
}