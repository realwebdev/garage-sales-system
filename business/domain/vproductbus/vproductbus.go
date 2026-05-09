// Package vproductbus provides business access to view product domain.
package vproductbus

import (
	"context"
	"fmt"

	"github.com/ardanlabs/service/business/sdk/page"
	"github.com/realwebdev/garage-sales-system/business/sdk/order"
)

// Storer interface declares the behavior this package needs to persist
// and retrieve data
type Storer interface {
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}

// ExtBusiness interface provides support for extension that wraps extra functionality
// around the core business business logic
type ExtBusiness interface {
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for view product access
type Business struct {
	storer Storer
}

// NewBusiness contructs a vproduct business API for use.
func NewBusiness(storer Storer, extension ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		storer: storer,
	})

	for i := len(extension) - 1; i >= 0; i-- {
		ext := extension[i]
		if ext != nil {
			b = ext(b)
		}
	}

	return b
}

// Query retrieves list of existing products.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Product, error) {
	users, err := b.storer.Query(ctx, filter, orderBy, page)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return users, nil
}

// Count returns the total number of products.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}
