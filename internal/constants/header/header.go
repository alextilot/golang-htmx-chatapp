package header

// -----------------------
// Request Headerss
// -----------------------
const (
	// Accept indicates the MIME types the client can accept.
	Accept = "Accept"

	// UserAgent identifies the client software initiating the request.
	UserAgent = "User-Agent"

	// Referer indicates the address of the previous web page from which a request was made.
	Referer = "Referer"

	// Origin indicates the origin of the request, used for CORS.
	Origin = "Origin"

	// XRequestID is a unique ID for the HTTP request, often used for logging.
	XRequestID = "X-Request-ID"

	// XCorrelationID is used to correlate logs and requests across services.
	XCorrelationID = "X-Correlation-ID"

	// XForwardedFor identifies the originating IP address of a client connecting through a proxy.
	XForwardedFor = "X-Forwarded-For"
)

// -----------------------
// Response & Security Headers
// -----------------------
const (
	// ContentType is the standard HTTP header for Content-Type.
	ContentType = "Content-Type"

	// Authorization carries authentication credentials (e.g., Bearer token).
	Authorization = "Authorization"

	// WWWAuthenticate prompts client authentication when 401 Unauthorized is returned.
	WWWAuthenticate = "WWW-Authenticate"

	// CacheControl controls caching behavior for responses.
	CacheControl = "Cache-Control"

	// ETag provides an entity tag for caching validation.
	ETag = "ETag"

	// LastModified indicates when the resource was last modified.
	LastModified = "Last-Modified"

	// Expires specifies the expiration date/time for the response.
	Expires = "Expires"

	// AccessControlAllowOrigin specifies permitted CORS origins.
	AccessControlAllowOrigin = "Access-Control-Allow-Origin"

	// AccessControlAllowMethods specifies permitted CORS HTTP methods.
	AccessControlAllowMethods = "Access-Control-Allow-Methods"

	// AccessControlAllowHeaders specifies permitted CORS headers.
	AccessControlAllowHeaders = "Access-Control-Allow-Headers"

	// AccessControlExposeHeaders specifies headers exposed to the browser in CORS.
	AccessControlExposeHeaders = "Access-Control-Expose-Headers"
)

// -----------------------
// HTMX Headers
// -----------------------
const (
	// HXBoosted is sent by HTMX when the request is triggered by hx-boost.
	HXBoosted = "HX-Boosted"

	// HXRequest is sent by HTMX to indicate an AJAX request.
	HXRequest = "HX-Request"

	// HXRedirect instructs HTMX to perform a client-side redirect.
	HXRedirect = "HX-Redirect"

	// HXTrigger triggers a client-side event on an element.
	HXTrigger = "HX-Trigger"

	// HXTriggerAfterSettle triggers events after content settles in HTMX.
	HXTriggerAfterSettle = "HX-Trigger-After-Settle"

	// HXSwap defines how HTMX should swap the returned content.
	HXSwap = "HX-Swap"

	// HXTarget specifies the target element selector for HTMX swaps.
	HXTarget = "HX-Target"

	// HXPushURL updates the browser's URL without a full page reload.
	HXPushURL = "HX-Push-URL"
)
