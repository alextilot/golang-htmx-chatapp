# Route Structure

## Target

Use explicit top-level namespaces for alternate application interfaces while treating the normal website as the default interface.

The general interface structure is:

```text
/<type>/<optional-version>/<resource>
```

The normal website is the exception and does not have a `<type>` prefix.

## Route Namespaces

```text
Website
    /

REST API
    /api/v1/

WebSocket
    /ws/

GraphQL
    /graphql

Operational
    /health
    /metrics
```

## Website

The normal website owns the root namespace.

Examples:

```text
/
/login
/logout
/chat
/chat/:room
/users/:username
```

Do not create an `/html/` prefix.

The absence of a type prefix means the request is for the normal website interface.

For example:

```text
/chat
```

is the HTML/HTMX chat page.

It must not become:

```text
/html/chat
```

## REST API

REST API routes must use the `api` namespace and an explicit version.

Examples:

```text
/api/v1/users
/api/v1/messages
/api/v1/rooms
```

Do not create an unversioned REST API namespace such as:

```text
/api/users
```

API versioning is part of the URL contract from the beginning.

Breaking API changes should use a new version:

```text
/api/v1/...
/api/v2/...
```

The API version applies to the API contract, not to the underlying services or repositories.

## WebSocket

WebSocket routes use the `ws` namespace:

```text
/ws/...
```

Examples:

```text
/ws/chat
/ws/rooms/:room
```

WebSocket routes do not require URL versioning initially.

Do not create:

```text
/ws/v1/...
```

unless a real breaking WebSocket protocol change requires it.

If versioning becomes necessary later, a versioned WebSocket namespace may be introduced:

```text
/ws/v2/...
```

## GraphQL

If GraphQL is introduced, use:

```text
/graphql
```

Do not add URL versioning by default.

Do not create:

```text
/graphql/v1
/graphql/v2
```

unless there is a concrete architectural reason to do so.

## Operational Routes

Operational endpoints are separate from application interfaces:

```text
/health
/metrics
```

These endpoints are not REST API resources and are not versioned under `/api`.

## Interface Categories

```text
Application Interfaces
├── Website
│   └── /
│
├── REST API
│   └── /api/v1/
│
├── WebSocket
│   └── /ws/
│
└── GraphQL
    └── /graphql

Operational
├── /health
└── /metrics
```

Do not introduce additional top-level interface types unless the application has a concrete need for them.

## Handler Mapping

Route namespaces map to transport-specific handlers:

```text
/...
    → handler/html/

/api/v1/...
    → handler/api/

/ws/...
    → handler/ws/
```

The URL namespace identifies the interface/transport.

It does not identify the business domain.

## Shared Application Logic

Different interfaces should act as adapters into the same service layer.

For example:

```text
/chat
/api/v1/messages
/ws/chat
```

may all ultimately use the same chat/message services.

The architecture should be:

```text
                    Application
                   /     |      \
                  /      |       \
              Website   REST      WS
                 │       │        │
                 └───────┼────────┘
                         │
                      Services
                         │
                    Repositories
```

Do not duplicate application/business logic between HTML, API, and WebSocket handlers.

Transport-specific behavior belongs in the corresponding handler package.

Application behavior belongs in `service/`.

## Routing Responsibilities

The router is responsible for:

- registering routes
- grouping routes
- attaching middleware
- connecting routes to handlers

The router must not contain application/business logic.

## Routing Rules

1. The normal website owns `/`.
2. Do not create an `/html/` namespace.
3. REST API routes use `/api/v1/...`.
4. REST API versioning is mandatory from the beginning.
5. Breaking REST API changes use a new API version.
6. WebSocket routes use `/ws/...`.
7. WebSocket routes are not versioned unless a real breaking protocol change requires it.
8. GraphQL uses `/graphql` without preemptive URL versioning.
9. Operational endpoints use `/health` and `/metrics`.
10. Route prefixes represent interfaces/transports, not business domains.
11. Different interfaces should reuse the same services.
12. Transport-specific request/response behavior belongs in handlers.
13. Business/application behavior belongs in services.
14. The router only handles HTTP route and middleware wiring.
