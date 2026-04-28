// Package userbus provides business access to user domain.
package userbus

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/ardanlabs/service/business/sdk/delegate"
	"github.com/ardanlabs/service/business/sdk/page"
	"github.com/ardanlabs/service/business/sdk/sqldb"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/business/sdk/order"
	"github.com/realwebdev/garage-sales-system/business/types/role"
	"golang.org/x/crypto/bcrypt"
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
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, usr User) error
	Update(ctx context.Context, usr User) error
	Delete(ctx context.Context, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
}

// ExtBusiness interface provides support for extensions that wrap extra
// functionality around the core business logic.
type ExtBusiness interface {
	NewWithTx(tx sqldb.CommitRollbacker) (ExtBusiness, error)
	Create(ctx context.Context, actorID uuid.UUID, nu NewUser) (User, error)
	Update(ctx context.Context, actorID uuid.UUID, usr User, uu UpdateUser) (User, error)
	Delete(ctx context.Context, actorID uuid.UUID, usr User) error
	Query(ctx context.Context, filter QueryFilter, orderBy order.By, page page.Page) ([]User, error)
	Count(ctx context.Context, filter QueryFilter) (int, error)
	QueryByID(ctx context.Context, userID uuid.UUID) (User, error)
	QueryByEmail(ctx context.Context, email mail.Address) (User, error)
	Authenticate(ctx context.Context, email mail.Address, password string) (User, error)
}

// Extension is a function that wraps a new layer of business logic
// around the existing business logic.
type Extension func(ExtBusiness) ExtBusiness

// Business manages the set of APIs for user access.
type Business struct {
	log       *logger.Logger
	storer    Storer
	delegate  *delegate.Delegate
	extension []Extension
}

// NewBusiness constructs a user business API for use.
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
			b = ext(b)
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

	return NewBusiness(b.log, b.delegate, storer, b.extension...), nil
}

// Create adds a new user to the system.
func (b *Business) Create(ctx context.Context, actorID uuid.UUID, nu NewUser) (User, error) {
	_ = actorID

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password.String()), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("generateFromPassword: %w", err)
	}

	now := time.Now()
	usr := User{
		ID:           uuid.New(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        cloneRoles(nu.Roles),
		PasswordHash: hash,
		Department:   nu.Department,
		Enabled:      true,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := b.storer.Create(ctx, usr); err != nil {
		return User{}, err
	}

	return usr, nil
}

// Update modifies information about a user.
func (b *Business) Update(ctx context.Context, actorID uuid.UUID, usr User, uu UpdateUser) (User, error) {
	_ = actorID

	if uu.Name != nil {
		usr.Name = *uu.Name
	}

	if uu.Email != nil {
		usr.Email = *uu.Email
	}

	if uu.Roles != nil {
		usr.Roles = cloneRoles(uu.Roles)
	}

	if uu.Department != nil {
		usr.Department = *uu.Department
	}

	if uu.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(uu.Password.String()), bcrypt.DefaultCost)
		if err != nil {
			return User{}, fmt.Errorf("generateFromPassword: %w", err)
		}

		usr.PasswordHash = hash
	}

	if uu.Enabled != nil {
		usr.Enabled = *uu.Enabled
	}

	usr.DateUpdated = time.Now()

	if err := b.storer.Update(ctx, usr); err != nil {
		return User{}, err
	}

	return usr, nil
}

// Delete removes a user from the system.
func (b *Business) Delete(ctx context.Context, actorID uuid.UUID, usr User) error {
	_ = actorID
	return b.storer.Delete(ctx, usr)
}

// Query retrieves users based on the supplied filter.
func (b *Business) Query(ctx context.Context, filter QueryFilter, orderBy order.By, pg page.Page) ([]User, error) {
	return b.storer.Query(ctx, filter, orderBy, pg)
}

// Count returns the number of matching users.
func (b *Business) Count(ctx context.Context, filter QueryFilter) (int, error) {
	return b.storer.Count(ctx, filter)
}

// QueryByID retrieves a user by ID.
func (b *Business) QueryByID(ctx context.Context, userID uuid.UUID) (User, error) {
	return b.storer.QueryByID(ctx, userID)
}

// QueryByEmail retrieves a user by email.
func (b *Business) QueryByEmail(ctx context.Context, email mail.Address) (User, error) {
	return b.storer.QueryByEmail(ctx, email)
}

// Authenticate verifies a user's credentials.
func (b *Business) Authenticate(ctx context.Context, email mail.Address, password string) (User, error) {
	usr, err := b.storer.QueryByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return User{}, ErrAuthenticationFailure
		}

		return User{}, err
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return User{}, ErrAuthenticationFailure
	}

	return usr, nil
}

func cloneRoles(roles []role.Role) []role.Role {
	if roles == nil {
		return nil
	}

	cloned := make([]role.Role, len(roles))
	copy(cloned, roles)

	return cloned
}
