package product

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          string             `bson:"_id" json:"id"`
	ResellerID  primitive.ObjectID `bson:"reseller_id" json:"reseller_id"`
	SupplierID  string             `bson:"supplier_id" json:"supplier_id"`
	BundleID    string             `bson:"bundle_id" json:"bundle_id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Size        string             `json:"size"`
	Type        string             `json:"type"`
	Grade       string             `json:"grade"`
	Price       float64            `json:"price"`
	Status      string             `json:"status"`
	ImageURL    string             `json:"image_url"`
	CreatedAt   string             `json:"created_at"`
	Rating      float64            `bson:"rating" json:"rating"`
}

func (p *Product) GenerateID() string {
	return primitive.NewObjectID().Hex()
}
