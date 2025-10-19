package core

import (
	"fmt"
	"sort"
	"strings"
)

// FieldErrors represents a collection of validation errors, where each map key
// corresponds to a field name and the value is a slice of error messages for that field.
type FieldErrors map[string][]string

func NewFieldErrors() FieldErrors {
	return make(FieldErrors)
}

func (f FieldErrors) Add(field, message string) {
	f[field] = append(f[field], message)
}

func (f FieldErrors) HasErrors() bool {
	return len(f) > 0
}

func (f FieldErrors) Error() string {
	if len(f) == 0 {
		return "no validation errors"
	}

	// Collect summaries like: "email: required, invalid"
	var parts []string
	for field, msgs := range f {
		parts = append(parts, fmt.Sprintf("%s: %s", field, strings.Join(msgs, ", ")))
	}

	return "validation failed: " + strings.Join(parts, "; ")
}

// Flatten combines all error messages from the FieldErrors map into a single string,
// with each message on a new line.
func (fe FieldErrors) Flatten() string {
	if len(fe) == 0 {
		return ""
	}
	keys := make([]string, 0, len(fe))
	for k := range fe {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var all []string
	for _, k := range keys {
		all = append(all, fe[k]...)
	}
	return strings.Join(all, "\n")
}

func (fe FieldErrors) FlattenByField() map[string]string {
	out := make(map[string]string)
	for field, msgs := range fe {
		out[field] = strings.Join(msgs, "; ")
	}
	return out
}
