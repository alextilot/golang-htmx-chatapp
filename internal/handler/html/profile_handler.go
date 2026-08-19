package html

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/apperr"
	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/labstack/echo/v5"
)

// Profile renders the logged-in user's profile page.
// GET /profile
func (h *Handler) Profile(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	if _, err := h.Services.UserService.GetByID(c.Request().Context(), p.ID); err != nil {
		if appErr, ok := apperr.As(err); ok {
			return sendResponse(c, Response{Status: statusFor(appErr.Kind)})
		}
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status: http.StatusOK,
		// TODO: Page: pages.ProfilePage(user),
	})
}

// ProfileEdit renders the profile edit form.
// GET /profile/edit
func (h *Handler) ProfileEdit(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	if _, err := h.Services.UserService.GetByID(c.Request().Context(), p.ID); err != nil {
		if appErr, ok := apperr.As(err); ok {
			return sendResponse(c, Response{Status: statusFor(appErr.Kind)})
		}
		return err // unexpected — let middleware handle
	}

	return sendResponse(c, Response{
		Status: http.StatusOK,
		// TODO: Page: forms.ProfileEditForm(user),
	})
}
