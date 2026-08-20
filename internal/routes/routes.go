// Package routes centralizes named paths for the website interface so
// templates and handlers don't hardcode URL strings in more than one place.
//
// This only covers internal website links (e.g. "Go Home" buttons on error
// pages) — it is not a router and does not replace internal/router, which
// remains the single place routes are registered.
package routes

// Route is a single named website path.
type Route struct {
	Path string
}

type routeTable struct {
	HomePage   Route
	LoginPage  Route
	SignupPage Route
	AboutPage  Route
	Profile    Route
	Logout     Route
	Groups     Route
}

// Routes is the shared table of named website paths.
var Routes = routeTable{
	HomePage:   Route{Path: "/"},
	LoginPage:  Route{Path: "/login"},
	SignupPage: Route{Path: "/signup"},
	AboutPage:  Route{Path: "/about"},
	Profile:    Route{Path: "/profile"},
	Logout:     Route{Path: "/logout"},
	Groups:     Route{Path: "/groups"},
}
