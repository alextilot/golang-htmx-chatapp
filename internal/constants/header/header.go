// Package header contains HTTP header name constants that echo/v5 does not
// already provide. For headers echo defines (Content-Type, Authorization,
// Cache-Control, Access-Control-*, X-Request-Id, etc.), use echo.Header*
// directly instead of duplicating them here.
package header

const (
	// UserAgent identifies the client software initiating the request.
	UserAgent = "User-Agent"

	// Referer indicates the address of the previous web page from which a request was made.
	Referer = "Referer"

	// ETag provides an entity tag for caching validation.
	ETag = "ETag"

	// Expires specifies the expiration date/time for the response.
	Expires = "Expires"
)
