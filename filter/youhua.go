package filter

import (
	"fmt"
	"image"
	"image/color"
	"image_pdf_shengmo/utils"
)

type point struct {
	x, y int
	gray uint8
	c    color.Color
}

func getKey(x, y int) string {
	return fmt.Sprintf("%d_%d", x, y)
}

type youhua struct {
	name string
}

// Filter filter youhua
func (instance *youhua) Filter(img *utils.Img) *image.RGBA {
	bounds := img.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()

	cache := make(map[string]point)
	for i := 1; i < dx-1; i++ {
		for j := 1; j < dy-1; j++ {
			c := img.GrayAt(i, j)

			var cha uint8 = 2
			// if c > 200 {
			// 	cha = 8
			// } else if c > 100 {
			// 	cha = 7
			// } else if c > 50 {
			// 	cha = 6
			// }
			hasSame := false
			// SearchSame:
			for x := i - 1; x < i+1; x++ {
				for y := j - 1; y < j+1; y++ {
					if x == i && y == j {
						continue
					}

					if _, ok := cache[getKey(x, y)]; ok {
						continue
					}
					cCheck := img.GrayAt(x, y)

					if utils.AbsUint8(c, cCheck) > cha {
						cache[getKey(x, y)] = point{
							x:    x,
							y:    y,
							gray: cCheck,
							c:    img.At(x, y),
						}
						hasSame = true
						// break SearchSame
					}
				}
			}

			if hasSame {
				cache[getKey(i, j)] = point{
					x:    i,
					y:    j,
					gray: c,
					c:    img.At(i, j),
				}
			}
		}
	}

	newRgba := image.NewRGBA(bounds)
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			newRgba.SetRGBA(i, j, color.RGBA{255, 255, 255, 255})
		}
	}
	for _, p := range cache {
		newRgba.SetRGBA(p.x, p.y, utils.GetGrayColor(p.gray))
		// newRgba.Set(p.x, p.y, p.c)
	}
	return newRgba
}
func (instance *youhua) Name() string {
	return instance.name
}
func (instance *youhua) SetOption(options interface{}) {

}
func (instance *youhua) GetDescription() string {
	return fmt.Sprintf("filter name = %s", instance.name)
}

// Youhua filter youhua
var Youhua IFilter = &youhua{"youhua"}
