package homebus

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/business/domain/userbus"
	"github.com/realwebdev/garage-sales-system/business/sdk/delegate"
	"github.com/realwebdev/garage-sales-system/business/sdk/order"
	"github.com/realwebdev/garage-sales-system/business/sdk/page"
	"github.com/realwebdev/garage-sales-system/business/sdk/sqldb"
	"github.com/realwebdev/garage-sales-system/foundation/logger"
)

// set of error variables for CRUD operations.
var (
	ErrNotFound     = errors.New("home not found")
	ErrUserDisabled = errors.New("user disabled")
)

// Storer interface declares the behavior this package need to persis and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, hme Home) error
	Update(ctx context.Context, hme Home) error
	Delete(ctx context.Context, hme Home) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Home, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, homeID uuid.UUID) (Home, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Home, error)
}

// ExtBusiness interface provides support for extensions that wraps
// extra functionality around the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, nh NewHome) (Home, error)
	Update(ctx context.Context, hme Home, uh UpdateHome) (Home, error)
	Delete(ctx context.Context, hme Home) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]Home, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, homeID uuid.UUID) (Home, error)
	QueryByUserID(ctx context.Context, userID uuid.UUID) ([]Home, error)
}

// Extension is a function that wraps a new layer of business logic
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for home api access.
type Business struct {
	log        *logger.Logger
	userBus    userbus.ExtBusiness
	delegate   *delegate.Delegate
	storer     Storer
	extensions []Extension
}

// NewBusiness constructs a home business API for use.
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

	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			extBus = ext(extBus)
		}
	}

	return extBus

}

// NewWithTx constructs a new domain value that will
// use the specified transaction in any store related calls.
func(b *Business)NewWithTx(tx sqldb.CommitRollbacker)(ExtBusiness, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	userBus, err := b.userBus.NewWithTx(tx)
	if err != nil
}