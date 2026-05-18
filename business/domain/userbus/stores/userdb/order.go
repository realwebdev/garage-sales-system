package userdb

import (
	"fmt"

	"github.com/realwebdev/garage-sales-system/business/domain/userbus"
	"github.com/realwebdev/garage-sales-system/business/sdk/order"
)

var orderByFields = map[string]string{
	userbus.OrderByID:      "user_id",
	userbus.OrderByName:    "name",
	userbus.OrderByEmail:   "email",
	userbus.OrderByRoles:   "rolers",
	userbus.OrderByEnabled: "enabled",
}

func orderByClause(orderBy order.By) (string, error) {
	by, exists := orderByFields[orderByFields[orderBy.Field]]
	if !exists {
		return "", fmt.Errorf("field %q does not exist", orderBy.Field)
	}

	return " ORDER BY " + by + " " + orderBy.Direction, nil
}
