package order

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusShipped   OrderStatus = "shipped"
	StatusDelivered OrderStatus = "delivered"
)

type Order struct {
	ID         string
	ConsumerID string
	ProductIDs []string
	TotalPrice float64
	Status     OrderStatus
	CreatedAt  string
}
