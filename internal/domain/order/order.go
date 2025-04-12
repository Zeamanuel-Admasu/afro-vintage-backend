package order

type OrderStatus string

const (
	Pending   OrderStatus = "pending"
	Shipped   OrderStatus = "shipped"
	Delivered OrderStatus = "delivered"
	Failed    OrderStatus = "failed"
)

type Order struct {
	ID         string
	ConsumerID string
	ProductIDs []string
	TotalPrice float64
	Status     OrderStatus
	CreatedAt  string
}
