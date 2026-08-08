# Garage Sales System - Mini Project

This project is a minimal implementation of the Ardan Labs Service pattern. It demonstrates structured logging, observability, middleware, Redis caching, and dependency inversion.

## Project Structure

```
.
├── api/
│   └── services/
│       └── sales/          # Service entrypoint (main.go)
├── app/
│   ├── domain/
│   │   ├── checkapp/       # Health check handlers
│   │   └── userapp/        # User CRUD handlers
│   └── sdk/
│       ├── errs/           # Error handling system
│       ├── mid/            # Middleware (logger, etc.)
│       └── mux/            # API router wiring
├── business/
│   ├── domain/
│   │   └── userbus/        # User business logic
│   │       ├── stores/
│   │       │   ├── usercache/  # Redis cache layer (Decorator)
│   │       │   └── userdb/     # PostgreSQL store
│   └── sdk/                # Shared business logic
└── foundation/
    ├── logger/             # Structured logging (slog)
    └── web/                # Custom web framework
```

## Key Concepts

### 1. Structured Logging
We use the `foundation/logger` package which wraps `log/slog`. It provides JSON output by default, which is essential for production log aggregation (Loki, ELK).

Example usage in `main.go` and `middleware`:
```go
log.Info(ctx, "request started", "method", r.Method, "path", r.URL.Path)
```

### 2. Dependency Inversion (DIP)
The `userbus` package defines a `Storer` interface. The business logic only depends on this interface, not on a concrete database implementation.
This allows us to:
- Swap PostgreSQL for an In-Memory store during tests.
- Wrap the store with a Cache layer (Decorator pattern) without changing business logic.

In `main.go`, we wire it up:
```go
storer := userdb.NewStore(...)
cache := usercache.NewStore(log, storer, time.Minute) // usercache implements Storer
userBus := userbus.NewBusiness(log, delegate, cache)
```

### 3. Middleware (Decorators)
Middleware are functions that wrap handlers to add cross-cutting concerns like logging, recovery, and metrics.
Our `web.App` supports global and per-route middleware.

Example `mid.Logger`:
```go
func Logger(log *logger.Logger) web.MidFunc {
    return func(handler web.HandlerFunc) web.HandlerFunc {
        return func(ctx context.Context, r *http.Request) web.Encoder {
            // ... logging ...
            return handler(ctx, r)
        }
    }
}
```

### 4. Redis Cache Layer
The `usercache` package implements the `userbus.Storer` interface. It wraps a "real" storer (like `userdb`).
- On `QueryByID`: It first checks Redis. If a miss, it calls the database and populates Redis.
- On `Update/Delete`: It invalidates the Redis entry.

### 5. Observability
- **Health Checks**: Handled by `checkapp` (`/health`).
- **Tracing/Metrics**: In a full project, we would use OpenTelemetry (OTel). We've laid the groundwork in `foundation/web` and `app/sdk/mid`.

## How to Run

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Run the service**:
   ```bash
   go run api/services/sales/main.go
   ```

3. **Test the endpoints**:
   - Health check: `curl http://localhost:8080/health`
   - List users: `curl http://localhost:8080/users`
   - Create user:
     ```bash
     curl -X POST http://localhost:8080/users -d '{
       "name": "Haseeb",
       "email": "haseeb@example.com",
       "roles": ["User"],
       "password": "password123"
     }'
     ```

## Deployment (Zarf / Docker)

### Docker Compose
Create a `docker-compose.yaml` with:
- `postgres`: For persistent storage.
- `redis`: For caching.
- `sales-api`: This application.
- `grafana/prometheus/tempo/loki`: The observability stack.

(Detailed compose and k8s files would be in the `zarf/` directory).
