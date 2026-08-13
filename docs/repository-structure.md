# Repository Structure

## Target

Use a pragmatic layered architecture. Keep technical concerns grouped by layer rather than organizing the code by domain or feature.

Do not introduce full DDD, feature-oriented packages, separate domain entities, or DTO mapping unless a concrete requirement appears later.

The current GORM models remain the persistence models.

## Structure

```text
golang-htmx-chatapp/
│
├── cmd/
│   └── main/
│       └── main.go
│
├── db/
│   ├── db.go
│   └── migrations.go
│
├── model/                 # GORM persistence models
│   ├── user.go
│   ├── message.go
│   └── ...
│
├── repository/            # Database/persistence access
│   ├── repositories.go    # NewRepositories()
│   ├── user.go
│   ├── message.go
│   └── ...
│
├── service/               # Application/business logic
│   ├── services.go        # NewServices()
│   ├── user.go
│   ├── message.go
│   ├── auth.go
│   └── ...
│
├── handler/               # HTTP presentation/transport
│   ├── html/
│   │   ├── handler.go
│   │   ├── user.go
│   │   ├── chatroom.go
│   │   └── requests.go
│   │
│   ├── api/
│   │   ├── handler.go
│   │   ├── user.go
│   │   ├── message.go
│   │   └── ...
│   │
│   └── ws/
│       ├── handler.go
│       └── ...
│
├── router/
│   └── router.go
│
└── web/
    ├── templates/
    └── static/
```

## Layers

### `model/`

Contains GORM persistence models.

Models represent how application data is persisted.

Models may contain:

- GORM tags
- database relationships
- persistence-specific fields
- GORM lifecycle behavior where appropriate

Do not create separate domain models or DTOs solely to satisfy an architectural pattern.

### `repository/`

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

### `service/`

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

### `handler/`

Contains HTTP transport and presentation logic.

Handlers are separated by interface:

```text
handler/html/    HTML/HTMX website
handler/api/     REST/JSON API
handler/ws/      WebSocket
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

### `router/`

Owns route registration and middleware wiring.

The router is responsible for:

- registering routes
- connecting routes to handlers
- configuring middleware
- configuring route groups

The router must not contain application/business logic.

### `web/`

Contains presentation assets:

- HTML templates
- HTMX templates/fragments
- CSS
- JavaScript
- static assets

## Dependency Direction

Dependencies should flow toward the application logic and persistence layers:

```text
Router
  ↓
Handlers
  ↓
Services
  ↓
Repositories
  ↓
GORM Models
  ↓
Database
```

Different HTTP interfaces must reuse the same service layer rather than duplicating business logic.

## Centralized Construction

Keep dependency construction centralized.

Repositories are constructed first:

```go
repos := repository.NewRepositories(db)
```

Services are constructed from the repositories:

```go
services := service.NewServices(repos)
```

Handlers are constructed from the services:

```go
htmlHandlers := html.NewHandlers(services)
apiHandlers := api.NewHandlers(services)
wsHandlers := ws.NewHandlers(services)
```

Handlers and repositories must not independently construct their own application dependencies.

The composition root is responsible for wiring the application together.

## Architectural Rules

1. Keep the repository organized by technical layer.
2. Do not reorganize the repository into domain/feature packages.
3. Keep GORM models in `model/`.
4. Replace `store/` with `repository/`.
5. Keep application/business logic in `service/`.
6. Keep HTTP presentation logic in `handler/`.
7. Separate HTML, REST API, and WebSocket handlers.
8. Keep route registration in `router/`.
9. Keep templates and static assets in `web/`.
10. Keep `NewRepositories()` as the centralized repository construction mechanism.
11. Keep `NewServices()` as the centralized service construction mechanism.
12. Services may orchestrate other services.
13. Do not introduce DTOs or separate domain models without a concrete requirement.
14. Do not duplicate business logic between HTTP interfaces.
15. Keep HTTP concerns out of services and repositories.

## Philosophy

Prefer explicit, boring layers over unnecessary abstraction.

The architecture should make these responsibilities obvious:

```text
Handlers       → transport/presentation
Services       → application/business behavior
Repositories   → persistence
Models         → GORM/database representation
Router         → HTTP wiring
```
