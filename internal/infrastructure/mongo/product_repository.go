package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoProductRepository struct {
	collection *mongo.Collection
}

func NewMongoProductRepository(db *mongo.Database) product.Repository {
	return &mongoProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *mongoProductRepository) AddProduct(ctx context.Context, p *product.Product) error {
	p.CreatedAt = time.Now().Format(time.RFC3339)
	_, err := r.collection.InsertOne(ctx, p)
	return err
}

func (r *mongoProductRepository) GetProductByID(ctx context.Context, id string) (*product.Product, error) {
	var p product.Product
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *mongoProductRepository) ListProductsByReseller(ctx context.Context, resellerID string, page, limit int) ([]*product.Product, error) {
	var products []*product.Product
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	resellerObjectID, err := primitive.ObjectIDFromHex(resellerID)
	if err != nil {
		return nil, fmt.Errorf("invalid reseller ID: %w", err)
	}

	cursor, err := r.collection.Find(ctx, bson.M{"reseller_id": resellerObjectID}, opts)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p product.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}
		products = append(products, &p)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return products, nil
}

func (r *mongoProductRepository) ListAvailableProducts(ctx context.Context, page, limit int) ([]*product.Product, error) {
	var products []*product.Product
	skip := (page - 1) * limit
	opts := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("database query failed: %w", err)
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p product.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, fmt.Errorf("failed to decode product: %w", err)
		}
		products = append(products, &p)
	}

	if err := cursor.Err(); err != nil {
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return products, nil
}

func (r *mongoProductRepository) DeleteProduct(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

func (r *mongoProductRepository) UpdateProduct(ctx context.Context, id string, updates map[string]interface{}) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updates})
	return err
}
func (r *mongoProductRepository) GetProductsByBundleID(ctx context.Context, bundleID string) ([]*product.Product, error) {
	var products []*product.Product

	cursor, err := r.collection.Find(ctx, bson.M{"bundle_id": bundleID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var p product.Product
		if err := cursor.Decode(&p); err != nil {
			return nil, err
		}
		products = append(products, &p)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
