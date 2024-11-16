package downloader

import (
	"errors"
	"fmt"
)

var (
	ErrDownloadedFileNotImage    = errors.New("downloaded file is not an image")
	ErrServerDoesNotExist        = errors.New("server does not exist")
	ErrImageDoesNotFoundOnServer = errors.New("image does not exist")
)

type FailToDownloadImageErr struct {
	statusCode int
}

func (e *FailToDownloadImageErr) Error() string {
	return fmt.Sprintf("fail to download image, err status code: %d", e.statusCode)
}
