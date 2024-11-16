package downloader

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"io"
	"log/slog"
	"net/http"
	"syscall"
)

// JPEGImageDownloader загрузчик изображений jpeg.
type JPEGImageDownloader struct {
	client *http.Client
}

// Download загрузить изображение по сети. Проксирует заголовки.
// Выдает ошибку, если удаленный сервер не найден или если в ответе вернулось не изображение, либо изображение не jpeg.
func (s *JPEGImageDownloader) Download(ctx context.Context, url string, header http.Header) (image.Image, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		slog.Error("fail to create http request", "error", err)
		return nil, err
	}
	for k, v := range header {
		req.Header[k] = v
	}
	resp, err := s.client.Do(req)
	if err != nil {
		slog.Error("fail to do http request to download file", "error", err)
		// удаленный сервер не существует
		if errors.Is(err, syscall.ECONNREFUSED) {
			return nil, ErrServerDoesNotExist
		}
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		slog.Error("resp status code is not OK", "statusCode", resp.StatusCode)
		// удаленный сервер существует, картинка не нашлась
		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrImageDoesNotFoundOnServer
		}
		// удаленный сервис вернул ошибку, проксируем код
		return nil, &FailToDownloadImageErr{statusCode: resp.StatusCode}
	}
	var b []byte
	if b, err = io.ReadAll(resp.Body); err != nil {
		slog.Error("fail to read image from response body", "error", err)
		return nil, err
	}
	contentType := http.DetectContentType(b)
	// загруженный файл не является изображением
	if contentType != "image/jpeg" {
		slog.Error("content type of downloaded image is not image/jpeg", "contentType", contentType)
		return nil, ErrDownloadedFileNotImage
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
func NewJPEGImageDownloader(client *http.Client) *JPEGImageDownloader {
	return &JPEGImageDownloader{client: client}
}
