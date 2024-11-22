// Package handler содержит код http обработчиков запросов.
package handler

import (
	"errors"
	"image"
	"image/jpeg"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/kpechenenko/img-previewer/internal/downloader" //nolint:depguard
	"github.com/kpechenenko/img-previewer/internal/service"    //nolint:depguard
)

// MakePreviewHandler обработчик http запроса на создание превью изображения.
// Скачивает изображение с внешнего ресурса, создает превью, отдает результат пользователю.
type MakePreviewHandler struct {
	srv service.NetPreviewerService
}

// NewMakePreviewHandler конструктор с параметрами.
func NewMakePreviewHandler(
	srv service.NetPreviewerService,
) *MakePreviewHandler {
	return &MakePreviewHandler{srv: srv}
}

const (
	// PreviewPrefix название эндпоинта для создания превью изображений.
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
		return nil, &InvalidParamErr{description: msg}
	}
	width, err := strconv.Atoi(p[:widthEndIdx])
	if err != nil {
		msg := "width is invalid number"
		slog.Info(msg)
		return nil, &InvalidParamErr{description: msg}
	}
	p = p[widthEndIdx+1:]
	heightEndIdx := strings.Index(p, "/")
	if heightEndIdx == -1 {
		msg := "height is empty"
		slog.Info(msg)
		return nil, &InvalidParamErr{description: msg}
	}
	height, err := strconv.Atoi(p[:heightEndIdx])
	if err != nil {
		msg := "height is invalid number"
		slog.Info(msg)
		return nil, &InvalidParamErr{description: msg}
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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slog.Info("try to make preview", "params", *p)
	var preview image.Image
	if preview, err = h.srv.DownloadImageAndMakePreview(
		r.Context(),
		"http://"+p.url,
		r.Header,
		p.width,
		p.height,
	); err != nil {
		var downloadImgErr *downloader.FailToDownloadImageErr
		if errors.As(err, &downloadImgErr) {
			http.Error(w, err.Error(), http.StatusBadGateway)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if err = jpeg.Encode(w, preview, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/jpeg")
}
