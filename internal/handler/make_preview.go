package handler

import (
	"errors"
	"image"
	"image/jpeg"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/kpechenenko/img-previewer/internal/downloader"
	"github.com/kpechenenko/img-previewer/internal/previewer"
)

// MakePreviewHandler обработчик http запроса на создание превью изображения.
// Скачивает изображение с внешнего ресурса, создает превью, отдает результат пользователю.
type MakePreviewHandler struct {
	previewer  previewer.Previewer
	downloader downloader.HTTPImageDownloader
}

func NewMakePreviewHandler(
	previewer previewer.Previewer,
	downloader downloader.HTTPImageDownloader,
) *MakePreviewHandler {
	return &MakePreviewHandler{previewer: previewer, downloader: downloader}
}

const (
	PreviewPrefix = "/preview/"
)

type pathParams struct {
	width  int
	height int
	url    string
}

func (h *MakePreviewHandler) parsePathParams(p string) (*pathParams, error) {
	p = strings.TrimPrefix(p, PreviewPrefix)
	widthEndIdx := strings.Index(p, "/")
	if widthEndIdx == -1 {
		msg := "width is empty"
		slog.Info(msg)
		return nil, errors.New(msg)
	}
	width, err := strconv.Atoi(p[:widthEndIdx])
	if err != nil {
		msg := "width is invalid number"
		slog.Info(msg)
		return nil, errors.New(msg)
	}
	p = p[widthEndIdx+1:]
	heightEndIdx := strings.Index(p, "/")
	if heightEndIdx == -1 {
		msg := "height is empty"
		slog.Info(msg)
		return nil, errors.New(msg)
	}
	height, err := strconv.Atoi(p[:heightEndIdx])
	if err != nil {
		msg := "height is invalid number"
		slog.Info(msg)
		return nil, errors.New(msg)
	}
	p = p[heightEndIdx+1:]
	return &pathParams{
		width:  width,
		height: height,
		url:    p,
	}, nil
}

func (h *MakePreviewHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error
	var p *pathParams
	if p, err = h.parsePathParams(r.RequestURI); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			slog.Error("error writing response", "error", err)
		}
		return
	}
	slog.Info("try to make preview", "params", *p)
	var img image.Image
	if img, err = h.downloader.Download(r.Context(), "http://"+p.url, r.Header); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	var preview image.Image
	if preview, err = h.previewer.MakePreview(img, p.width, p.height); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
	if err = jpeg.Encode(w, preview, nil); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
