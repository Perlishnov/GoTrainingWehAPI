package service

import (
	"context"
	"errors"
	"testing"

	"github.com/Perlishnov/gotrainingproject/internal/models"
	"github.com/Perlishnov/gotrainingproject/mocks"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_CreateUser(t *testing.T) {
	// Silence logger during tests
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	type args struct {
		req *models.CreateUserRequest
	}
	tests := []struct {
		name      string
		setupMock func(mockDAO *mocks.UserDAO)
		args      args
		wantUser  func(user *models.User) bool // predicate to validate result
		wantErr   bool
		errContain string
	}{
		{
			name: "success - user created",
			setupMock: func(mockDAO *mocks.UserDAO) {
				mockDAO.On("GetByEmail", mock.Anything, "john@example.com").
					Return(nil, nil)
				mockDAO.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(nil).
					Run(func(args mock.Arguments) {
						u := args.Get(1).(*models.User)
						u.ID = "mock-id-123" // simulate ID set by DAO
					})
			},
			args: args{
				req: &models.CreateUserRequest{
					Name:     "John Doe",
					Email:    "john@example.com",
					Password: "secret123",
				},
			},
			wantUser: func(u *models.User) bool {
				return u.Name == "John Doe" &&
					u.Email == "john@example.com" &&
					u.Role == "user" &&
					u.ID == "mock-id-123" &&
					u.Password == "" // password cleared
			},
			wantErr: false,
		},
		{
			name: "fail - email already exists",
			setupMock: func(mockDAO *mocks.UserDAO) {
				mockDAO.On("GetByEmail", mock.Anything, "existing@example.com").
					Return(&models.User{ID: "some-id", Email: "existing@example.com"}, nil)
			},
			args: args{
				req: &models.CreateUserRequest{
					Name:     "Jane",
					Email:    "existing@example.com",
					Password: "pass",
				},
			},
			wantErr:   true,
			errContain: "already exists",
		},
		{
			name: "fail - role invalid",
			setupMock: func(mockDAO *mocks.UserDAO) {
				mockDAO.On("GetByEmail", mock.Anything, "invalid@example.com").
					Return(nil, nil)
			},
			args: args{
				req: &models.CreateUserRequest{
					Name:     "Bad Role",
					Email:    "invalid@example.com",
					Password: "pass",
					Role:     "superuser",
				},
			},
			wantErr:   true,
            errContain: "invalid role: must be 'user' or 'admin'", // updated
		},
		{
			name: "success - role admin allowed",
			setupMock: func(mockDAO *mocks.UserDAO) {
				mockDAO.On("GetByEmail", mock.Anything, "admin@example.com").
					Return(nil, nil)
				mockDAO.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(nil).
					Run(func(args mock.Arguments) {
						u := args.Get(1).(*models.User)
						u.ID = "admin-id"
					})
			},
			args: args{
				req: &models.CreateUserRequest{
					Name:     "Admin User",
					Email:    "admin@example.com",
					Password: "adminpass",
					Role:     "admin",
				},
			},
			wantUser: func(u *models.User) bool {
				return u.Role == "admin"
			},
			wantErr: false,
		},
		{
			name: "fail - GetByEmail database error",
			setupMock: func(mockDAO *mocks.UserDAO) {
				mockDAO.On("GetByEmail", mock.Anything, "db@error.com").
					Return(nil, errors.New("connection refused"))
			},
			args: args{
				req: &models.CreateUserRequest{
					Name:     "DB Fail",
					Email:    "db@error.com",
					Password: "pass",
				},
			},
			wantErr:   true,
			errContain: "check existing user",
		},
		{
			name: "fail - Create database error",
			setupMock: func(mockDAO *mocks.UserDAO) {
				mockDAO.On("GetByEmail", mock.Anything, "create@fail.com").
					Return(nil, nil)
				mockDAO.On("Create", mock.Anything, mock.AnythingOfType("*models.User")).
					Return(errors.New("duplicate key error"))
			},
			args: args{
				req: &models.CreateUserRequest{
					Name:     "Create Fail",
					Email:    "create@fail.com",
					Password: "pass",
				},
			},
			wantErr:   true,
			errContain: "create user in DB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDAO := new(mocks.UserDAO)
			tt.setupMock(mockDAO)
			svc := NewUserService(mockDAO, logger)

			got, err := svc.CreateUser(context.Background(), tt.args.req)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContain != "" {
					assert.Contains(t, err.Error(), tt.errContain)
				}
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				if tt.wantUser != nil {
					assert.True(t, tt.wantUser(got), "user did not meet expectations")
				}
			}
			mockDAO.AssertExpectations(t)
		})
	}
}
func TestUserService_GetUserByID(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	tests := []struct {
		name      string
		userID    string
		setupMock func(mockDAO *mocks.UserDAO)
		wantUser  *models.User
		wantErr   bool
	}{
		{
			name:   "success - user found",
			userID: "123",
			setupMock: func(mockDAO *mocks.UserDAO) {
				mockDAO.On("GetByID", mock.Anything, "123").
					Return(&models.User{ID: "123", Name: "Found"}, nil)
			},
			wantUser: &models.User{ID: "123", Name: "Found", Password: ""},
			wantErr:  false,
		},
		{
			name:   "not found",
			userID: "999",
			setupMock: func(mockDAO *mocks.UserDAO) {
				mockDAO.On("GetByID", mock.Anything, "999").
					Return(nil, nil)
			},
			wantUser: nil,
			wantErr:  true,
		},
		{
			name:   "database error",
			userID: "456",
			setupMock: func(mockDAO *mocks.UserDAO) {
				mockDAO.On("GetByID", mock.Anything, "456").
					Return(nil, errors.New("db error"))
			},
			wantUser: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDAO := new(mocks.UserDAO)
			tt.setupMock(mockDAO)
			svc := NewUserService(mockDAO, logger)

			got, err := svc.GetUserByID(context.Background(), tt.userID)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUser, got)
			}
			mockDAO.AssertExpectations(t)
		})
	}
}

func TestUserService_ValidateUserAccess(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)
	svc := NewUserService(nil, logger) // DAO not needed for this method

	tests := []struct {
		name       string
		userID     string
		targetID   string
		role       string
		shouldAllow bool
	}{
		{"admin can access any", "1", "2", "admin", true},
		{"user accessing self", "1", "1", "user", true},
		{"user accessing other", "1", "2", "user", false},
		{"empty role defaults to user", "1", "2", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			allowed := svc.ValidateUserAccess(tt.userID, tt.targetID, tt.role)
			assert.Equal(t, tt.shouldAllow, allowed)
		})
	}
}