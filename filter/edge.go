package filter

import (
	"fmt"
	"image"
	"image_pdf_shengmo/utils"
	"math"
)

var sqrtVal = make(map[int]uint8)
var kernelX = [][]int{
	{-1, 0, 1},
	{-2, 0, 2},
	{-1, 0, 1},
}
var kernelY = [][]int{
	{-1, -2, -1},
	{0, 0, 0},
	{1, 2, 1},
}

func init() {
	for i := 0; i < 255*255; i++ {
		sqrtVal[i] = 255 - uint8(math.Sqrt(float64(i)))
	}
}

// Edge filter
type edge struct {
	name    string
	grayCut uint8
}

// https://github.com/miguelmota/sobel/blob/master/sobel.js
// Filter fitler for Edge
func (edge *edge) Filter(img *utils.Img) *image.RGBA {
	bounds := img.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	for x := 1; x < dx-1; x++ {
		for y := 1; y < dy-1; y++ {
			c1 := int(img.GrayAt(x-1, y-1))
			c2 := int(img.GrayAt(x, y-1))
			c3 := int(img.GrayAt(x+1, y-1))
			c4 := int(img.GrayAt(x-1, y))
			c5 := int(img.GrayAt(x, y))
			c6 := int(img.GrayAt(x+1, y))
			c7 := int(img.GrayAt(x-1, y+1))
			c8 := int(img.GrayAt(x, y+1))
			c9 := int(img.GrayAt(x+1, y+1))

			pixelX := kernelX[0][0]*c1 +
				kernelX[0][1]*c2 +
				kernelX[0][2]*c3 +
				kernelX[1][0]*c4 +
				kernelX[1][1]*c5 +
				kernelX[1][2]*c6 +
				kernelX[2][0]*c7 +
				kernelX[2][1]*c8 +
				kernelX[2][2]*c9
			pixelY := kernelY[0][0]*c1 +
				kernelY[0][1]*c2 +
				kernelY[0][2]*c3 +
				kernelY[1][0]*c4 +
				kernelY[1][1]*c5 +
				kernelY[1][2]*c6 +
				kernelY[2][0]*c7 +
				kernelY[2][1]*c8 +
				kernelY[2][2]*c9

			var magnitude = sqrtVal[pixelX*pixelX+pixelY*pixelY]
			if edge.grayCut > 0 && magnitude > edge.grayCut {
				magnitude = 255
			}
			newRgba.SetRGBA(x, y, utils.GetGrayColor(magnitude))
		}
	}
	return newRgba
}
func (edge *edge) Name() string {
	return edge.name
}
func (edge *edge) SetOption(options interface{}) {
	grayCut := uint8(options.(uint))
	edge.grayCut = grayCut
}
func (edge *edge) GetDescription() string {
	return fmt.Sprintf("filter name = %s, sobel, grayCut = %d", edge.name, edge.grayCut)
}

// Edge filter instance
var Edge IFilter = &edge{
	name:    "edge",
	grayCut: 200,
}
