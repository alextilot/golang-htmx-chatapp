package html

import (
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
	"github.com/labstack/echo/v5"
)

// Profile renders the logged-in user's profile page.
// GET /profile
func (h *Handler) Profile(c *echo.Context) error {
	p := auth.PrincipalFromEcho(c)

	if _, err := h.Services.UserService.GetByID(c.Request().Context(), p.ID); err != nil {
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to load profile"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			Errors: fe,
		})
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
		fe := validation.FieldErrors{validation.FieldServer: {"Failed to load profile"}}
		return sendResponse(c, Response{
			Status: http.StatusInternalServerError,
			Errors: fe,
		})
	}

	return sendResponse(c, Response{
		Status: http.StatusOK,
		// TODO: Page: forms.ProfileEditForm(user),
	})
}
