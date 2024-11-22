package downloader

import (
	"errors"
	"fmt"
)

var (
	// ErrDownloadedFileNotImage ошибка скачанный файл не изображение.
	ErrDownloadedFileNotImage = errors.New("downloaded file is not an image")
	// ErrServerDoesNotExist ошибка сервер не существует.
	ErrServerDoesNotExist = errors.New("server does not exist")
	// ErrImageDoesNotFoundOnServer ошибка изображение не нашлось на сервере.
	ErrImageDoesNotFoundOnServer = errors.New("image does not exist")
)

// FailToDownloadImageErr ошибка при загрузке изображения с сервера.
type FailToDownloadImageErr struct {
	statusCode int
}

func (e *FailToDownloadImageErr) Error() string {
	return fmt.Sprintf("fail to download image, err status code: %d", e.statusCode)
}
