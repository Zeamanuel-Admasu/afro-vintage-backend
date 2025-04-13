package mongo

import (
    "context"
    "github.com/Zeamanuel-Admasu/afro-vintage-backend/internal/domain/bundleorder"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
)

type BundleOrderRepository struct {
    collection *mongo.Collection
}

func NewBundleOrderRepository(db *mongo.Database) *BundleOrderRepository {
    return &BundleOrderRepository{
        collection: db.Collection("bundle_orders"),
    }
}

func (r *BundleOrderRepository) CreateOrder(ctx context.Context, order *bundleorder.BundleOrder) error {
    _, err := r.collection.InsertOne(ctx, order)
    return err
}

func (r *BundleOrderRepository) GetOrdersByBuyerID(ctx context.Context, buyerID string) ([]*bundleorder.BundleOrder, error) {
    var orders []*bundleorder.BundleOrder
    cursor, err := r.collection.Find(ctx, bson.M{"buyer_id": buyerID})
    if err != nil {
        return nil, err
    }
    if err = cursor.All(ctx, &orders); err != nil {
        return nil, err
    }
    return orders, nil
}

func (r *BundleOrderRepository) GetOrderByBundleID(ctx context.Context, bundleID string) (*bundleorder.BundleOrder, error) {
    var order bundleorder.BundleOrder
    err := r.collection.FindOne(ctx, bson.M{"bundle_id": bundleID}).Decode(&order)
    if err != nil {
        return nil, err
    }
    return &order, nil
}

func (r *BundleOrderRepository) GetOrdersBySellerID(ctx context.Context, sellerID string) ([]*bundleorder.BundleOrder, error) {
    var orders []*bundleorder.BundleOrder
    cursor, err := r.collection.Find(ctx, bson.M{"seller_id": sellerID})
    if err != nil {
        return nil, err
    }
    if err = cursor.All(ctx, &orders); err != nil {
        return nil, err
    }
    return orders, nil
}