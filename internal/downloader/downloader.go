// Package downloader downloader содержит код загрузчика изображений по http.
package downloader

import (
	"context"
	"image"
	"net/http"
)

// HTTPImageDownloader загрузчик изображений по сети.
type HTTPImageDownloader interface {
	// Download загрузить изображение по сети. Проксирует заголовки.
	// Выдает ошибку, если удаленный сервер не найден или если в ответе вернулось не изображение.
	Download(ctx context.Context, url string, header http.Header) (image.Image, error)
}
