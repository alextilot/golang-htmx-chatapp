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
	view := &forms.LoginFormModelView{
		Username: nil,
		Password: nil,
		Other:    nil,
	}

	if ok := req.bind(etx, view); !ok {
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(view))
	}

	user, err := h.userService.GetByUsername(req.Username)
	if err != nil {
		view.Other = ErrorUsernameNotFound
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(view))
	}

	if user == nil {
		view.Other = ErrorInvalidCredentials
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(view))
	}
	if !user.CheckPassword(req.Password) {
		view.Other = ErrorInvalidCredentials
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(view))
	}

	// JWT tokens for signed in users.
	if err := services.GenerateTokensAndSetCookies(user, etx); err != nil {
		view.Other = ErrorUserAuthentication
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(view))
	}

	etx.Response().Header().Set("HX-Redirect", "/chatroom")
	return etx.String(http.StatusTemporaryRedirect, "Success")
}

func (h *Handler) SignUp(etx echo.Context) error {
	var user = &model.User{}
	req := &userRegisterRequest{}
	view := &forms.SignUpFormModelView{
		Username:       nil,
		Password:       nil,
		RepeatPassword: nil,
		Other:          nil,
	}

	// Validate user inputs!
	if ok := req.bind(etx, user, view); !ok {
		return web.Render(etx, http.StatusUnauthorized, forms.SignUpForm(view))
	}

	// Create user, this should error out if user already exists.
	if err := h.userService.Create(user); err != nil {
		view.Username = ErrorUsernameAlreadyExists
		return web.Render(etx, http.StatusUnauthorized, forms.SignUpForm(view))
	}

	// JWT tokens for signed in users.
	if err := services.GenerateTokensAndSetCookies(user, etx); err != nil {
		view.Other = ErrorUserAuthentication
		return web.Render(etx, http.StatusUnauthorized, forms.SignUpForm(view))
	}

	etx.Response().Header().Set("HX-Redirect", "/chatroom")
	return etx.String(http.StatusTemporaryRedirect, "Success")
}

func (h *Handler) Logout(etx echo.Context) error {
	services.RemoveTokensAndCookies(etx)
	etx.Response().Header().Set("HX-Redirect", "/")
	return etx.String(http.StatusTemporaryRedirect, "Success")
}
