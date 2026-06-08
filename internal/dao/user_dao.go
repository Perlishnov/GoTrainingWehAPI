package dao

import (
	"context"
	"github.com/Perlishnov/gotrainingproject/internal/models"
)


type UserDAO interface{
	Create(ctx context.Context,user *models.User ) error
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id int64) error
}
