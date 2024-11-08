package handler

import (
	"log/slog"
	"net/http"
)

type PingHandler struct{}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		slog.Error("error to write response", "error", err)
		return
	}
}

func NewPingHandler() *PingHandler {
	return &PingHandler{}
}
