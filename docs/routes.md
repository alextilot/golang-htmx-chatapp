# Route Design Guidelines

## HTML vs API Routes

### Recommended Pattern

- **HTML routes:** Use “normal” paths for user-facing pages and HTMX form endpoints.
  Example:
  /login         → HTML page or HTML partial
  /signup        → HTML page or HTML partial
  /chat          → HTML chat page

- **API / JSON routes:** Prefix paths with `/api/` for JSON endpoints.
  Example:
  /api/login     → JSON login endpoint
  /api/users     → JSON user data
  /api/chat      → JSON chat messages

### Why This Works Well

1. **Clear separation of concerns:** You immediately know the response type from the URL.  
2. **Easier for front-end and external clients:** No guesswork on headers required.  
3. **Scalable:** You can version APIs easily (`/api/v1/...`).  
4. **Shared business logic:** HTML and JSON handlers call the same service layer, keeping logic DRY.  
5. **HTMX-friendly:** HTMX requests still hit the HTML routes, returning partials or full pages as needed.

### Notes

- Polymorphic endpoints (one route serving both HTML and JSON) are possible, but harder to maintain, document, and scale.  
- Keep the service layer separate from handlers so formatting (HTML vs JSON) is only handled in the handler.
