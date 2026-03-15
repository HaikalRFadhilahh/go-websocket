package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"websocket/handler"

	"github.com/gorilla/websocket"
)

type apiServer struct {
	Host string
	Port string
	Url  string
}

type optsAPIServer func(*apiServer)

func NewAPIServer(opts ...optsAPIServer) *apiServer {
	// INIT SERVER
	server := &apiServer{
		Host: "127.0.0.1",
		Port: "8000",
		Url:  "http://127.0.0.1:8000",
	}

	// opts
	for _, f := range opts {
		f(server)
	}

	// rtn
	return server
}

func (s *apiServer) Init() {
	// DEPS INJ
	upgrader := &websocket.Upgrader{}
	webSocketHandler := handler.NewWebSocketHandler(handler.WithCustomUpgrader(upgrader))

	// ROUTER INIT
	r := http.NewServeMux()

	// ROUTING
	r.HandleFunc("GET /websocket", webSocketHandler.ClientWebSocket)
	r.HandleFunc("POST /broadcast", webSocketHandler.BroadcastMessage)

	// Server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%v:%v", s.Host, s.Port),
		Handler: r,
	}

	go s.listenAndServe(httpServer)

	shutdownChannel := make(chan os.Signal, 1)
	defer close(shutdownChannel)
	signal.Notify(shutdownChannel, syscall.SIGTERM, syscall.SIGINT)
	<-shutdownChannel

	shutdownContext, closeShutdownContext := context.WithTimeout(context.Background(), 5*time.Second)
	defer closeShutdownContext()

	if err := httpServer.Shutdown(shutdownContext); err != nil {
		slog.Error(fmt.Sprintf("Error while Shutdown Server : %v", err.Error()))
	} else {
		slog.Info("Success Shutdown Server Properly!")
	}
}

func (s *apiServer) listenAndServe(muxserver *http.Server) {
	slog.Info(fmt.Sprintf("Server Running On %v", s.Url))
	if err := muxserver.ListenAndServe(); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error(fmt.Sprintf("Error Server : %v", err.Error()))
			return
		}

		return
	}
}

func WithCustomHost(h string) optsAPIServer {
	return func(s *apiServer) {
		s.Host = h
	}
}

func WithCustomPort(p string) optsAPIServer {
	return func(s *apiServer) {
		s.Port = p
	}
}

func WithCustomUrl(u string) optsAPIServer {
	return func(s *apiServer) {
		s.Url = u
	}
}
