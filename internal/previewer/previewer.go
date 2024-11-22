// Package previewer алгоритмы для создания превью изображений.
package previewer

import (
	"image"
)

// Previewer позволяет создавать превью изображений.
type Previewer interface {
	// MakePreview создать превью изображения.
	MakePreview(img image.Image, width, height int) (image.Image, error)
}
