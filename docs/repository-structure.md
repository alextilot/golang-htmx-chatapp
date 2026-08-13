# Repository Structure

## Target

Use a pragmatic layered architecture with all application code under `internal/`.

The repository may contain multiple executable applications. Each executable lives under `cmd/` and composes the shared application layers from `internal/`.

Keep technical concerns grouped by layer rather than organizing the code by domain or feature.

Do not introduce full DDD, feature-oriented packages, separate domain entities, or DTO mapping unless a concrete requirement appears later.

The current GORM models remain the persistence models.

## Structure

```text
golang-htmx-chatapp/
│
├── cmd/                            # Executable applications
│   ├── server/
│   │   └── main.go
│   │
│   ├── worker/
│   │   └── main.go
│   │
│   └── ...
│
├── internal/                       # Application-private code
│   │
│   ├── db/
│   │   ├── db.go
│   │   └── migrations.go
│   │
│   ├── model/                     # GORM persistence models
│   │   ├── user.go
│   │   ├── message.go
│   │   └── ...
│   │
│   ├── repository/                # Database/persistence access
│   │   ├── repositories.go        # NewRepositories()
│   │   ├── user.go
│   │   ├── message.go
│   │   └── ...
│   │
│   ├── service/                   # Application/business logic
│   │   ├── services.go            # NewServices()
│   │   ├── user.go
│   │   ├── message.go
│   │   ├── auth.go
│   │   └── ...
│   │
│   ├── handler/                   # HTTP presentation/transport
│   │   ├── html/
│   │   │   ├── handler.go
│   │   │   ├── user.go
│   │   │   ├── chatroom.go
│   │   │   └── requests.go
│   │   │
│   │   ├── api/
│   │   │   ├── handler.go
│   │   │   ├── user.go
│   │   │   ├── message.go
│   │   │   └── ...
│   │   │
│   │   └── ws/
│   │       ├── handler.go
│   │       └── ...
│   │
│   ├── router/
│   │   └── router.go
│   │
│   └── web/
│       ├── templates/
│       └── static/
│
├── migrations/
│   └── ...
│
├── go.mod
└── ...
```

## `cmd/`

`cmd/` contains executable applications.

Each application has its own `main.go` and is responsible for composing and starting the functionality it needs.

Examples:

```text
cmd/server/main.go
cmd/worker/main.go
cmd/migrate/main.go
```

Applications should not contain reusable business logic.

The normal server application should:

1. Load configuration.
2. Initialize infrastructure.
3. Construct repositories.
4. Construct services.
5. Construct handlers.
6. Construct the router.
7. Start the server.

Additional applications may reuse the same `internal/` layers.

Do not duplicate service or repository implementations between applications.

## `internal/`

`internal/` contains application-private code.

All shared application layers belong under `internal/`.

```text
internal/
├── db/
├── model/
├── repository/
├── service/
├── handler/
├── router/
└── web/
```

Code outside the module cannot import these packages directly.

## `internal/model/`

Contains GORM persistence models.

Models represent how application data is persisted.

Models may contain:

- GORM tags
- database relationships
- persistence-specific fields
- GORM lifecycle behavior where appropriate

Do not create separate domain models or DTOs solely to satisfy an architectural pattern.

## `internal/repository/`

Owns persistence and database access.

Repositories are responsible for:

- GORM operations
- queries
- creating records
- updating records
- deleting records
- loading `model` types
- persisting `model` types

Repositories must not contain:

- HTTP logic
- HTML/template rendering
- JSON response formatting
- JWT/cookie handling
- application workflows

Replace the existing `store/` layer with `repository/`.

Repositories are constructed centrally through `NewRepositories()`.

## `internal/service/`

Contains application and business logic.

Services are responsible for:

- application rules
- business operations
- coordinating repositories
- orchestrating workflows
- calling other services when required

Services must not depend on HTTP presentation concerns.

Services must not directly handle:

- HTTP requests
- HTTP responses
- Echo handlers/context
- HTML templates
- JSON response formatting

Services may depend on and orchestrate other services when a workflow requires it.

Services are constructed centrally through `NewServices()`.

## `internal/handler/`

Contains HTTP transport and presentation logic.

