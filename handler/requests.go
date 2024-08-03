package handler

import (
	"errors"
	"fmt"

	"github.com/alextilot/golang-htmx-chatapp/model"
	"github.com/alextilot/golang-htmx-chatapp/web/forms"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

var (
	//TODO: Change password hashing error
	ErrorDuringPasswordHashing = errors.New("Password uses invalid characters.")
	ErrorUsernameRequired      = errors.New("Username is required.")
	ErrorUsernameMin           = errors.New("Username minimum is 2 characters.")
	ErrorUsernameMax           = errors.New("Username maximum is 20 characters.")
	ErrorPasswordRequired      = errors.New("Password is required.")
	ErrorPasswordMin           = errors.New("Password minimum is 5 characters.")
	ErrorPasswordEqual         = errors.New("Passwords do not match.")
)

func msgForUsernameError(tag string, msg string) error {
	switch tag {
	case "required":
		return ErrorUsernameRequired
	case "min":
		return ErrorUsernameMin
	case "max":
		return ErrorUsernameMax
	}
	return errors.New(msg)
}

func msgForPasswordError(tag string, msg string) error {
	switch tag {
	case "required":
		return ErrorPasswordRequired
	case "min":
		return ErrorPasswordMin
	case "eqfield":
		return ErrorPasswordEqual
	}
	return errors.New(msg)
}

type userLoginRequest struct {
	Username string `form:"username" validate:"required,min=2,max=20"`
	Password string `form:"password" validate:"required,min=5"`
}

func (r *userLoginRequest) bind(c echo.Context, v *forms.LoginFormModelView) bool {
	if err := c.Bind(r); err != nil {
		v.Other = err
		return false
	}
	if err := c.Validate(r); err != nil {

		validator := err.(validator.ValidationErrors)
		for _, ve := range validator {
			fmt.Printf("%v - %v\n", ve.Field(), ve.Tag())

			key := ve.Field()
			tag := ve.Tag()
			msg := ve.Error()

			if key == "Username" {
				v.Username = msgForUsernameError(tag, msg)
			} else if key == "Password" {
				v.Password = msgForPasswordError(tag, msg)
			} else {
				v.Other = errors.New(msg)
			}
		}
		return false
	}
	return true
}

type userRegisterRequest struct {
	Username       string `form:"username" validate:"required,min=2,max=20"`
	Password       string `form:"password" validate:"required,min=5,eqfield=RepeatPassword"`
	RepeatPassword string `form:"repeatPassword" validate:"required,min=5,eqfield=Password"`
}

func (r *userRegisterRequest) bind(c echo.Context, u *model.User, v *forms.SignUpFormModelView) bool {

	if err := c.Bind(r); err != nil {
		v.Other = err
		return false
	}

	if err := c.Validate(r); err != nil {

		validator := err.(validator.ValidationErrors)
		for _, ve := range validator {
			fmt.Printf("%v - %v\n", ve.Field(), ve.Tag())

			key := ve.Field()
			tag := ve.Tag()
			msg := ve.Error()

			if key == "Username" {
				v.Username = msgForUsernameError(tag, msg)
			} else if key == "Password" {
				v.Password = msgForPasswordError(tag, msg)
			} else if key == "RepeatPassword" {
				v.RepeatPassword = msgForPasswordError(tag, msg)
			} else {
				v.Other = errors.New(msg)
			}
		}
		return false
	}

	u.Username = r.Username
	h, err := u.HashPassword(r.Password)
	if err != nil {
		v.Other = ErrorDuringPasswordHashing
		return false
	}
	u.Password = h
	return true
}

type messagesRequest struct {
	Count int `param:"count" validate:"required,numeric"`
	Ref   int `query:"ref" validate:"required,numeric"`
}

var (
	ErrorFieldIsRequired    = errors.New("Missing required field.")
	ErrorParamMustBeNumeric = errors.New("Param must be numeric.")
)

func errorForMessages(key string, tag string, msg string) error {
	switch tag {
	case "required":
		return ErrorFieldIsRequired
	case "numeric":
		return ErrorParamMustBeNumeric
	}
	return errors.New(msg)
}

func (r *messagesRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {

		validator := err.(validator.ValidationErrors)
		for _, ve := range validator {
			key := ve.Field()
			tag := ve.Tag()
			msg := ve.Error()

			switch key {
			case "Count":
				return errorForMessages(key, tag, msg)
			case "Ref":
				return errorForMessages(key, tag, msg)
			default:
				return errors.New(msg)
			}
		}
	}
	return nil
}
