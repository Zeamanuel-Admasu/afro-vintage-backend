package bundleorder

import "time"

type BundleOrder struct {
    ID              string    `bson:"_id" json:"id"`
    BuyerID         string    `bson:"buyer_id" json:"buyer_id"` // Reseller's ID
    SellerID        string    `bson:"seller_id" json:"seller_id"` // Supplier's ID
    BundleID        string    `bson:"bundle_id" json:"bundle_id"`
    WarehouseStatus string    `bson:"warehouse_status" json:"warehouse_status"` // "unpacked", "completed"
    CreatedAt       time.Time `bson:"created_at" json:"created_at"`
}
