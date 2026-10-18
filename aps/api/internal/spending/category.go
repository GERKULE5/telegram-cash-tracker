package spending

type Category string

const (
	CategoryFood          Category = "food"
	CategoryTransport     Category = "transport"
	CategoryHousing       Category = "housing"
	CategoryEntertainment Category = "entertainment"
	CategoryHealth        Category = "health"
	CategoryPurchases     Category = "purchases"
	CategoryOther         Category = "other"
)

func (c Category) IsValid() bool {
	switch c {
	case CategoryFood,
		CategoryTransport,
		CategoryHousing,
		CategoryEntertainment,
		CategoryHealth,
		CategoryPurchases,
		CategoryOther:
		return true
	default:
		return false
	}
}
