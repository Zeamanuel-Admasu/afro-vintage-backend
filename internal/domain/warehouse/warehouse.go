package warehouse

type WarehouseItem struct {
	ID         string `bson:"_id"`
	ResellerID string `bson:"reseller_id"`
	BundleID   string `bson:"bundle_id"`
	ProductID  string `bson:"product_id"`
	Status     string `bson:"status"` // listed, skipped, pending
	CreatedAt  string `bson:"created_at"`
}
