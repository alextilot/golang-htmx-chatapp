package handler

import (
	"net/http"

	"errors"
	"github.com/alextilot/golang-htmx-chatapp/model"
	"github.com/alextilot/golang-htmx-chatapp/services"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/forms"
	"github.com/labstack/echo/v4"
)

var (
	ErrorUsernameNotFound   = errors.New("Username not found.")
	ErrorInvalidCredentials = errors.New("Incorrect username or password.")
	ErrorUserAuthentication = errors.New("Unable to authenticate user.")

	ErrorUsernameAlreadyExists = errors.New("Username already exists.")
)

func (h *Handler) Login(etx echo.Context) error {
	req := &userLoginRequest{}
	errs := forms.LoginErrorsModelView{
		Username: nil,
		Password: nil,
		Other:    nil,
	}

	// Validate user inputs!
	// if err := req.bind(etx); err != nil {
	// 	validation := err.(validator.ValidationErrors)
	// 	for _, _ := range validation {
	// 		//errs.Add(v.Field(), fmt.Sprintf("%v", v.Tag()))
	// 	}
	// 	return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(errs))
	// }

	user, err := h.userService.GetByUsername(req.Username)
	if err != nil {
		errs.Other = ErrorUsernameNotFound
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(errs))
	}

	if user == nil {
		errs.Other = ErrorInvalidCredentials
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(errs))
	}
	if !user.CheckPassword(req.Password) {
		errs.Other = ErrorInvalidCredentials
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(errs))
	}

	// JWT tokens for signed in users.
	if err := services.GenerateTokensAndSetCookies(user, etx); err != nil {
		errs.Other = ErrorUserAuthentication
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(errs))
	}

	etx.Response().Header().Set("HX-Redirect", "/chatroom")
	return etx.String(http.StatusTemporaryRedirect, "Success")
}

func (h *Handler) SignUp(etx echo.Context) error {
	var user = &model.User{}
	req := &userRegisterRequest{}
	errs := forms.SignUpErrorsModelView{
		Username:       nil,
		Password:       nil,
		RepeatPassword: nil,
		Other:          nil,
	}

	// Validate user inputs!
	if eMap := req.bind(etx, user); eMap != nil {
		for key, value := range eMap {
			if key == "Username" {
				errs.Username = value
			} else if key == "Password" {
				errs.Password = value
			} else if key == "RepeatPassword" {
				errs.RepeatPassword = value
			} else {
				errs.Other = value
			}
		}
		return web.Render(etx, http.StatusUnauthorized, forms.SignUpForm(errs))
	}

	// Create user, this should error out if user already exists.
	if err := h.userService.Create(user); err != nil {
		errs.Username = ErrorUsernameAlreadyExists
		return web.Render(etx, http.StatusUnauthorized, forms.SignUpForm(errs))
	}

	// JWT tokens for signed in users.
	if err := services.GenerateTokensAndSetCookies(user, etx); err != nil {
		errs.Other = ErrorUserAuthentication
		return web.Render(etx, http.StatusUnauthorized, forms.SignUpForm(errs))
	}

	etx.Response().Header().Set("HX-Redirect", "/chatroom")
	return etx.String(http.StatusTemporaryRedirect, "Success")
}

func (h *Handler) Logout(etx echo.Context) error {
	services.RemoveTokensAndCookies(etx)
	etx.Response().Header().Set("HX-Redirect", "/")
	return etx.String(http.StatusTemporaryRedirect, "Success")
}
