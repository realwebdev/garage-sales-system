package productbus

import (
	"time"

	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/business/types/money"
	"github.com/realwebdev/garage-sales-system/business/types/name"
	"github.com/realwebdev/garage-sales-system/business/types/quantity"
)

// Product represents an individual product.
type Product struct {
	ID uuid.UUID
	UserID uuid.UUID
	Name name.Name
	Cost money.Money
	Quantity quantity.Quantity
	DateCreated time.Time
	DateUpdated time.Time
}

// NewProduct is what we require from clients when adding a product.
type NewProduct struct {
	UserID uuid.UUID
	Name name.Name
	Cost money.Money
	Quanitity quantity.Quantity
}

// UpdateProduct defines what information may be provided to modify an 
// existing Product. All fields are optional so clients can send just the
// fields they want changed. It uses pointer fields that was provided as
// explicitly blank. Normally we do not want to use pointers to basic types
// but we make exception around marshalling and unmarshalling.
type UpdateProduct struct {
	Name *name.Name
	Cost *money.Money
	Quantity quantity.Quantity
}