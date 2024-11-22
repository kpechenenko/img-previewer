// Package service бизнес логика приложения.
package service

import (
	"context"
	"image"
	"net/http"
)

// NetPreviewerService создает превью для изображений, которые загружает из сети.
type NetPreviewerService interface {
	DownloadImageAndMakePreview(
		ctx context.Context,
		downloadURL string,
		proxyHeader http.Header,
		width int,
		height int,
	) (image.Image, error)
}
