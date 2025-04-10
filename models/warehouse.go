package models

type WarehouseItemResponse struct {
	ItemID     string `json:"item_id"`
	BundleID   string `json:"bundle_id"`
	ResellerID string `json:"reseller_id"`
	Status     string `json:"status"` // unpacked, listed, skipped, damaged
}
