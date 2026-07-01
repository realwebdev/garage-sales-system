package main

import (
	"context"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/realwebdev/garage-sales-system/app/sdk/mux"
	"github.com/realwebdev/garage-sales-system/business/domain/userbus"
	"github.com/realwebdev/garage-sales-system/business/domain/userbus/stores/usercache"
	"github.com/realwebdev/garage-sales-system/business/sdk/delegate"
	"github.com/realwebdev/garage-sales-system/business/sdk/order"
	"github.com/realwebdev/garage-sales-system/business/sdk/page"
	"github.com/realwebdev/garage-sales-system/business/sdk/sqldb"
	"github.com/realwebdev/garage-sales-system/foundation/logger"
)

func main() {
	log := logger.New(os.Stdout, logger.LevelInfo, "SALES", nil)

	if err := run(log); err != nil {
		log.Error(context.Background(), "startup", "err", err)
		os.Exit(1)
	}
}

func run(log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// Start Tracing / Metrics (Placeholder for now)

	// -------------------------------------------------------------------------
	// Setup Business Layer
	//
	// This demonstrates DEPENDENCY INVERSION:
	// 1. The 'userbus' package defines the 'Storer' interface it needs.
	// 2. We can provide any implementation of 'Storer' here.
	// 3. Here we use an in-memory 'mockStorer' for simple running.
	// 4. We can wrap it in a 'usercache' (Decorator Pattern) which also implements 'Storer'.
	// 5. 'userBus' doesn't know (or care) if it's talking to memory, cache, or a real DB.
	
	storer := newMockStorer()

	// Optional: Wrap the storer in a cache layer.
	// storer = usercache.NewStore(log, storer, time.Minute)

	delegate := delegate.New(log)
	userBus := userbus.NewBusiness(log, delegate, storer)

	// -------------------------------------------------------------------------
	// Start API Service

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	cfg := mux.Config{
		Build:    "develop",
		Log:      log,
		Shutdown: shutdown,
		UserBus:  userBus,
	}

	api := http.Server{
		Addr:         ":8080",
		Handler:      mux.WebAPI(cfg),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		log.Info(context.Background(), "startup", "status", "api router started", "host", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(context.Background(), "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(context.Background(), "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}

// -------------------------------------------------------------------------
// Mock Storer for demonstration

type mockStorer struct {
	users map[string]userbus.User
}

func newMockStorer() *mockStorer {
	return &mockStorer{
		users: make(map[string]userbus.User),
	}
}

func (m *mockStorer) NewWithTx(tx sqldb.CommitRollbacker) (userbus.Storer, error) { return m, nil }
func (m *mockStorer) Create(ctx context.Context, usr userbus.User) error {
	m.users[usr.ID.String()] = usr
	return nil
}
func (m *mockStorer) Update(ctx context.Context, usr userbus.User) error {
	m.users[usr.ID.String()] = usr
	return nil
}
func (m *mockStorer) Delete(ctx context.Context, usr userbus.User) error {
	delete(m.users, usr.ID.String())
	return nil
}
func (m *mockStorer) Query(ctx context.Context, filter userbus.QueryFilter, orderBy order.By, pg page.Page) ([]userbus.User, error) {
	var users []userbus.User
	for _, u := range m.users {
		users = append(users, u)
	}
	return users, nil
}
func (m *mockStorer) Count(ctx context.Context, filter userbus.QueryFilter) (int, error) {
	return len(m.users), nil
}
func (m *mockStorer) QueryByID(ctx context.Context, userID uuid.UUID) (userbus.User, error) {
	u, ok := m.users[userID.String()]
	if !ok {
		return userbus.User{}, userbus.ErrNotFound
	}
	return u, nil
}
func (m *mockStorer) QueryByEmail(ctx context.Context, email mail.Address) (userbus.User, error) {
	for _, u := range m.users {
		if u.Email.Address == email.Address {
			return u, nil
		}
	}
	return userbus.User{}, userbus.ErrNotFound
}
