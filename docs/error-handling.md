# Error Handling & Validation

## Target

Give every interface (website, REST API, WebSocket, and any future interface) the same transport-neutral vocabulary for describing why an operation failed, and validate input in exactly one place: the service that owns the operation.

Do not let validation rules live in more than one layer, and do not let any layer decide what an operation's input must look like except the service that performs the operation.

## Structure

```text
internal/
├── apperr/          # Shared failure vocabulary
├── validation/       # Field-format primitives
├── service/           # Input types, validation, business logic, failure classification
├── handler/            # Request shape, transport-specific response translation
└── repository/          # Persistence
```

## `internal/apperr`

Owns the one vocabulary every layer between `repository` and the transport agrees on:

- A small, fixed set of failure kinds: invalid input, unauthorized, forbidden, not found, conflict.
- A field-level detail bag, for failures tied to specific input fields.

Do not put HTTP status codes, response rendering, or any transport-specific concept here.

Do not put the definition of what makes a field's format valid here — that belongs to `internal/validation`. This package only carries the *result* of validation, not the rules behind it.

## `internal/validation`

Owns field-format primitives that are true independent of any operation or transport — what a valid password looks like, and similar low-level rules.

Do not put any operation's input shape, any business rule, or the failure vocabulary here. If something here starts describing what a specific operation needs rather than what a field's format is, it belongs in `internal/service` instead.

This package should have exactly one caller.

## `internal/service`

Owns, for every operation it exposes:

- The shape of its own input. Not shared with, and not shaped by, any transport.
- All validation of that input — both format rules and business rules (uniqueness, cross-entity checks, and so on) — run unconditionally, every time the operation is called, regardless of which transport called it.
- Classifying every failure it can produce into the shared failure vocabulary. The service is the only layer that knows *why* something failed, so it is the only layer allowed to decide that.

Do not let a service depend on any transport's request format, HTTP concepts, or response format.

## `internal/handler` (one package per transport)

Owns:

- The transport's own request shape (form fields, JSON body, or whatever the transport speaks).
- A boring, one-to-one mapping from that request shape into the service's input type.
- Translating the service's outcome — success, a classified failure, or an unexpected failure — into this transport's own response representation (HTTP status and page for the website, HTTP status and body for the REST API, and so on).

Do not validate input here beyond what's required to bind it. If a handler is deciding whether input is *acceptable* rather than just well-formed enough to read, that decision is misplaced and belongs in the service.

Do not let one transport's handler package know about another transport's request or response shape.

## `internal/repository`

Reports whether something exists, was stored, or failed, in its own neutral terms.

Do not reference the shared failure vocabulary here, and do not decide what a failure means to the caller — that is the service's job. A repository answers "does this exist," not "should this be allowed" or "what does the caller show the user."

## Rules

1. Every operation's input shape and validation rules live in exactly one service method.
2. Failures a caller needs to react to are classified into the shared failure vocabulary; failures nothing anticipated are not classified and are handled generically.
3. Only a service decides which failure kind applies. No other layer makes that decision.
4. A repository never speaks the shared failure vocabulary. It reports outcomes in its own terms; the service interprets them.
5. Field-format rules are defined once, independently of any specific operation or transport.
6. Every transport maps the same shared failure vocabulary to its own protocol's representation. No transport invents its own status/error scheme.
7. A transport validates only enough to bind a request. All content and business validation happens in the service.
8. Two transports are allowed to have differently-shaped requests for what is conceptually the same operation — that divergence is expected, not a bug to reconcile.

## Philosophy

Validate once, in the layer that owns the rule. Classify failures once, in a vocabulary every interface can read. Let each transport decide, on its own, how that vocabulary should look to its callers.
