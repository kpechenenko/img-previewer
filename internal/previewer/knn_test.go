package previewer

import (
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/assert" //nolint:depguard
)

// TestKNNPreviewer_MakePreview проверить создание превью изображений, обработку параметров.
func TestKNNPreviewer_MakePreview(t *testing.T) {
	tests := []struct {
		name    string
		imgPath string
		width   int
		height  int
		err     error
	}{
		{
			name:    "make preview, wrong negative width",
			imgPath: "./testdata/gopher_275_183.jpeg",
			width:   -1,
			err:     ErrInvalidImageSize,
		},
		{
			name:    "make preview, wrong zero width",
			imgPath: "./testdata/gopher_275_183.jpeg",
			width:   0,
			err:     ErrInvalidImageSize,
		},
		{
			name:    "make preview, wrong negative height",
			imgPath: "./testdata/gopher_275_183.jpeg",
			height:  -1,
			err:     ErrInvalidImageSize,
		},
		{
			name:    "make preview, wrong zero height",
			imgPath: "./testdata/gopher_275_183.jpeg",
			height:  0,
			err:     ErrInvalidImageSize,
		},
		{
			name:    "make preview, new height bigger than source height",
			imgPath: "./testdata/gopher_275_183.jpeg",
			height:  200,
			err:     ErrInvalidImageSize,
		},
		{
			name:    "make preview, new width bigger than source width",
			imgPath: "./testdata/gopher_275_183.jpeg",
			height:  300,
			err:     ErrInvalidImageSize,
		},
		{
			name:    "make preview, correct width and height",
			imgPath: "./testdata/gopher_275_183.jpeg",
			width:   200,
			height:  100,
			err:     nil,
		},
		{
			name:    "make preview, correct width and height",
			imgPath: "./testdata/gopher_275_183.jpeg",
			width:   147,
			height:  170,
			err:     nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			img, err := readImgFromFile(tc.imgPath)
			assert.NoError(t, err)
			srv := NewKNNImageCompressor()
			resizedImg, err := srv.MakePreview(img, tc.width, tc.height)
			if tc.err != nil {
				assert.ErrorContains(t, err, tc.err.Error())
			} else {
				assert.Equal(t, resizedImg.Bounds().Dx(), tc.width)
				assert.Equal(t, resizedImg.Bounds().Dy(), tc.height)
			}
		})
	}
}

// TestKNNPreviewer_CheckContent проверить, что превью создаются корректно.
// Суть проверки - сравнить результат работы алгоритма с заранее созданным превью.
func TestKNNPreviewer_CheckContent(t *testing.T) {
	// Напрямую сравнить preview и img не получилось
	// т.к. в кишках при загрузке изображения из файла используется image.YCbCr, а я делал превью в image.RGBA
	// Как конвертировать image.RGBA в image.YCbCr не нашел
	// Поэтому в тесте создаю превью, временно записываю его на диск, сравниваю содержимое с уже сохраненным превью.
	// Буду рад услышать как можно избежать построения такого костыля :)
	imgPath := "./testdata/gopher_275_183.jpeg"
	img, err := readImgFromFile(imgPath)
	assert.NoError(t, err)

	srv := NewKNNImageCompressor()
	var preview image.Image
	preview, err = srv.MakePreview(img, 200, 50)
	assert.NoError(t, err)
	tmpFile := "./testdata/current_preview.jpeg"

	err = writeImgToFile(preview, tmpFile)
	assert.NoError(t, err)
	defer os.Remove(tmpFile)

	previewPath := "./testdata/gopher_200_50.jpeg"
	ok, err := filesHaveSameContent(previewPath, tmpFile)
	assert.NoError(t, err)
	assert.True(t, ok)
}
