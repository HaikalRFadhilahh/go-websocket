package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"websocket/data"
	"websocket/helper"

	"github.com/gorilla/websocket"
)

type webSocketHandler struct {
	upgrader        *websocket.Upgrader
	webSocketClient []data.WebSocketClient
}

type optsWebSocketHandler func(*webSocketHandler)

func (wsh *webSocketHandler) BroadcastMessage(w http.ResponseWriter, r *http.Request) {
	// CONFIG
	w.Header().Set("Content-Type", "application/json")

	// DATA REQUEST
	var broadcaseRequest data.BroadcastWebSocket

	// Json Binding
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&broadcaseRequest); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"statusCode": http.StatusBadRequest,
			"status":     "error",
			"message":    err.Error(),
		})
		return
	}

	if broadcaseRequest.Message == "" || broadcaseRequest.Username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"statusCode": http.StatusBadRequest,
			"status":     "error",
			"message":    "Username or Message must Filled",
		})
		return
	}

	if len(wsh.webSocketClient) < 1 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"statusCode": http.StatusNotFound,
			"status":     "error",
			"message":    "No One Users Connected to WebSocket",
		})
		return
	}

	var counterBroadcastMessage int
	for _, d := range wsh.webSocketClient {
		if err := d.Connection.WriteMessage(websocket.TextMessage, fmt.Appendf(nil, "%v - %v", broadcaseRequest.Username, broadcaseRequest.Message)); err != nil {
			helper.FilterStruct(wsh.webSocketClient, func(w data.WebSocketClient) bool {
				if w.Connection == d.Connection {
					return false
				}

				return true
			})
			continue
		}
		counterBroadcastMessage++
	}

	if counterBroadcastMessage < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]any{
			"statusCode": http.StatusInternalServerError,
			"status":     "error",
			"message":    "No One Recieved Broadcast Message",
		})
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"statusCode": http.StatusOK,
			"status":     "success",
			"message":    fmt.Sprintf("Message Success Broadcast to %v Users", counterBroadcastMessage),
		})
		return
	}
}

func (wsh *webSocketHandler) ClientWebSocket(w http.ResponseWriter, r *http.Request) {
	// WEB Socket
	ws, err := wsh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error(fmt.Sprintf("Error Upgrade Connection to Web Socket : %v", err.Error()))
		return
	}
	defer ws.Close()

	wsh.webSocketClient = append(wsh.webSocketClient, data.WebSocketClient{
		Connection: ws,
		Active:     true,
	})

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			wsh.webSocketClient = helper.FilterStruct(wsh.webSocketClient, func(d data.WebSocketClient) bool {
				if d.Connection == ws {
					return false
				} else {
					return true
				}
			})
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				slog.Error(fmt.Sprintf("Error Read Message Because Connection Closed : %v", err.Error()))
				return
			}
			slog.Error(fmt.Sprintf("Error Unexpedted : %v", err.Error()))
			return
		}

		err = ws.WriteMessage(websocket.TextMessage, fmt.Appendf(nil, "Message Recieved : %v", string(msg)))
		if err != nil {
			wsh.webSocketClient = helper.FilterStruct(wsh.webSocketClient, func(d data.WebSocketClient) bool {
				if d.Connection == ws {
					return false
				} else {
					return true
				}
			})
			return
		}

		fmt.Println(wsh.webSocketClient)
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
