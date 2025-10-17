// Package headerval contains common HTTP header values.
// It centralizes header values like Content-Types to avoid duplication and typos.
package headerval

// -----------------------
// Content Types
// -----------------------
const (
	// ContentTypeJSON is the MIME type for JSON responses.
	ContentTypeJSON = "application/json"

	// ContentTypeHTML is the MIME type for HTML responses with UTF-8 charset.
	ContentTypeHTML = "text/html; charset=utf-8"

	// ContentTypeFormURLEncoded is the MIME type for form submissions.
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"

	// ContentTypeMultipartForm is the MIME type for multipart form submissions (file uploads).
	ContentTypeMultipartForm = "multipart/form-data"
)
