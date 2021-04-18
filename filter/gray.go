package filter

import (
	"fmt"
	"image"
	"image/color"
	"image_pdf_shengmo/utils"
)

var step = 1 // 查看相邻像素灰度的步长（步长越大计算量越大）
// var scale float32 = 1.1  // 得到的灰度平均值的缩放系统
// var distance float64 = 7 // 计算相邻像素灰度值的容差
var scaleDefault float32 = 1.1 // 得到的灰度平均值的缩放系统
var distanceDefault uint8 = 7  // 计算相邻像素灰度值的容差

type gray struct {
	name     string
	scale    float32
	distance uint8
}

func (instance *gray) Filter(img *utils.Img) *image.RGBA {
	bounds := img.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)

	var valTotal uint64 = 0
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			gUint8 := img.GrayAt(i, j)

			valTotal += uint64(gUint8)
			newRgba.SetRGBA(i, j, color.RGBA{gUint8, gUint8, gUint8, 255})
		}
	}
	var valPerc = uint8(float32(valTotal/uint64(dx*dy)) / instance.scale)
	// fmt.Println("valPerc =", valPerc)
	// 生成灰度图
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			gUint8 := img.GrayAt(i, j)

			if gUint8 > valPerc {
				gUint8 = 255
			} else {
				// gUint8 = 80
			}
			newRgba.SetRGBA(i, j, color.RGBA{gUint8, gUint8, gUint8, 255})
		}
	}

	// 对边缘造成的阴影进行处理
	var queue = make([]point, 0)
	for x := 0; x < dx; x++ {
		gray1 := img.GrayAt(x, 0)
		gray2 := img.GrayAt(x, dy-1)

		queue = append(queue, point{
			x:    x,
			y:    0,
			gray: gray1,
		}, point{
			x:    x,
			y:    dy - 1,
			gray: gray2,
		})
	}
	for y := 0; y < dy; y++ {
		gray1 := img.GrayAt(0, y)
		gray2 := img.GrayAt(dx-1, y)

		queue = append(queue, point{
			x:    0,
			y:    y,
			gray: gray1,
		}, point{
			x:    dx - 1,
			y:    y,
			gray: gray2,
		})
	}

	checkedList := make(map[string]point)

	for len(queue) > 0 {
		p := queue[0]
		queue = queue[1:]
		keyCache := getKey(p.x, p.y)
		if _, ok := checkedList[keyCache]; ok {
			continue
		}
		hasSame := false
		for x, j1 := p.x-step, p.x+step; x <= j1; x++ {
			if x < 0 || x >= dx {
				continue
			}
			for y, j2 := p.y-step, p.y+step; y <= j2; y++ {
				if y < 0 || y >= dy {
					continue
				}
				if x == p.x && y == p.y {
					continue
				}
				key := getKey(x, y)
				if _, ok := checkedList[key]; ok {
					continue
				}

				gUint8 := img.GrayAt(x, y)

				if utils.AbsUint8(gUint8, p.gray) <= instance.distance {
					hasSame = true
					queue = append(queue, point{
						x:    x,
						y:    y,
						gray: gUint8,
					})
				}
			}
		}
		if hasSame {
			p.gray = 255
		}
		checkedList[keyCache] = p
	}
	for _, p := range checkedList {
		if p.gray > 0 {
			newRgba.SetRGBA(p.x, p.y, color.RGBA{255, 255, 255, 255})
		}
	}
	return newRgba
}

// Name get name
func (instance *gray) Name() string {
	return instance.name
}

func (instance *gray) SetOption(options interface{}) {
	opt := options.(map[string]interface{})
	if scale, ok := opt["scale"]; ok {
		instance.SetScale(float32(scale.(float64)))
	}
	if distance, ok := opt["distance"]; ok {
		instance.SetDistance(uint8(distance.(uint)))
	}
}
func (instance *gray) GetDescription() string {
	return fmt.Sprintf("filter name = %s, scale = %f, distance = %d", instance.name, instance.scale, instance.distance)
}

// SetScale set scale
func (instance *gray) SetScale(scale float32) {
	if scale > 0 {
		instance.scale = scale
	}
}

// SetDistance set distance
func (instance *gray) SetDistance(distance uint8) {
	if distance > 0 {
		instance.distance = distance
	}
}

// Gray filter gray
var Gray IFilter = &gray{
	name:     "gray",
	scale:    scaleDefault,
	distance: distanceDefault,
}
