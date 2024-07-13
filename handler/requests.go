package handler

import (
	"github.com/alextilot/golang-htmx-chatapp/model"
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
	Username       string `form:"username" validate:"required,min=2,max=20"`
	Password       string `form:"password" validate:"required,min=5,eqfield=RepeatPassword"`
	RepeatPassword string `form:"repeatPassword" validate:"required,min=5,eqfield=Password"`
}

func (r *userRegisterRequest) bind(c echo.Context, u *model.User) error {
	if err := c.Bind(r); err != nil {
		return err
	}
	if err := c.Validate(r); err != nil {
		return err
	}
	u.Username = r.Username
	h, err := u.HashPassword(r.Password)
	if err != nil {
		return err
	}
	u.Password = h
	return nil
}
