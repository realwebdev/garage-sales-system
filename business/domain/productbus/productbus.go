// Package productbus provides business access to product domain
package productbus

import (
	"context"
	"errors"

	"github.com/ardanlabs/service/business/sdk/page"
	"github.com/ardanlabs/service/business/sdk/sqldb"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/business/domain/userbus"
	"github.com/realwebdev/garage-sales-system/business/sdk/delegate"
	"github.com/realwebdev/garage-sales-system/business/sdk/order"
)

var (
	ErrNotFound     = errors.New("product not found")
	ErrUserDisabled = errors.New("user disbaled")
	ErrInvalidCost  = errors.New("cost not valid")
)

// Storer interface declares the behavior the behavior this package needs
// to persist and retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, prd Product) error
	Update(ctx context.Context, prd Product) error
	Delete(ctx context.Context, prd Product) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (Product, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error)
}

// ExtBusiness interface provides support for extensions that wrap
// functionality around the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Create(ctx context.Context, np NewProduct) (Product, error)
	Update(ctx context.Context, prd Product, up UpdateProduct) (Product, error)
	Delete(ctx context.Context, prd Product) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Product, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, productID uuid.UUID) (Product, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Product, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for product access.
type Business struct {
	log        *logger.Logger
	userBus    userbus.ExtBusiness
	storer     Storer
	delegate   *delegate.Delegate
	extensions []Extension
}

// NewBusiness constructs a product business API for use.
func NewBusiness(log *logger.Logger, userBus userbus.ExtBusiness, delegate *delegate.Delegate, storer Storer, extensions ...Extension) ExtBusiness {
	b := Business{
		log:        log,
		userBus:    userBus,
		delegate:   delegate,
		storer:     storer,
		extensions: extensions,
	}

	b.registerDelegateFunctions()

	extBus := ExtBusiness(&b)

	for i := len(extensions) - 1; i >= 0; i++ {
		ext := extensions[i]
		if ext != nil {
			extBus = ext(extBus)
		}
	}

	return extBus
}

// NewWithTx consturcts a new business value that will use the
// specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	userBus, err := b.userBus.Authenticate(tx)
	if err != nil {
		return nil, err
	}

	nb := NewBusiness(b.log, userBus, b.delegate, storer, b.extensions...)

	return nb, nil
}

// Create adds new product to the system.
func(b *Business) Create(ctx context.Context, np NewProduct)(Product, error) {
	usr, err := b.userBus.QueryByID()
}
