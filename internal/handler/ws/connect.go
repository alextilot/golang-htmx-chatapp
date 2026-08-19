package ws

import (
	"context"
	"net/http"

	"github.com/alextilot/golang-htmx-chatapp/internal/auth"
	"github.com/alextilot/golang-htmx-chatapp/internal/realtime"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v5"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// TODO: restrict to your actual origin(s) before going to production.
		return true
	},
}

// connect requires auth.Service.AuthMiddleware to have already run.
func connect(echoCtx *echo.Context, ctx context.Context, hub *realtime.Hub, groupID string) error {
	p := auth.PrincipalFromEcho(echoCtx)
	if !p.Authenticated {
		return echo.ErrUnauthorized
	}

	conn, err := upgrader.Upgrade(echoCtx.Response(), echoCtx.Request(), nil)
	if err != nil {
		return err
	}

	client := newClient(conn, hub, p.ID, p.Username, groupID)

	go client.ReadPump(echoCtx, ctx)
	go client.WritePump(echoCtx, ctx)

	return nil
}
