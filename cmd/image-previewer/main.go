package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kpechenenko/img-previewer/internal/app"
)

var (
	previewerApp *app.PreviewerApp
	exit         = make(chan bool)
)

func main() {
	go listenSignals()
	go start()

	<-exit
	slog.Info("see you soon...")
}

func start() {
	var cfg *config
	var err error
	if cfg, err = readConfigFromEnv(); err != nil {
		slog.Error("read config from env", "error", err)
		os.Exit(1)
	}
	previewerApp = app.NewPreviewer(cfg.server.Addr)
	previewerApp.Start()
}

func listenSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGINT)

	for {
		s := <-signalCh

		switch s {
		case syscall.SIGHUP:
			slog.Info("SIGHUP received")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			err := previewerApp.Stop(ctx)
			if err != nil {
				slog.Error("can't stop previewer with SIGHUP", "error", err)
				os.Exit(1)
			}
			cancel()
			start()
		case syscall.SIGINT, syscall.SIGTERM:
			slog.Info("SIGTERM or SIGINT received")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			err := previewerApp.Stop(ctx)
			if err != nil {
				slog.Error("can't stop previewer with SIGTERM or SIGINT: %v", "error", err)
				os.Exit(1)
			}
			cancel()
			exit <- true
		}
	}
}
