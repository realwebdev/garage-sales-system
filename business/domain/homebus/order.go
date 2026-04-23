package homebus

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByID, order.ASC)

// Set of fields that the results can be ordered by.
const (
	OrderByID     = "a"
	OrderByType   = "b"
	OrderByUserID = "c"
)
