package usercache

import (
	"context"
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/business/domain/userbus"
	"github.com/realwebdev/garage-sales-system/business/sdk/order"
	"github.com/realwebdev/garage-sales-system/business/sdk/page"
	"github.com/realwebdev/garage-sales-system/business/sdk/sqldb"
	"github.com/realwebdev/garage-sales-system/foundation/logger"
)

// Store manages the set of APIs for user caching.
type Store struct {
	log    *logger.Logger
	storer userbus.Storer
	ttl    time.Duration
}

// NewStore constructs a cache-layer store.
func NewStore(log *logger.Logger, storer userbus.Storer, ttl time.Duration) *Store {
	return &Store{
		log:    log,
		storer: storer,
		ttl:    ttl,
	}
}

// NewWithTx implements the userbus.Storer interface.
func (s *Store) NewWithTx(tx sqldb.CommitRollbacker) (userbus.Storer, error) {
	return s.storer.NewWithTx(tx)
}

// Create adds a new user to the database and then to the cache.
func (s *Store) Create(ctx context.Context, usr userbus.User) error {
	if err := s.storer.Create(ctx, usr); err != nil {
		return err
	}

	// For simplicity, we just log here. In real life, you'd save to Redis.
	s.log.Info(ctx, "usercache.Create", "user_id", usr.ID)

	return nil
}

// Update modifies a user in the database and invalidates the cache.
func (s *Store) Update(ctx context.Context, usr userbus.User) error {
	if err := s.storer.Update(ctx, usr); err != nil {
		return err
	}

	s.log.Info(ctx, "usercache.Update", "user_id", usr.ID, "status", "cache invalidated")

	return nil
}

// Delete removes a user from the database and the cache.
func (s *Store) Delete(ctx context.Context, usr userbus.User) error {
	if err := s.storer.Delete(ctx, usr); err != nil {
		return err
	}

	s.log.Info(ctx, "usercache.Delete", "user_id", usr.ID, "status", "cache removed")

	return nil
}

// QueryByID retrieves a user from the cache or the database.
func (s *Store) QueryByID(ctx context.Context, userID uuid.UUID) (userbus.User, error) {
	// 1. Check cache (mocked here)
	s.log.Info(ctx, "usercache.QueryByID", "user_id", userID, "status", "cache miss")

	// 2. If miss, check database
	usr, err := s.storer.QueryByID(ctx, userID)
	if err != nil {
		return userbus.User{}, err
	}

	// 3. Populate cache
	s.log.Info(ctx, "usercache.QueryByID", "user_id", userID, "status", "cache populated")

	return usr, nil
}

// Implement other methods by just delegating to the inner storer for now.

func (s *Store) Query(ctx context.Context, filter userbus.QueryFilter, orderBy order.By, pg page.Page) ([]userbus.User, error) {
	return s.storer.Query(ctx, filter, orderBy, pg)
}

func (s *Store) Count(ctx context.Context, filter userbus.QueryFilter) (int, error) {
	return s.storer.Count(ctx, filter)
}

func (s *Store) QueryByEmail(ctx context.Context, email mail.Address) (userbus.User, error) {
	return s.storer.QueryByEmail(ctx, email)
}
