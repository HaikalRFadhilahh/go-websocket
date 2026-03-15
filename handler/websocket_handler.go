package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"websocket/data"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader        *websocket.Upgrader
	webSocketClient []data.WebSocketClient
}

type optsWebSocketHandler func(*webSocketHandler)

func (wsh *webSocketHandler) BroadcastMessage(w http.ResponseWriter, r *http.Request) {

}

func (wsh *webSocketHandler) ClientWebSocket(w http.ResponseWriter, r *http.Request) {
	// WEB Socket
	ws, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Error Upgrade Connection to Web Socket : %v", err.Error()))
		return
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				slog.Error(fmt.Sprintf("Error Read Message Because Connection Closed : %v", err.Error()))
				return
			}

			slog.Error(fmt.Sprintf("Error Unexpedted : %v", err.Error()))
			return
		}

		err = ws.WriteMessage(websocket.TextMessage, fmt.Appendf(nil, "Message Recieved : %v", string(msg)))
		if err != nil {
			
			return
		}
	}
}

func NewWebSocketHandler(opts ...optsWebSocketHandler) *webSocketHandler {
	// Web Socket Handler
	wsh := &webSocketHandler{
		upgrader: nil,
	}

	// Implementing Opts
	for _, f := range opts {
		f(wsh)
	}

	// rtn
	return wsh
}

func WithCustomUpgrader(u *websocket.Upgrader) optsWebSocketHandler {
	return func(wsh *webSocketHandler) {
		wsh.upgrader = u
	}
}

func WithCustomWebSocketClient(c []data.WebSocketClient) optsWebSocketHandler {
	return func(wsh *webSocketHandler) {
		wsh.webSocketClient = c
	}
}
