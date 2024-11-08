package downloader

import (
	"bytes"
	"context"
	"image"
	"image/jpeg"
	"io"
	"log/slog"
	"net/http"
)

// JPEGImageDownloader загрузчик изображений jpeg.
type JPEGImageDownloader struct{}

// Download загрузить изображение по сети. Проксирует заголовки.
// Выдает ошибку, если удаленный сервер не найден или если в ответе вернулось не изображение, либо изображение не jpeg.
func (s *JPEGImageDownloader) Download(ctx context.Context, url string, header http.Header) (image.Image, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("fail to create http request", "error", err)
		return nil, err
	}
	req.Header = header
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		slog.Error("fail to do http request to download file", "error", err)
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		slog.Error("resp status code is not OK", "statusCode", resp.StatusCode)
		return nil, err
	}
	var b []byte
	if b, err = io.ReadAll(resp.Body); err != nil {
		slog.Error("fail to read image from response body", "error", err)
		return nil, err
	}
	contentType := http.DetectContentType(b)
	if contentType != "image/jpeg" {
		slog.Error("content type of downloaded image is not image/jpeg", "contentType", contentType)
		return nil, err
	}
	var img image.Image
	img, err = jpeg.Decode(bytes.NewReader(b))
	if err != nil {
		slog.Error("fail to decode jpeg image", "error", err)
		return nil, err
	}
	return img, nil
}

// NewJPEGImageDownloader конструктор.
func NewJPEGImageDownloader() *JPEGImageDownloader {
	return &JPEGImageDownloader{}
}
