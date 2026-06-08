package service

import (
    "context"
    "testing"

    "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "github.com/Perlishnov/gotrainingproject/internal/models"
    "github.com/Perlishnov/gotrainingproject/mocks"
)

func TestCreateUser_Success(t *testing.T) {
    mockDAO := new(mocks.UserDAO)
    logger := logrus.New()
    logger.SetLevel(logrus.FatalLevel) // silence logs in tests
    svc := NewUserService(mockDAO, logger)

    req := &models.CreateUserRequest{
        Name:     "John",
        Email:    "john@example.com",
        Password: "secret123",
    }

    mockDAO.On("GetByEmail", mock.Anything, req.Email).Return(nil, nil)
    mockDAO.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).Return(nil).Run(func(args mock.Arguments) {
        u := args.Get(1).(*models.User)
        u.ID = 1
    })

    user, err := svc.CreateUser(context.Background(), req)

    assert.NoError(t, err)
    assert.Equal(t, int64(1), user.ID)
    assert.Equal(t, "John", user.Name)
    mockDAO.AssertExpectations(t)
}

func TestCreateUser_EmailExists(t *testing.T) {
    mockDAO := new(mocks.UserDAO)
    logger := logrus.New()
    logger.SetLevel(logrus.FatalLevel)
    svc := NewUserService(mockDAO, logger)

    req := &models.CreateUserRequest{
        Email: "existing@example.com",
    }
    mockDAO.On("GetByEmail", mock.Anything, req.Email).Return(&models.User{ID: 99}, nil)

    _, err := svc.CreateUser(context.Background(), req)

    assert.Error(t, err)
    assert.Contains(t, err.Error(), "already exists")
    mockDAO.AssertExpectations(t)
}