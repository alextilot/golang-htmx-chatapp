package handler

import (
	"context"
	"github.com/alextilot/golang-htmx-chatapp/services"
	"github.com/alextilot/golang-htmx-chatapp/web"
	"github.com/alextilot/golang-htmx-chatapp/web/forms"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func PostLogin(etx echo.Context, ctx context.Context, userService *services.UserService) error {
	time.Sleep(1 * time.Second)
	username := etx.FormValue("username")
	password := etx.FormValue("password")

	var errorMessages []string

	// Validate user input
	validate := validator.New()
	err := validate.Var(username, "required,min=2,max=20")
	if err != nil {
		errorMessages = append(errorMessages, "Username is Required, minimum: 2, maximum 20")
	}

	err = validate.Var(password, "required,min=5,max=20")
	if err != nil {
		errorMessages = append(errorMessages, "Password is Required, minimum: 5, maximum 20")
	}

	if len(errorMessages) != 0 {
		component := forms.LoginForm(strings.Join(errorMessages, "\n"))
		return web.Render(etx, http.StatusUnauthorized, component)
	}

	// Check login information
	loggedInUser, err := userService.LoginUser(username, password)
	if loggedInUser == nil || err != nil {
		component := forms.LoginForm("Invalid login information")
		return web.Render(etx, http.StatusUnauthorized, component)
	}

	// JWT tokens for signed in users.
	err = services.GenerateTokensAndSetCookies(loggedInUser, etx)
	if err != nil {
		component := forms.LoginForm("Unexpected Error: JwtToken failed to generate")
		return web.Render(etx, http.StatusUnauthorized, component)
	}

	etx.Response().Header().Set("HX-Redirect", "/chatroom")
	return etx.String(http.StatusTemporaryRedirect, "Successful")
}

func PostSignup(etx echo.Context, ctx context.Context, userService *services.UserService) error {
	time.Sleep(1 * time.Second)
	username := etx.FormValue("username")
	password := etx.FormValue("password")
	repeatPassword := etx.FormValue("repeatPassword")

	var errorMessages []string

	// Validate input data
	validate := validator.New()
	err := validate.Var(username, "required,min=2,max=20")
	if err != nil {
		errorMessages = append(errorMessages, "Username required, minimum: 2, maximum 20")
	}

	err = validate.Var(password, "required,min=5,max=20")
	if err != nil {
		errorMessages = append(errorMessages, "Password required, minimum: 5, maximum 20")
	}

	if password != repeatPassword {
		errorMessages = append(errorMessages, "Passwords do not match")
	}

	if len(errorMessages) != 0 {
		component := forms.SignupForm(strings.Join(errorMessages, "\n"))
		return web.Render(etx, http.StatusUnauthorized, component)
	}

	// validate unique username
	users, err := userService.GetUsers(username)
	if err != nil || len(users) > 0 {
		component := forms.SignupForm("User with that name already exists")
		return web.Render(etx, http.StatusUnauthorized, component)
	}

	// create user
	newUser, err := userService.CreateUser(username, password)
	if err != nil {
		component := forms.SignupForm("Error creating user")
		return web.Render(etx, http.StatusUnauthorized, component)
	}

	// JWT tokens for signed in users.
	err = services.GenerateTokensAndSetCookies(newUser, etx)
	if err != nil {
		component := forms.SignupForm("Unexpected Error: JwtToken failed to generate")
		return web.Render(etx, http.StatusUnauthorized, component)
	}

	etx.Response().Header().Set("HX-Redirect", "/chatroom")
	return etx.String(http.StatusTemporaryRedirect, "Successful")
}
