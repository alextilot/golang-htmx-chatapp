package handler

import (
	"errors"
	"fmt"

	"github.com/alextilot/golang-htmx-chatapp/model"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type userLoginRequest struct {
	Username string `form:"username" validate:"required,min=2,max=20"`
	Password string `form:"password" validate:"required,min=5"`
}

func (r *userLoginRequest) bind(c echo.Context) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	return nil
}

type userRegisterRequest struct {
	Username       string `form:"username" validate:"required,min=5,max=20"`
	Password       string `form:"password" validate:"required,min=5,eqfield=RepeatPassword"`
	RepeatPassword string `form:"repeatPassword" validate:"required,min=5,eqfield=Password"`
}

var (
	ErrorDuringPasswordHashing = errors.New("TODO: this is really bad.")
)

func msgForTag(tag string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "max":
		return "Maximum of 20 characters"
	case "min":
		return "Minimum of 5 characters"
	case "eqfield":
		return "Passwords must match"
	}
	return tag
}

func (r *userRegisterRequest) bind(c echo.Context, u *model.User) map[string]error {
	out := make(map[string]error)

	if err := c.Bind(r); err != nil {
		out["Other"] = err
		return out
	}

	if err := c.Validate(r); err != nil {

		validator := err.(validator.ValidationErrors)
		for _, ve := range validator {
			fmt.Printf("%v - %v\n", ve.Field(), ve.Tag())

			key := ve.Field()
			str := msgForTag(ve.Tag())
			out[key] = errors.New(str)
		}
		return out
	}

	u.Username = r.Username
	h, err := u.HashPassword(r.Password)
	if err != nil {
		out["Other"] = ErrorDuringPasswordHashing
		return out
	}
	u.Password = h
	return nil
}
