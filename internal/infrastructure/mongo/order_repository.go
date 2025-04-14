package mongo

import (
	"context"

	"github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/order"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoOrderRepository struct {
	collection *mongo.Collection
}

func (r *mongoOrderRepository) UpdateOrderStatus(ctx context.Context, orderID string, status order.OrderStatus) error {
	filter := bson.M{"order_id": orderID}
	update := bson.M{"$set": bson.M{"status": status}}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *mongoOrderRepository) DeleteOrder(ctx context.Context, orderID string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": orderID})
	return err
}

func (r *mongoOrderRepository) CreateOrder(ctx context.Context, o *order.Order) error {
	_, err := r.collection.InsertOne(ctx, o)
	return err
}

func NewMongoOrderRepository(db *mongo.Database) order.Repository {
	return &mongoOrderRepository{
		collection: db.Collection("orders"),
	}
}

func (r *mongoOrderRepository) GetOrdersByConsumer(ctx context.Context, consumerID string) ([]*order.Order, error) {
	var orders []*order.Order
	cursor, err := r.collection.Find(ctx, bson.M{"consumer_id": consumerID}) // Query by consumer_id
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var o order.Order
		if err := cursor.Decode(&o); err != nil {
			return nil, err
		}
		orders = append(orders, &o)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *mongoOrderRepository) GetOrderByID(ctx context.Context, orderID string) (*order.Order, error) {
	var o order.Order
	err := r.collection.FindOne(ctx, bson.M{"_id": orderID}).Decode(&o)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &o, nil
}
