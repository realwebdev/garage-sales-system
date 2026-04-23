package vproductbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/business/types/money"
	"github.com/realwebdev/garage-sales-system/business/types/name"
	"github.com/realwebdev/garage-sales-system/business/types/quantity"
)

// PRoduct represents an individual product with extended information.
type Product struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        name.Name
	UserName    name.Name
	Cost        money.Money
	Quanitity   quantity.Quantity
	DateCreated time.Time
	DateUpdated time.Time
}
