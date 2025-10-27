package handler

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/response"
	"github.com/alextilot/golang-htmx-chatapp/internal/routes"
	"github.com/alextilot/golang-htmx-chatapp/internal/usercontext"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation/schema"
	"github.com/alextilot/golang-htmx-chatapp/web/forms"
	"github.com/labstack/echo/v4"
)

func (h *Handler) Login(c echo.Context) error {
	// 1. Get input, parse, sanitize, validate
	input, err := schema.HandleInput(c, schema.SanitizeUserLogin, schema.ValidateUserLogin)

	if err != nil {
		if ierr, ok := err.(*schema.InputError); ok {
			return response.Send(c, response.Response{
				Status: http.StatusBadRequest,
				View:   forms.LoginForm(ierr.Fields),
				Errors: ierr.Fields,
			})
		}
		return err // unexpected, let middleware handle
	}

	// 2. Check login information
	user, err := h.Services.UserService.Login(input.Username, input.Password)
	if user == nil || err != nil {
		fe := validation.FieldErrors{
			validation.FieldAuth: {"Invalid login information"},
		}

		return response.Send(c, response.Response{
			Status: http.StatusUnauthorized,
			View:   forms.LoginForm(fe),
			Errors: fe,
		},
		)
	}

	// 3. Generate JWT Tokens
	if err := auth.SetUserAuthContext(user, c); err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to generate JWT tokens"}}

		return response.Send(c, response.Response{
			Status: http.StatusInternalServerError,
			View:   forms.LoginForm(fe),
			Errors: fe,
		},
		)
	}

	// 4. Return Response
	return response.Send(c, response.Response{
		Status:   http.StatusOK,
		Data:     map[string]any{"user": user},
		Redirect: "/chatroom",
	})
}

func (h *Handler) SignUp(c echo.Context) error {
	// 1. Get input, parse, sanitize, validate
	input, err := schema.HandleInput(c, schema.SanitizeUserCreate, schema.ValidateUserCreate)
	if err != nil {
		if ierr, ok := err.(*schema.InputError); ok {
			return response.Send(c, response.Response{
				Status: http.StatusBadRequest,
				View:   forms.SignupForm(ierr.Fields),
				Errors: ierr.Fields,
			})
		}
		return err // unexpected, let middleware handle
	}

	// Delegate signup and additional business-level validation to service
	newUser, fe, err := h.Services.UserService.Signup(input.Username, input.Email, input.Password)
	if err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Error creating user"}}
		return response.Send(c, response.Response{
			Status: http.StatusInternalServerError,
			View:   forms.SignupForm(fe),
			Errors: fe,
		})
	}

	if len(fe) > 0 {
		return response.Send(c, response.Response{
			Status: http.StatusConflict,
			View:   forms.SignupForm(fe),
			Errors: fe,
		})
	}

	// 4. Generate JWT tokens
	if err := auth.SetUserAuthContext(newUser, c); err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to generate JWT tokens"}}
		return response.Send(c, response.Response{
			Status: http.StatusInternalServerError,
			View:   forms.SignupForm(fe),
			Errors: fe,
		})
	}

	// 5. Return success (handles HTML, JSON, HTMX)
	return response.Send(c, response.Response{
		Status:   http.StatusOK,
		Data:     map[string]any{"user": newUser},
		Redirect: "/chatroom",
	})
}

func (h *Handler) Logout(c echo.Context) error {
	// Clear auth cookies
	auth.Clear(c)
	// Reset the user context.
	usercontext.Set(c, usercontext.Default())

	return response.Send(c, response.Response{
		Status:   http.StatusSeeOther,
		Data:     map[string]any{"logout": true},
		Redirect: routes.Routes.HomePage.Path,
	})
}
