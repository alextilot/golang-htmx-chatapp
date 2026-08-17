package html

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/routes"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation/schema"
	"github.com/alextilot/golang-htmx-chatapp/web/forms"
	"github.com/labstack/echo/v5"
)

func (h *Handler) Login(c *echo.Context) error {
	// 1. Bind, sanitize, validate
	input, err := schema.HandleInput(c, schema.SanitizeUserLogin, schema.ValidateUserLogin)
	if err != nil {
		if ierr, ok := err.(*schema.InputError); ok {
			return sendResponse(c, Response{
				Status: http.StatusBadRequest,
				Page:   forms.LoginForm(ierr.Fields),
				Errors: ierr.Fields,
			})
		}
		return err // unexpected, let middleware handle
	}

	// 2. Check login information
	user, err := h.Services.UserService.Login(c.Request().Context(), input.Username, input.Password)
	if err != nil {
		fe := validation.FieldErrors{
			validation.FieldAuth: {"Invalid login information"},
		}
		return sendResponse(c, Response{
			Status: http.StatusUnauthorized,
			Page:   forms.LoginForm(fe),
			Errors: fe,
		})
	}

	// 3. Generate JWT tokens and attach the authenticated principal
	if err := h.Auth.CreateUserSession(user, c); err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to generate session"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			Page:   forms.LoginForm(fe),
			Errors: fe,
		})
	}

	// 4. Redirect to the app
	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: "/chatroom",
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
	// 1. Bind, sanitize, validate
	input, err := schema.HandleInput(c, schema.SanitizeUserCreate, schema.ValidateUserCreate)
	if err != nil {
		if ierr, ok := err.(*schema.InputError); ok {
			return sendResponse(c, Response{
				Status: http.StatusBadRequest,
				Page:   forms.SignupForm(ierr.Fields),
				Errors: ierr.Fields,
			})
		}
		return err // unexpected, let middleware handle
	}

	// 2. Delegate signup and additional business-level validation to the service
	user, fe, err := h.Services.UserService.Signup(
		c.Request().Context(),
		input.Username,
		input.Email,
		input.Password,
	)
	if err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Error creating user"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			Page:   forms.SignupForm(fe),
			Errors: fe,
		})
	}

	if len(fe) > 0 {
		return sendResponse(c, Response{
			Status: http.StatusConflict,
			Page:   forms.SignupForm(fe),
			Errors: fe,
		})
	}

	// 3. Generate JWT tokens and attach the authenticated principal
	if err := h.Auth.CreateUserSession(user, c); err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to generate session"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			Page:   forms.SignupForm(fe),
			Errors: fe,
		})
	}

	// 4. Redirect to the app
	return sendResponse(c, Response{
		Status:   http.StatusOK,
		Redirect: "/chatroom",
	})
}
