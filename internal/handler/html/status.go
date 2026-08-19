package html

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/apperr"
)

// statusFor maps a classified service failure to the HTTP status this
// website responds with. Other transports (api, ...) own their own mapping
// from the same apperr.Kind values — see docs/error-handling.md.
func statusFor(kind apperr.Kind) int {
	switch kind {
	case apperr.KindUnauthorized:
		return http.StatusUnauthorized
	case apperr.KindForbidden:
		return http.StatusForbidden
	case apperr.KindNotFound:
		return http.StatusNotFound
	case apperr.KindConflict:
		return http.StatusConflict
	default: // apperr.KindInvalid
		return http.StatusBadRequest
	}
}

// fieldsOf returns appErr's field-level detail for rendering, falling back
// to a single FieldGlobal-keyed notice for kinds that don't carry field
// detail (Unauthorized, Forbidden, NotFound).
func fieldsOf(appErr *apperr.Error) apperr.FieldErrors {
	if appErr.Fields.HasErrors() {
		return appErr.Fields
	}
	return apperr.FieldErrors{apperr.FieldGlobal: {appErr.Message}}
}
