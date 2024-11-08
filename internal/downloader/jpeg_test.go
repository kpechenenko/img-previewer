package downloader

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// JpegImageDownloader_TestProxyHeaders проверка того,
// что загрузчик проксирует заголовки при отправке запроса на загрузку изображения.
// Для этого в одной горутине поднимается вебсервер, в другой - отправляется запрос из загрузчика.
func TestJpegImageDownloaderProxyHeaders(t *testing.T) {
	d := NewJPEGImageDownloader()
	header := http.Header{
		"Header1": {"1"},
		"Header2": {"2"},
		"Header3": {"3"},
	}
	address := "127.0.0.1:9000"

	var actualHeader http.Header
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(_ http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		actualHeader = r.Header
	})
	srv := http.Server{Addr: address, Handler: mux, ReadHeaderTimeout: 5 * time.Second}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			assert.FailNow(t, err.Error())
		}
	}()
	_, _ = d.Download(context.Background(), "http://"+address, header)
	err := srv.Shutdown(context.Background())
	assert.NoError(t, err)
	for h := range header {
		assert.Contains(t, actualHeader, h)
	}
}
