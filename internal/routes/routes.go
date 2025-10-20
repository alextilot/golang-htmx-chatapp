package routes

import "net/http"

type Role int

const (
	Public Role = iota
	Authenticated
)

type Route struct {
	Path        string
	Method      string
	Roles       []Role
	Description string
}

var Routes = struct {
	// Authentication
	LoginPage  Route
	SignupPage Route
	Logout     Route

	// User profile
	Profile     Route
	ProfileEdit Route

	// Chat / Chat rooms
	ChatList    Route
	ChatRoom    Route
	ChatJoin    Route
	ChatLeave   Route
	ChatMessage Route
	ChatSocket  Route

	// Global / Misc
	HomePage    Route
	AboutPage   Route
	NotFound    Route
	ServerError Route
}{
	// Authentication
	LoginPage: Route{
		Path:        "/login",
		Method:      http.MethodGet,
		Roles:       []Role{Public},
		Description: "Show login page",
	},
	SignupPage: Route{
		Path:        "/signup",
		Method:      http.MethodGet,
		Roles:       []Role{Public},
		Description: "Show signup page",
	},
	Logout: Route{
		Path:        "/logout",
		Method:      http.MethodPost,
		Roles:       []Role{Authenticated},
		Description: "Log out the current user",
	},

	// User profile
	Profile: Route{
		Path:        "/profile",
		Method:      http.MethodGet,
		Roles:       []Role{Authenticated},
		Description: "Show user profile",
	},
	ProfileEdit: Route{
		Path:        "/profile/edit",
		Method:      http.MethodGet,
		Roles:       []Role{Authenticated},
		Description: "Edit user profile page",
	},

	// Chat / Chat rooms
	ChatList: Route{
		Path:        "/chat",
		Method:      http.MethodGet,
		Roles:       []Role{Authenticated},
		Description: "List of all chat rooms",
	},
	ChatRoom: Route{
		Path:        "/chat/:id",
		Method:      http.MethodGet,
		Roles:       []Role{Authenticated},
		Description: "Show messages in a specific chat room",
	},
	ChatJoin: Route{
		Path:        "/chat/:id/join",
		Method:      http.MethodPost,
		Roles:       []Role{Authenticated},
		Description: "Join a chat room",
	},
	ChatLeave: Route{
		Path:        "/chat/:id/leave",
		Method:      http.MethodPost,
		Roles:       []Role{Authenticated},
		Description: "Leave a chat room",
	},
	ChatMessage: Route{
		Path:        "/chat/:id/message",
		Method:      http.MethodPost,
		Roles:       []Role{Authenticated},
		Description: "Post a new message in a chat room",
	},
	ChatSocket: Route{
		Path:        "/ws/chat/:id",
		Method:      http.MethodGet,
		Roles:       []Role{Authenticated},
		Description: "WebSocket endpoint for real-time chat messages",
	},

	// Global / Misc
	HomePage: Route{
		Path:        "",
		Method:      http.MethodGet,
		Roles:       []Role{Public},
		Description: "Landing page or redirect to chat",
	},
	AboutPage: Route{
		Path:        "/about",
		Method:      http.MethodGet,
		Roles:       []Role{Public},
		Description: "About or help page",
	},
	NotFound: Route{
		Path:        "/404",
		Method:      http.MethodGet,
		Roles:       []Role{Public},
		Description: "Custom 404 page",
	},
	ServerError: Route{
		Path:        "/500",
		Method:      http.MethodGet,
		Roles:       []Role{Public},
		Description: "Custom 500 server error page",
	},
}
