package core

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func Validator() *validator.Validate {
	return validate
}

func ValidateStruct(s any) FieldErrors {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	errors := make(FieldErrors)
	for _, e := range err.(validator.ValidationErrors) {
		msg := GetErrorMessage(e.Tag(), e.Field(), e.Param())

		if e.Tag() == "eqfield" && e.Param() == "Password" {
			msg = "Passwords do not match."
		}

		errors[e.Field()] = append(errors[e.Field()], msg)
	}

	return errors
}
