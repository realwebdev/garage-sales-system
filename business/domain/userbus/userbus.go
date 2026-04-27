// Package userbus provides business access to user domain
package userbus

import (
	"context"
	"errors"
	"net/mail"

	"github.com/ardanlabs/service/business/sdk/delegate"
	"github.com/ardanlabs/service/business/sdk/page"
	"github.com/ardanlabs/service/business/sdk/sqldb"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/business/sdk/order"
	"github.com/vertica/vertica-sql-go/logger"
)

// Set of error variables for CRUD operations.

var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer interface declares the behavior this package needs to
// persist and retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollBacker) (Storer, error)
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// ExtBusiness interface provides support for extensions that wrap extra functionality around
// the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Create(ctx context.Context, actorID uuid.UUID, nu NewUser) (User, error)
	Update(ctx context.Context, actorID uuid.UUID, usr User, uu UpdateUser) (User, error)
	Delete(ctx context.Context, actorID uuid.UUID, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.page) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
	Authenticate(ctx context.Context, email mail.Address, password string) (User, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for user Access.
type Business struct {
	log       *logger.Logger // TODO: Change the import as it will be a custome logger
	storer    Storer
	delegate  *delegate.Delegate // TODO: Change the import as it will be a custome logger
	extension []Extension
}

// NewBusiness constructs a user business API for use
func NewBusiness(log *logger.Logger, delegate *delegate.Delegate, storer Storer, extensions ...Extension) ExtBusiness {
	b := ExtBusiness(&Business{
		log:       log,
		delegate:  delegate,
		storer:    storer,
		extension: extensions,
	})

	for i := len(extensions) - 1; i >= 0; i-- {
		ext := extensions[i]
		if ext != nil {
			b = ext
		}
	}

	return b
}

// NewWithTx constructs a new business value that will use
// the specified transaction in any store related calls.
func (b *Business) NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error) {
	storer, err := b.storer.NewWithTx(tx)
	if err != nil {
		return nil, err
	}

	nb := NewBusiness(b.log, b.delegate, storer, b.extension...)

	return nb, nil
}
