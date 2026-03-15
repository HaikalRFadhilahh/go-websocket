package data

import "github.com/gorilla/websocket"

type WebSocketClient struct {
	Name       string          `json:"name"`
	Connection *websocket.Conn `json:"connection"`
}
