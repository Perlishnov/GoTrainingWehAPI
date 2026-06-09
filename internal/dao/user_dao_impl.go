package dao

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Perlishnov/gotrainingproject/internal/models"
)


type UserDAOImpl struct{
	db * sql.DB
}

func NewUserDAO( db *sql.DB) UserDAO  {
	return &UserDAOImpl{db: db}
}

func (d *UserDAOImpl) Create(ctx context.Context,user *models.User ) error {
	query := `INSERT INTO users (name, email, password, role, created_at, updated_at)
	VALUES (?,?,?,?,?,?)
	`
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	res, err := d.db.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.Role, now, now)
	if err != nil {
		return fmt.Errorf("Failed to create user: %w", err)
	}
	id, _ := res.LastInsertId()
	user.ID = id
	return nil
}

func (d *UserDAOImpl) GetByID(ctx context.Context, id int64) (*models.User, error)  {
	query := `SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE id =?`

	var u models.User
	err := d.db.QueryRowContext(ctx, query, id).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user %w", err)
	}
	return &u, nil
}

func (d *UserDAOImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, name, email, password, role, created_at, updated_at FROM users WHERE email =?`

	var u models.User
	err := d.db.QueryRowContext(ctx, query, email).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user %w", err)
	}
	return &u, nil
}

func (d *UserDAOImpl) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
    query := `SELECT id, name, email, password, role, created_at, updated_at FROM users LIMIT ? OFFSET ?`
    rows, err := d.db.QueryContext(ctx, query, limit, offset)
    if err != nil {
        return nil, fmt.Errorf("failed to get users: %w", err)
    }
    defer rows.Close()
    var users []models.User
    for rows.Next() {
        var u models.User
        if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
            return nil, fmt.Errorf("failed to scan user: %w", err)
        }
        users = append(users, u)
    }
    return users, nil
}

func (d *UserDAOImpl) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET name=?, email=?, role=?, updated_at=? WHERE id=?`
	user.UpdatedAt = time.Now()
	_, err := d.db.ExecContext(ctx,query,user.Name, user.Email, user.Role, user.UpdatedAt, user.ID)

	return err
}

func (d *UserDAOImpl) Delete(ctx context.Context, id int64) error {
	_, err := d.db.ExecContext(ctx, `DELETE FROM users WHERE id=?`,id)
	return err
}
