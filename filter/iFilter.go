package filter

import (
	"image"
	"image_pdf_shengmo/utils"
)

// IFilter interface filter
type IFilter interface {
	Filter(img *utils.Img) *image.RGBA
	Name() string
	GetDescription() string
	SetOption(options interface{})
}
