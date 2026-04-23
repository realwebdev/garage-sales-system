package homebus

import (
	"time"

	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/business/types/home"
)

// QueryFilter holds the availble fields a query
// can be filtered on. we are using pointer semantics
// because the With api mutates the value.
type QueryFilter struct {
	ID *uuid.UUID
	UserID *uuid.UUID
	Type *home.Home
	StartCreatedDate *time.Time
	EndCreatedDate *time.Time
} 