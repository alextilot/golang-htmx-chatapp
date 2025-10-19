package schema

import (
	"fmt"

	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
	"github.com/labstack/echo/v4"
)

type InputError struct {
	Kind   string
	Fields validation.FieldErrors
}

func (e *InputError) Error() string {
	return fmt.Sprintf("input error (%s)", e.Kind)
}

func HandleInput[T any](
	c echo.Context,
	sanitizer func(T) T,
	validator func(T) validation.FieldErrors,
) (T, error) {
	var input T
	if err := c.Bind(&input); err != nil {
		return input, &InputError{
			Kind: validation.FieldRequest,
			Fields: validation.FieldErrors{
				validation.FieldRequest: {err.Error()},
			},
		}
	}

	if sanitizer != nil {
		input = sanitizer(input)
	}

	if validator != nil {
		if errs := validator(input); errs != nil {
			return input, &InputError{
				Kind:   validation.FieldValidation,
				Fields: errs,
			}
		}
	}

	return input, nil
}
