package previewer

import (
	"image"
	"log/slog"
	"math"
)

// KNNPreviewer создает превью при помощи уменьшения изображение
// с использованием алгоритма ближайшего соседа для выбора сжимаемого пикселя.
type KNNPreviewer struct{}

// MakePreview создать превью изображения.
// https://medium.com/@epcm18/image-resampling-in-image-processing-f7b597ee78a8
func (s *KNNPreviewer) MakePreview(img image.Image, width, height int) (image.Image, error) {
	if width <= 0 || width > img.Bounds().Dx() {
		slog.Info(" width out of bounds img width")
		return nil, ErrInvalidImageSize
	}
	if height <= 0 || height > img.Bounds().Dy() {
		slog.Info("height out of bounds img height")
		return nil, ErrInvalidImageSize
	}
	preview := image.NewRGBA(image.Rect(0, 0, width, height))

	xRation := float64(img.Bounds().Dy()) / float64(height)
	yRation := float64(img.Bounds().Dx()) / float64(width)

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			imgX := int(math.Round(float64(x) * yRation))
			imgY := int(math.Round(float64(y) * xRation))
			preview.Set(x, y, img.At(imgX, imgY))
		}
	}

	return preview, nil
}

// NewKNNImageCompressor конструктор.
func NewKNNImageCompressor() *KNNPreviewer {
	return &KNNPreviewer{}
}
