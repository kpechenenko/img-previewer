package downloader

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	//nolint:depguard
	"github.com/stretchr/testify/assert"
)

// TestDownloaderProxyHeaders проверка того, что загрузчик проксирует заголовки
// при отправке запроса на загрузку изображения.
func TestDownloaderProxyHeaders(t *testing.T) {
	header := http.Header{
		"Header1": []string{"1"},
		"Header2": []string{"2"},
		"Header3": []string{"3"},
	}
	var actualHeader http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		actualHeader = r.Header
	}))
	d := NewJPEGImageDownloader(srv.Client())
	_, _ = d.Download(context.Background(), srv.URL+"/test", header)
	for h := range header {
		assert.Contains(t, actualHeader, h)
	}
}

// TestDownloaderReturnErrWhenServerExistAndImageDoesNotExist проверяет тип возвращаемой ошибки при обращении к
// несуществующему изображению на существующем сервере.
func TestDownloaderReturnErrWhenServerExistAndImageDoesNotExist(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	d := NewJPEGImageDownloader(srv.Client())
	_, err := d.Download(context.Background(), srv.URL, nil)
	assert.ErrorIs(t, err, ErrImageDoesNotFoundOnServer)
}

// TestDownloaderProxyOriginalErr проверяет, что ошибка от оригинального сервера проксиурется загрузчиком.
func TestDownloaderProxyOriginalErr(t *testing.T) {
	statusCode := http.StatusBadRequest
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(statusCode)
	}))
	d := NewJPEGImageDownloader(srv.Client())
	_, err := d.Download(context.Background(), srv.URL, nil)
	var downloadImgErr *FailToDownloadImageErr
	if errors.As(err, &downloadImgErr) {
		assert.Equal(t, statusCode, downloadImgErr.statusCode)
	} else {
		assert.FailNow(t, "error type should be *FailToDownloadImageErr")
	}
}

// TestDownloaderDownloadImageCorrect проверить, что загрузчик корректно скачивает изображения.
func TestDownloaderDownloadImageCorrect(t *testing.T) {
	srv := httptest.NewServer(
		http.StripPrefix(
			"/testdata",
			http.FileServer(http.Dir("./testdata")),
		),
	)
	d := NewJPEGImageDownloader(srv.Client())
	img, err := d.Download(context.Background(), srv.URL+"/testdata/gopher_200_50.jpeg", nil)
	assert.NoError(t, err)
	assert.NotNil(t, img)
}
