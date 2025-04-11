package product

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          string             `bson:"_id"`
	ResellerID  primitive.ObjectID `bson:"reseller_id"`
	Title       string
	Description string
	Size        string
	Type        string
	Grade       string
	Price       float64
	ImageURL    string
	CreatedAt   string
}

func (p *Product) GenerateID() string {
	return primitive.NewObjectID().Hex()
}
