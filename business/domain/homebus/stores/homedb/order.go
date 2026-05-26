package homedb

import (
	"fmt"

	"github.com/realwebdev/garage-sales-system/business/domain/homebus"
	"github.com/realwebdev/garage-sales-system/business/sdk/order"
)

var OrderByFields = map[string]string{
	homebus.OrderByID:     "home_id",
	homebus.OrderByType:   "type",
	homebus.OrderByUserID: "user_id",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := OrderByFields[orderBy.Field]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
