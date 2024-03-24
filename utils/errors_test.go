package utils

import (
	"fmt"
	"testing"
)

func TestErrors(t *testing.T) {
	output := NewValidatorError(nil)
	fmt.Println("value is", output.Errors)

	// name := "Gladys"
	// want := regexp.MustCompile(`\b` + name + `\b`)
	// msg, err := Hello("Gladys")
	// if !want.MatchString(msg) || err != nil {
	// 	t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	// }
}
