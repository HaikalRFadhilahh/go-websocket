package main

import (
	"log/slog"
	"os"
	"websocket/internal/server"

	"github.com/joho/godotenv"
)

func main() {
	// LOAD ENV
	if err := godotenv.Load(); err != nil {
		slog.Error("Error Load .env Files!")
	}

	// Init Server
	s := server.NewAPIServer(
		server.WithCustomHost(os.Getenv("APP_HOST")),
		server.WithCustomPort(os.Getenv("APP_PORT")),
		server.WithCustomUrl(os.Getenv("APP_URL")),
	)

	// RUN Server
	s.Init()
}
