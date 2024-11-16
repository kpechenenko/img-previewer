package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/kpechenenko/img-previewer/internal/downloader"
	"github.com/kpechenenko/img-previewer/internal/handler"
	"github.com/kpechenenko/img-previewer/internal/middleware"
	"github.com/kpechenenko/img-previewer/internal/previewer"
	"github.com/kpechenenko/img-previewer/internal/service"
)

// PreviewerApp веб сервер с api для создания превью изображений.
type PreviewerApp struct {
	srv *http.Server
}

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

func (a *PreviewerApp) Start() {
	go func() {
		slog.Info("starting previewer app at " + a.srv.Addr)
		if err := a.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("Failed to start previewer app", "error", err)
		}
	}()
}

func (a *PreviewerApp) Stop(ctx context.Context) error {
	slog.Info("stopping previewer app")
	if err := a.srv.Shutdown(ctx); err != nil {
		slog.Error("Failed to stop app app", "error", err)
		return err
	}
	return nil
}
