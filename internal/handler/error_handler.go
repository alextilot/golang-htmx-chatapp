package handler

import (
	"github.com/alextilot/golang-htmx-chatapp/internal/response"
	"github.com/alextilot/golang-htmx-chatapp/internal/validation"
	"github.com/alextilot/golang-htmx-chatapp/web/pages"
	"github.com/labstack/echo/v4"
	"log"
	"net/http"
)

// Define the custom error handler.
func (h *Handler) HTTPErrorHandler(err error, c echo.Context) {
	// Don't do anything if the response is already committed.
	if c.Response().Committed {
		return
	}

	// Extract the HTTP error code.
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}

	c.Logger().Error(err)

	// Render a custom error page based on the status code.
	switch code {
	case http.StatusNotFound:
		log.Println("Error serving 404")
		response.Send(c, response.Response{
			Status:       code,
			HTMLTemplate: pages.NotFoundPage(),
			Errors:       validation.NewFieldErrors(),
		})
	case http.StatusInternalServerError:
		log.Println("Error serving 500")
		response.Send(c, response.Response{
			Status:       code,
			HTMLTemplate: pages.ServerErrorPage(),
			Errors:       validation.NewFieldErrors(),
		})
	default:
		if err != nil {
			log.Println("Error serving generic error page:", err)
		}
	}
}
