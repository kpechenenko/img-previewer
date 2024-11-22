package handler

import (
	"log/slog"
	"net/http"
)

// PingHandler обработчик пинг запросов от пользователя.
type PingHandler struct{}

func (h *PingHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	if _, err := w.Write([]byte("pong")); err != nil {
		slog.Error("error to write response", "error", err)
		return
	}
}

// NewPingHandler конструктор с параметрами.
func NewPingHandler() *PingHandler {
	return &PingHandler{}
}
