package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
    collection *mongo.Collection
}

func (r *mongoUserRepository) ListUsersByRole(ctx context.Context, role user.Role) ([]*user.User, error) {
    var users []*user.User
    cursor, err := r.collection.Find(ctx, bson.M{"role": string(role)})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        var u user.User
        if err := cursor.Decode(&u); err != nil {
            return nil, err
        }
        users = append(users, &u)
    }

    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return users, nil
}

func NewMongoUserRepository(db *mongo.Database) user.Repository {
    return &mongoUserRepository{
        collection: db.Collection("users"),
    }
}

func (r *mongoUserRepository) CreateUser(ctx context.Context, u *user.User) error {
	u.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, u)
	return err
}

func (r *mongoUserRepository) DeleteUser(ctx context.Context, id string) error {
    _, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
    return err
}

func (r *mongoUserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
    var u user.User
    err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&u)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *mongoUserRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
    var u user.User
    err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&u)
    if err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *mongoUserRepository) FindUserByUsername(ctx context.Context, username string) (*user.User, error) {
    var u user.User
    err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&u)
    if err != nil {
        return nil, errors.New("user not found")
    }
    return &u, nil
}

func (r *mongoUserRepository) UpdateUser(ctx context.Context, id string, updatedData map[string]interface{}) error {
    _, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updatedData})
    return err
}