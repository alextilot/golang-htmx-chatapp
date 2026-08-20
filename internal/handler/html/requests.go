package html

import "github.com/labstack/echo/v5"

// The types below are the website's own bind targets — dumb data carriers
// with no validation and no business meaning. Mapping them into a service's
// input type is this handler package's entire involvement in "shape."

type LoginRequest struct {
	Username string `form:"username" json:"username"`
	Password string `form:"password" json:"password"`
}

type SignupRequest struct {
	Username       string `form:"username" json:"username"`
	Email          string `form:"email" json:"email"`
	Password       string `form:"password" json:"password"`
	RepeatPassword string `form:"repeatPassword" json:"repeatPassword"`
}

type GroupCreateRequest struct {
	Name string `form:"name" json:"name"`
	Type string `form:"type" json:"type"`
}

type GroupUpdateRequest struct {
	Name        string `form:"name" json:"name"`
	Description string `form:"description" json:"description"`
}

type GroupAddMemberRequest struct {
	UserID string `form:"userID" json:"userID"`
}

// BindInput binds the current request into T using echo's content-negotiated
// binder (form or JSON). Sanitizing and validating T is the service's job,
// not this function's — see docs/error-handling.md.
func BindInput[T any](c *echo.Context) (T, error) {
	var input T
	err := c.Bind(&input)
	return input, err
}
