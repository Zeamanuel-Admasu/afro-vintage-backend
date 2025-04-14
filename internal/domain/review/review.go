package review

type Review struct {
	ID        string
	OrderID   string
	ProductID string
	UserID    string
	Rating    int
	Comment   string
	CreatedAt string
}