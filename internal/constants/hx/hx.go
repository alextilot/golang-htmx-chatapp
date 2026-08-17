// Package hx contains HTMX-specific HTTP header names. These have no
// equivalent in echo/v5, which is not HTMX-aware.
package hx

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
