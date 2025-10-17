package utils

import (
	"github.com/a-h/templ"
	"maps"
	"strings"
)

// ExtractAttribute extracts a typed attribute from templ.Attributes and removes it.
// If the key doesn't exist or the type doesn't match, it returns the zero value for T.
func ExtractAttribute[T any](attrs templ.Attributes, key string) (T, templ.Attributes) {
	var zero T
	if v, ok := attrs[key]; ok {
		if typed, ok := v.(T); ok {
			delete(attrs, key)
			return typed, attrs
		}
		delete(attrs, key)
	}
	return zero, attrs
}

// MergeAttrs merges two templ.Attributes maps.
// Keys in overrides will replace those in defaults.
func MergeAttrs(defaults, overrides templ.Attributes) templ.Attributes {
	merged := templ.Attributes{}
	maps.Copy(merged, defaults)
	maps.Copy(merged, overrides)
	return merged
}

// FieldHasErrors returns true if the provided error slice has entries.
func FieldHasErrors(fieldErrors []string) bool {
	return len(fieldErrors) > 0
}

// HasErrorClass appends an error class when fieldErrors exist.
func HasErrorClass(base string, fieldErrors []string, errorClass string) string {
	if FieldHasErrors(fieldErrors) {
		return base + " " + errorClass
	}
	return base
}

// JoinErrors concatenates error messages into one string.
func JoinErrors(fieldErrors []string) string {
	if FieldHasErrors(fieldErrors) {
		return strings.Join(fieldErrors, ", ")
	}
	return ""
}
