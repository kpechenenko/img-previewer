// Package app содержит код приложения превью изображений.
package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/kpechenenko/img-previewer/internal/downloader" //nolint:depguard
	"github.com/kpechenenko/img-previewer/internal/handler"    //nolint:depguard
	"github.com/kpechenenko/img-previewer/internal/middleware" //nolint:depguard
	"github.com/kpechenenko/img-previewer/internal/previewer"  //nolint:depguard
	"github.com/kpechenenko/img-previewer/internal/service"    //nolint:depguard
)

// PreviewerApp веб сервер с api для создания превью изображений.
type PreviewerApp struct {
	srv *http.Server
}

// NewPreviewer конструктор с параметрами.
func NewPreviewer(addr string) *PreviewerApp {
	mux := http.NewServeMux()
	mux.Handle(
		handler.PreviewPrefix,
		handler.NewMakePreviewHandler(
			service.NewHTTPPreviewerService(
				previewer.NewKNNImageCompressor(),
				downloader.NewJPEGImageDownloader(http.DefaultClient),
			),
		),
	)
	mux.Handle("GET /ping", handler.NewPingHandler())
	loggedMux := middleware.NewRequestLogger(mux)
	srv := http.Server{Addr: addr, Handler: loggedMux, ReadTimeout: 5 * time.Second}
	return &PreviewerApp{srv: &srv}
}

// Start запустить приложение.
func (a *PreviewerApp) Start() {
	go func() {
		slog.Info("starting previewer app at " + a.srv.Addr)
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start previewer app", "error", err)
		}
	}()
}

// Stop остановить работу приложения.
func (a *PreviewerApp) Stop(ctx context.Context) error {
	slog.Info("stopping previewer app")
	if err := a.srv.Shutdown(ctx); err != nil {
		slog.Error("Failed to stop app app", "error", err)
		return err
	}
	return nil
}
