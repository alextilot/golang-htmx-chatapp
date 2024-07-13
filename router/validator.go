package router

import "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {
	// validate := validator.New(validator.WithRequiredStructEnabled())
	validate := validator.New()
	return &Validator{
		validator: validate,
	}
}

func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
