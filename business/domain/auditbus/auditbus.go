// Package auditbus provides a business logic layer for handling audit events.
package auditbus

import (
	"context"

	"github.com/ardanlabs/service/business/sdk/page"
	"github.com/realwebdev/garage-sales-system/business/sdk/order"
	"github.com/realwebdev/garage-sales-system/foundation/logger"
)

type Storer interface {
	Create(ctx context.Context, audit Audit) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Audit, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}

// ExtBusiness interface provides support for extensions that wrap extra functionality
// around the core business logic.
type ExtBusiness interface {
	Create(ctx context.Context, na NewAudit) (Audit, error)
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Audit, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for audit access.
type Business struct {
	log    *logger.Logger
	strore Storer
}

// // NewBusiness constructs a audit business API for use.
// func NewBusiness(log *logger.Logger, storer Storer, extensions ...ExtBusiness)(ExtBusiness){
// 	b := ExtBusiness(&Business{
// 		log: log,
// 		strore: storer,
// 	})

// }
