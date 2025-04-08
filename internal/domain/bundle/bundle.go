package bundle

type SortingLevel string

const (
	Sorted     SortingLevel = "sorted"
	SemiSorted SortingLevel = "semi_sorted"
	Unsorted   SortingLevel = "unsorted"
)

type Bundle struct {
	ID                 string
	SupplierID         string
	Title              string
	Description        string
	SampleImage        string
	Quantity           int
	Grade              string
	SortingLevel       SortingLevel
	EstimatedBreakdown map[string]int `bson:"estimatedBreakdown,omitempty"` // Only for semi
	Type               string         `bson:"type,omitempty"`               // For sorted/semi
	Price              float64
	Status             string
	CreatedAt          string
}
