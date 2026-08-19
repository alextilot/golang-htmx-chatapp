// Package apperr is the transport-neutral failure vocabulary every service
// returns and every handler reads. See docs/error-handling.md for the
// ownership rules this package follows.
package apperr

import (
	"errors"
	"sort"
	"strings"
	"unicode"

	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
	"github.com/go-playground/validator/v10"
)

// Kind classifies why an operation failed. A service is the only layer
// that decides which Kind applies to a given failure.
type Kind int

const (
	KindInvalid       Kind = iota + 1 // bad input / business-rule violation — see Fields
	KindUnauthorized                  // missing/invalid credentials
	KindForbidden                     // authenticated, not permitted
	KindNotFound                      // resource does not exist
	KindConflict                      // conflicts with existing state (e.g. taken username)
)

// FieldErrors collects validation messages keyed by field name.
type FieldErrors map[string][]string

// FieldGlobal is the conventional key for a non-field-specific notice
// (e.g. "we couldn't read that request"). It is a rendering convention,
// not one of the Kind values above.
const FieldGlobal = "global"

func NewFieldErrors() FieldErrors {
	return make(FieldErrors)
}

func (f FieldErrors) Add(field, message string) {
	f[field] = append(f[field], message)
}

func (f FieldErrors) HasErrors() bool {
	return len(f) > 0
}

// Flatten combines every message into a single newline-separated string.
func (f FieldErrors) Flatten() string {
	if len(f) == 0 {
		return ""
	}
	keys := make([]string, 0, len(f))
	for k := range f {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var all []string
	for _, k := range keys {
		all = append(all, f[k]...)
	}
	return strings.Join(all, "\n")
}

func (f FieldErrors) FlattenByField() map[string]string {
	out := make(map[string]string)
	for field, msgs := range f {
		out[field] = strings.Join(msgs, "; ")
	}
	return out
}

// Error is the shared, transport-neutral shape every service returns for a
// classified/expected failure. Anything that isn't an *Error is treated by
// callers as unexpected/infra and handled generically.
type Error struct {
	Kind    Kind
	Fields  FieldErrors // field-level detail; empty for non-field kinds
	Message string      // human-readable summary for non-field kinds
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Fields.Flatten()
}

func Invalid(fields FieldErrors) *Error      { return &Error{Kind: KindInvalid, Fields: fields} }
func Unauthorized(msg string) *Error         { return &Error{Kind: KindUnauthorized, Message: msg} }
func Forbidden(msg string) *Error            { return &Error{Kind: KindForbidden, Message: msg} }
func NotFound(msg string) *Error             { return &Error{Kind: KindNotFound, Message: msg} }
func Conflict(fields FieldErrors) *Error     { return &Error{Kind: KindConflict, Fields: fields} }

// As is a thin convenience wrapper around errors.As so callers don't repeat
// the boilerplate at every call site.
func As(err error) (*Error, bool) {
	var e *Error
	ok := errors.As(err, &e)
	return e, ok
}

// ValidateStruct runs struct-tag validation (using the field-format rules
// registered in internal/validation) against s. It returns nil on success,
// or an Invalid(...) *Error on failure — the one call every service's
// input-validation step uses.
func ValidateStruct(s any) *Error {
	err := validation.Validator().Struct(s)
	if err == nil {
		return nil
	}

	fields := NewFieldErrors()
	for _, e := range err.(validator.ValidationErrors) {
		msg := validation.GetErrorMessage(e.Tag(), e.Field(), e.Param())
		if e.Tag() == "eqfield" && e.Param() == "Password" {
			msg = "Passwords do not match."
		}
		key := lowerFirst(e.Field())
		fields.Add(key, msg)
	}
	return Invalid(fields)
}

// lowerFirst lowercases a Go struct field name's first rune, so validation
// keys ("Username") match the lowercase form-field convention templates
// already key on ("username").
func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
