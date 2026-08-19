# UI Structure

## Target

Keep `web/` organized by technical role — the same principle `docs/repository-structure.md`
applies to `internal/` — so that adding a feature's UI (e.g. group chat: a sidebar,
a profile widget, a message list) doesn't turn any one folder into a dumping ground.

Two rules do most of the work:

1. A component's package name always matches its own directory name, never its parent's.
2. Every `components/` subdirectory is either **generic** (no domain types) or a
   **feature namespace** (named after the domain it renders) — never a flat pile of
   unrelated `.templ` files.

## Structure

```text
web/
├── layouts/          # Page shells: Document, LayoutBase, LayoutFullHeight
├── pages/             # One templ per route. Thin: composes layouts + components.
│   └── status/        # Subfolder once a family of related pages exists.
├── partials/          # Site-wide chrome present on every page (nav bar, footer, theme picker)
├── components/
│   ├── form/           # Generic field-level primitives (Input, Button, Select, ...)
│   └── chat/            # Feature namespace: chatroom, message bubble, and (as
│                          #   group-chat work lands) sidebar, group list item,
│                          #   profile widget, member list — anything that renders
│                          #   a Group/Message/User belongs here, not in pages/.
├── forms/              # Composed, page-level forms that own a submission flow
├── static/             # CSS, JS, images
├── utils/
└── templ.go
```

`web/` lives at the repository root, not under `internal/` — nothing outside the
module needs Go's `internal/` import boundary to be enforced against templates.

## `web/layouts/`

Page shells only: `<html>`/`<head>`/`<body>` skeleton, sticky header, optional
footer. A layout knows nothing about what feature it's wrapping.

## `web/pages/`

One file per route. A page composes a layout plus one or more top-level
components — it must not reach into `service`/`repository` types directly or
contain business logic; that stays in `internal/handler`, which passes the page
whatever data it needs to render.

## `web/partials/`

Chrome that's present on (nearly) every page regardless of feature — nav bar,
footer, theme dropdown. If it's specific to one page or feature, it belongs in
`components/`, not here.

## `web/components/`

Never place a `.templ` file directly in `components/` — always in a named
subdirectory, so the directory-per-package convention below always holds. Each
subdirectory is one of:

- **Generic** — presentation-only, doesn't know about `model.Group`,
  `model.Message`, `model.User`, etc. Example: `components/form/` (`Input`,
  `Button`, `Select`, `Checkbox`, `RadioGroup`, `Textarea`, `Alert`, `Error`).
- **Feature namespace** — named after the domain/feature it renders, and owns
  every component tied to that feature. Example: `components/chat/`.

A feature folder is allowed to grow — keep it flat as long as everything in it
belongs to that one feature. Only split a piece out further once it needs to be
reused *outside* that feature.

Don't add a `components/ui/` wrapper preemptively to group "generic" folders —
introduce it only once there are enough generic categories (`form/`, `avatar/`,
`dropdown/`, ...) that a flat `components/` is genuinely hard to scan.

### Package naming

A directory's package name always matches the directory itself — e.g.
`components/chat/` declares `package chat`, `components/form/` declares
`package form`. Do not let a subpackage declare the same package name as its
parent (`components/form/` previously declared `package components`, identical
to `components/`'s own package name — harmless until a file needed to import
both, at which point one side would require an alias). Matching Go's normal
directory-equals-package convention avoids that class of bug entirely.

## `web/forms/`

Composed, page-level forms — `LoginForm`, `SignupForm` — that own a submission
flow: the `<form>` element, `action`, `hx-target`, `hx-swap`. They're built out
of `components/form/` primitives, which know nothing about where they submit to.
Rule of thumb: if it knows its own route/action, it's `forms/`; if it's a
reusable field with no route knowledge, it's `components/form/`.

## Rules

1. Never place a `.templ` file directly in `components/` — every component lives in a named subdirectory.
2. A `components/` subdirectory is either generic (no domain types) or a feature namespace (named after the domain it renders).
3. A package's name always matches its own directory, never its parent's.
4. Pages stay thin: compose layouts + components, no direct repository/service calls, no business logic.
5. Partials are for chrome present on every page; anything page- or feature-specific belongs in `components/`.
6. `forms/` holds full submission flows; `components/form/` holds field-level primitives with no route/action knowledge.
7. Don't introduce a `components/ui/` wrapper until there are enough generic categories that a flat `components/` is hard to scan.
8. Don't reorganize `web/` into per-route or per-feature top-level folders — keep the layout/page/partial/component/form separation by technical role.

## Philosophy

Prefer explicit, boring layers over unnecessary abstraction — the same
philosophy `docs/repository-structure.md` applies to `internal/`. A feature
should be easy to find (its components live under one named folder) without
the presentation layer being reorganized around features wholesale.
