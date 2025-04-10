package models

type OrderRequest struct {
	ProductID string `json:"product_id"`
}

type OrderResponse struct {
	ID      string          `json:"id"`
	BuyerID string          `json:"buyer_id"`
	Product ProductResponse `json:"product"`
	Status  string          `json:"status"`
}
