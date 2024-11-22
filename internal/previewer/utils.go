package previewer

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
)

// readImgFromFile загрузить изображение из файла.
func readImgFromFile(filename string) (image.Image, error) {
	f, err := os.Open(filename) //nolint:gosec
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	return img, err
}

// writeImgToFile записать изображение в файл.
func writeImgToFile(img image.Image, filename string) error {
	f, err := os.Create(filename) //nolint:gosec
	if err != nil {
		return err
	}
	defer f.Close()
	if err = jpeg.Encode(f, img, nil); err != nil {
		return err
	}
	return nil
}

// filesHaveSameContent у файлов одинаковое содержимое в байтах.
func filesHaveSameContent(filename1, filename2 string) (bool, error) {
	f1, err := os.ReadFile(filename1) //nolint:gosec
	if err != nil {
		return false, err
	}
	f2, err := os.ReadFile(filename2) //nolint:gosec
	if err != nil {
		return false, err
	}
	eq := bytes.Equal(f1, f2)
	return eq, nil
}
