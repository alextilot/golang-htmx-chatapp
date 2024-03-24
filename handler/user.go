package handler

import (
	"fmt"
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/model"
	"github.com/alextilot/golang-htmx-chatapp/services"
	"github.com/alextilot/golang-htmx-chatapp/utils"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/forms"
	"github.com/go-playground/validator/v10"

	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(etx echo.Context) error {
	req := &userLoginRequest{}
	errMap := utils.NewUIError()

	// Validate user inputs!
	if err := req.bind(etx); err != nil {
		errs := err.(validator.ValidationErrors)
		for _, v := range errs {
			errMap.Add(v.Field(), fmt.Sprintf("%v", v.Tag()))
		}
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(errMap))
	}

	user, err := h.userService.GetByUsername(req.Username)
	if err != nil {
		errMap.Add("other", "Invalid login information")
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(errMap))
	}
	if user != nil {
		errMap.Add("other", "Invalid login information")
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(errMap))
	}
	if !user.CheckPassword(req.Password) {
		errMap.Add("other", "Invalid login information")
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(errMap))
	}

	// JWT tokens for signed in users.
	if err := services.GenerateTokensAndSetCookies(user, etx); err != nil {
		return web.Render(etx, http.StatusUnauthorized, forms.LoginForm(errMap))
	}

	etx.Response().Header().Set("HX-Redirect", "/chatroom")
	return etx.String(http.StatusTemporaryRedirect, "Successful")
}

func (h *Handler) SignUp(etx echo.Context) error {
	var user = &model.User{}
	req := &userRegisterRequest{}

	// Validate user inputs!
	if err := req.bind(etx, user); err != nil {
		fmt.Println(utils.NewValidatorError(err))
		return etx.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	// Create user, this should error out if user already exists.
	if err := h.userService.Create(user); err != nil {
		fmt.Println(utils.NewError(err))
		return etx.JSON(http.StatusUnprocessableEntity, utils.NewError(err))
	}

	// JWT tokens for signed in users.
	if err := services.GenerateTokensAndSetCookies(user, etx); err != nil {
		out := utils.NewValidatorError(err)
		fmt.Println(out)
		component := forms.SignUpForm(utils.NewValidatorError(err))
		return web.Render(etx, http.StatusUnauthorized, component)
	}

	etx.Response().Header().Set("HX-Redirect", "/chatroom")
	return etx.String(http.StatusTemporaryRedirect, "Successful")
}
