package warehouse

type WarehouseItem struct {
	ID         string
	ResellerID string
	BundleID   string
	ProductID  string
	Status     string // listed, skipped, pending
	CreatedAt  string
}
