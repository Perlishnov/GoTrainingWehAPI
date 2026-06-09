package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Perlishnov/gotrainingproject/internal/models"
	"github.com/Perlishnov/gotrainingproject/internal/utils"
)

// UserDAOMongo implements the UserDAO interface using MongoDB.
type UserDAOMongo struct {
	collection *mongo.Collection
}

// NewUserDAOMongo creates a new MongoDB DAO instance.
// It ensures a unique index on the email field.
func NewUserDAOMongo(db *mongo.Database, logger *logrus.Logger) UserDAO {
	collection := db.Collection("users")

	// Create a unique index on the email field to prevent duplicates.
	_, err := collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {

		if !(utils.IsIndexAlreadyExistsError(err)) {
            logger.WithError(err).Fatal("failed to create unique index on users.email")
        }
	}

	return &UserDAOMongo{collection: collection}
}



// Create inserts a new user document into MongoDB.
// It generates a unique int64 ID for the user using a nanosecond timestamp.
func (d *UserDAOMongo) Create(ctx context.Context, user *models.User) error {
	// Generate a new ID using current timestamp (nanoseconds) for uniqueness.
	// This avoids needing a separate sequence or ObjectID conversion.
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	_, err := d.collection.InsertOne(ctx, user)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("user with email %s already exists", user.Email)
		}
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetByID retrieves a user by its numeric ID.
func (d *UserDAOMongo) GetByID(ctx context.Context, id string) (*models.User, error) {
	filter := bson.M{"_id": id}
	var user models.User
	err := d.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// GetByEmail retrieves a user by their email address.
func (d *UserDAOMongo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	filter := bson.M{"email": email}
	var user models.User
	err := d.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

// GetAll retrieves a paginated list of users.
func (d *UserDAOMongo) GetAll(ctx context.Context, limit, offset int) ([]models.User, error) {
	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(int64(limit))
	}
	if offset > 0 {
		findOptions.SetSkip(int64(offset))
	}

	cursor, err := d.collection.Find(ctx, bson.M{}, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %w", err)
	}
	return users, nil
}

// Update replaces an existing user document.
func (d *UserDAOMongo) Update(ctx context.Context, user *models.User) error {
	user.UpdatedAt = time.Now()
	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}
	result, err := d.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("user with id %s not found", user.ID)
	}
	return nil
}

// Delete removes a user document from MongoDB.
func (d *UserDAOMongo) Delete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	result, err := d.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("user with id %s not found", id)
	}
	return nil
}