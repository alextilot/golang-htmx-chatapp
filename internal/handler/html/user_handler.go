package html

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/apperr"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/routes"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
	"github.com/alextilot/golang-htmx-chatapp/web/forms"
	"github.com/labstack/echo/v5"
)

func (h *Handler) Login(c *echo.Context) error {
	// 1. Bind
	req, err := BindInput[LoginRequest](c)
	if err != nil {
		fe := apperr.FieldErrors{apperr.FieldGlobal: {"We couldn't read that request."}}
		return sendResponse(c, Response{
			Status: http.StatusBadRequest,
			Page:   forms.LoginForm(fe, forms.LoginFormValues{}),
			Errors: fe,
		})
	}
	values := forms.LoginFormValues{Username: req.Username}

	// 2. Call the service — it validates and checks credentials internally
	user, err := h.Services.UserService.Login(c.Request().Context(), service.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})
	if appErr, ok := apperr.As(err); ok {
		fe := fieldsOf(appErr)
		return sendResponse(c, Response{
			Status: statusFor(appErr.Kind),
			Page:   forms.LoginForm(fe, values),
			Errors: fe,
		})
	}
	if err != nil {
		return err // unexpected, let middleware handle
	}

	// 3. Generate JWT tokens and attach the authenticated principal
	if err := h.Auth.CreateUserSession(user, c); err != nil {
		fe := apperr.FieldErrors{apperr.FieldGlobal: {"Failed to generate session"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			Page:   forms.LoginForm(fe, values),
			Errors: fe,
		})
	}

	// 4. Redirect to the app
	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: routes.Routes.Groups.Path,
	})
}

func (h *Handler) Logout(c *echo.Context) error {
	h.Auth.Clear(c)
	auth.SetEchoPrincipal(c, auth.AnonymousPrincipal())

	return sendResponse(c, Response{
		Status:   http.StatusSeeOther,
		Redirect: routes.Routes.HomePage.Path,
	})
}

func (h *Handler) SignUp(c *echo.Context) error {
	// 1. Bind
	req, err := BindInput[SignupRequest](c)
	if err != nil {
		fe := apperr.FieldErrors{apperr.FieldGlobal: {"We couldn't read that request."}}
		return sendResponse(c, Response{
			Status: http.StatusBadRequest,
			Page:   forms.SignupForm(fe, forms.SignupFormValues{}),
			Errors: fe,
		})
	}
	values := forms.SignupFormValues{Username: req.Username, Email: req.Email}

	// 2. Delegate signup — sanitizing, validation, and business rules all
	// happen inside the service.
	user, err := h.Services.UserService.Signup(c.Request().Context(), service.SignupInput{
		Username:       req.Username,
		Email:          req.Email,
		Password:       req.Password,
		RepeatPassword: req.RepeatPassword,
	})
	if appErr, ok := apperr.As(err); ok {
		fe := fieldsOf(appErr)
		return sendResponse(c, Response{
			Status: statusFor(appErr.Kind),
			Page:   forms.SignupForm(fe, values),
			Errors: fe,
		})
	}
	if err != nil {
		return err // unexpected, let middleware handle
	}

	// 3. Generate JWT tokens and attach the authenticated principal
	if err := h.Auth.CreateUserSession(user, c); err != nil {
		fe := apperr.FieldErrors{apperr.FieldGlobal: {"Failed to generate session"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			Page:   forms.SignupForm(fe, values),
			Errors: fe,
		})
	}

	// 4. Redirect to the app
	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: routes.Routes.Groups.Path,
	})
}