Handlers are separated by interface:

```text
internal/handler/html/    HTML/HTMX website
internal/handler/api/     REST/JSON API
internal/handler/ws/      WebSocket
```

Handlers are responsible for:

- parsing HTTP input
- binding requests
- transport-specific validation
- calling services
- formatting responses
- rendering HTML
- serializing JSON
- managing WebSocket connections

Handlers must not contain core application/business workflows.

HTML, API, and WebSocket handlers may use the same services.

## `internal/router/`

Owns HTTP route registration and middleware wiring.

The router is responsible for:

- registering routes
- connecting routes to handlers
- configuring middleware
- configuring route groups

The router must not contain application/business logic.

## `internal/web/`

Contains presentation assets:

- HTML templates
- HTMX templates/fragments
- CSS
- JavaScript
- static assets

## `internal/db/`

Contains database infrastructure and initialization.

This package is responsible for:

- opening database connections
- configuring GORM
- database initialization
- database-specific setup

Database schema migrations should remain separate if they are managed as explicit migration files.

## `migrations/`

Contains database migration files when migrations are managed outside the Go database initialization package.

Migrations are persistence infrastructure, not application/business logic.

## Dependency Direction

Dependencies should flow through the layers:

```text
cmd/*
  ↓
internal/router
  ↓
internal/handler
  ↓
internal/service
  ↓
internal/repository
  ↓
internal/model
  ↓
Database
```

Infrastructure initialization may be shared by multiple applications:

```text
cmd/server ──┐
cmd/worker ──┼──→ internal/db
cmd/... ─────┘
```

Different HTTP interfaces must reuse the same service layer rather than duplicating business logic.

## Centralized Construction

Application construction should happen in the executable/composition root.

For example:

```go
repos := repository.NewRepositories(db)
services := service.NewServices(repos)

htmlHandlers := html.NewHandlers(services)
apiHandlers := api.NewHandlers(services)
wsHandlers := ws.NewHandlers(services)

router := router.New(
    htmlHandlers,
    apiHandlers,
    wsHandlers,
)
```

The executable decides which infrastructure and interfaces it needs.

Repositories, services, and handlers must not independently construct their own application dependencies.

Do not create a global dependency container.

## Multiple Applications

Multiple executables may use different subsets of the application layers.

For example:

```text
cmd/server
    → db
    → repository
    → service
    → handler
    → router

cmd/worker
    → db
    → repository
    → service

cmd/migrate
    → database/migration infrastructure
```

A worker should not initialize HTTP handlers or a router.

A migration executable should not initialize application services unless required by the migration mechanism.

Each executable should compose only the dependencies it needs.

## Architectural Rules

1. Put application-private code under `internal/`.
2. Keep executable entry points under `cmd/`.
3. Organize `internal/` by technical layer.
4. Do not reorganize the repository into domain/feature packages.
5. Keep GORM models in `internal/model/`.
6. Replace `store/` with `internal/repository/`.
7. Keep application/business logic in `internal/service/`.
8. Keep HTTP presentation logic in `internal/handler/`.
9. Separate HTML, REST API, and WebSocket handlers.
10. Keep route registration in `internal/router/`.
11. Keep database infrastructure in `internal/db/`.
12. Keep templates and static assets in `internal/web/`.
13. Keep executable-specific composition in `cmd/*`.
14. Multiple applications may reuse the same `internal/` layers.
15. Do not duplicate business logic between applications.
16. Keep `NewRepositories()` as the centralized repository construction mechanism.
17. Keep `NewServices()` as the centralized service construction mechanism.
18. Services may orchestrate other services.
19. Do not introduce DTOs or separate domain models without a concrete requirement.
20. Do not duplicate business logic between HTTP interfaces.
21. Keep HTTP concerns out of services and repositories.
22. Each executable should initialize only the layers it actually needs.

## Philosophy

Prefer explicit, boring layers over unnecessary abstraction.

The architecture should make these responsibilities obvious:

```text
cmd/*         → application entry points/composition
router        → HTTP routing and middleware
handlers      → transport/presentation
services      → application/business behavior
repositories  → persistence
models        → GORM/database representation
db            → database infrastructure
```
