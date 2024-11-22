package service

import (
	"context"
	"errors"
	"image"
	"net/http"
	"strings"

	"github.com/kpechenenko/img-previewer/internal/downloader" //nolint:depguard
	"github.com/kpechenenko/img-previewer/internal/previewer"  //nolint:depguard
)

// HTTPPreviewerService создает превью для изображений, которые загружает по http.
type HTTPPreviewerService struct {
	previewer  previewer.Previewer
	downloader downloader.HTTPImageDownloader
}

// DownloadImageAndMakePreview скачать изображение и создать превью.
func (s *HTTPPreviewerService) DownloadImageAndMakePreview(
	ctx context.Context,
	downloadURL string,
	proxyHeader http.Header,
	width int,
	height int,
) (image.Image, error) {
	if !strings.HasPrefix(downloadURL, "http://") {
		return nil, errors.New("downloadURL must be starts with http://")
	}
	var img image.Image
	var err error
	if img, err = s.downloader.Download(ctx, downloadURL, proxyHeader); err != nil {
		return nil, err
	}
	var preview image.Image
	if preview, err = s.previewer.MakePreview(img, width, height); err != nil {
		return nil, err
	}
	return preview, nil
}

// NewHTTPPreviewerService конструктор с параметрами.
func NewHTTPPreviewerService(
	previewer previewer.Previewer,
	downloader downloader.HTTPImageDownloader,
) *HTTPPreviewerService {
	return &HTTPPreviewerService{
		previewer:  previewer,
		downloader: downloader,
	}
}
