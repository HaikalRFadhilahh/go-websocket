package data

import "github.com/gorilla/websocket"

type WebSocketClient struct {
	Connection *websocket.Conn `json:"connection"`
	Active     bool            `json:"active"`
}

type BroadcastWebSocket struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}
