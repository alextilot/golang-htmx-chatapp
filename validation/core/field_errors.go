package core

import (
	"sort"
	"strings"
)

// FieldErrors represents a collection of validation errors, where each map key
// corresponds to a field name and the value is a slice of error messages for that field.
type FieldErrors map[string][]string

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
