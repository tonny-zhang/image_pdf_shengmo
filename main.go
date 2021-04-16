package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/signintech/gopdf"
)

type img struct {
	src    string
	width  int
	height int
}
type point struct {
	x       int
	y       int
	gray    uint8
	grayNew uint8
}

var step = 1
var scale float32 = 1.1
var distance float64 = 6

//图片灰化处理
func hdImage(m image.Image) *image.RGBA {
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	var valTotal uint64 = 0
	var cache = make(map[uint8][]point)

	// 生成灰度图
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			_, g, _, _ := colorRgb.RGBA()
			gUint8 := uint8(g >> 8)
			cache[gUint8] = append(cache[gUint8], point{
				x:    i,
				y:    j,
				gray: gUint8,
			})
			valTotal += uint64(gUint8)
			newRgba.SetRGBA(i, j, color.RGBA{gUint8, gUint8, gUint8, 255})
		}
	}

	// 计算分割点
	valPerc := uint8(float32(valTotal/uint64(dx*dy)) / scale)
	// if valPerc < 180 {
	// 	valPerc += 30
	// }
	// valPerc += 50
	// fmt.Println("valPerc = ", valPerc)

	var listToSet = make([]point, 0)
	var i uint8 = 255
	var queue = make([]point, 0)
	var checkedList = make(map[string]bool)
	for ; i > valPerc; i-- {
		list := cache[i]
		// list := cache[140]
		if len(list) > 0 {
			queue = append(queue, list...)
		} else {
			// continue
		}
		fmt.Println(i, len(queue))
		// for len(queue) > 0 {
		// 	p := queue[0]
		// 	queue = queue[1:]

		// 	keyCache := fmt.Sprintf("%d_%d", p.x, p.y)
		// 	if _, ok := checkedList[keyCache]; ok {
		// 		continue
		// 	}
		// 	hasSame := false
		// 	for x, j1 := p.x-step, p.x+step; x <= j1; x++ {
		// 		if x < 0 || x >= dx {
		// 			continue
		// 		}
		// 		for y, j2 := p.y-step, p.y+step; y <= j2; y++ {
		// 			if y < 0 || y >= dy {
		// 				continue
		// 			}
		// 			if x == p.x && y == p.y {
		// 				continue
		// 			}
		// 			key := fmt.Sprintf("%d_%d", x, y)
		// 			if _, ok := checkedList[key]; ok {
		// 				continue
		// 			}

		// 			colorRgb := m.At(x, y)
		// 			_, g, _, _ := colorRgb.RGBA()
		// 			gUint8 := uint8(g >> 8)
		// 			if math.Abs(float64(gUint8)-float64(p.gray)) <= 5 {
		// 				hasSame = true
		// 				queue = append([]point{
		// 					{
		// 						x:    x,
		// 						y:    y,
		// 						gray: gUint8,
		// 					},
		// 				}, queue...)

		// 				fmt.Println("add", x, y)
		// 			}
		// 		}
		// 	}
		// 	checkedList[keyCache] = true
		// 	if hasSame {
		// 		p.grayNew = 255
		// 		listToSet = append(listToSet, p)
		// 	}
		// }

		for len(queue) > 0 {
			// fmt.Println("------", i, len(queue))
			var queueTmp = make([]point, 0)
			for _, p := range queue {
				keyCache := fmt.Sprintf("%d_%d", p.x, p.y)
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
						key := fmt.Sprintf("%d_%d", x, y)
						if _, ok := checkedList[key]; ok {
							continue
						}
						colorRgb := m.At(x, y)
						_, g, _, _ := colorRgb.RGBA()
						gUint8 := uint8(g >> 8)
						if math.Abs(float64(gUint8)-float64(p.gray)) <= distance {
							hasSame = true

							queueTmp = append(queueTmp, point{
								x:    x,
								y:    y,
								gray: gUint8,
							})
						}
					}
				}
				checkedList[keyCache] = true
				if hasSame {
					p.grayNew = 255
					listToSet = append(listToSet, p)
				} else {
					// newRgba.SetRGBA(p.x, p.y, color.RGBA{0, 0, 255, 255})
				}
				// fmt.Println("set", p.x, p.y, p.gray, hasSame, len(queueTmp))
			}
			queue = queueTmp
		}
		// break
	}

	fmt.Println("len(listToSet) =", len(listToSet))
	for _, p := range listToSet {
		if p.grayNew > 0 {
			// fmt.Println(p.x, p.y, p.grayNew)
			newRgba.SetRGBA(p.x, p.y, color.RGBA{p.grayNew, p.grayNew, p.grayNew, 255})
		}
	}

	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			_, g, _, _ := colorRgb.RGBA()
			gUint8 := uint8(g >> 8)
			// gUint8 := uint8(float32(r>>8) * 1.5) // uint8(r>>8) + 30
			if gUint8 > valPerc {
				newRgba.SetRGBA(i, j, color.RGBA{255, 255, 255, 255})
			} else {
				// newRgba.SetRGBA(i, j, color.RGBA{0, 0, 0, 255})
			}
		}
	}
	return newRgba
}

//图片编码
func encode(inputName string, rgba *image.RGBA) {
	file, _ := os.Create(inputName)
	defer file.Close()

	if strings.HasSuffix(inputName, "jpg") || strings.HasSuffix(inputName, "jpeg") {
		jpeg.Encode(file, rgba, &jpeg.Options{
			Quality: 100,
		})
	} else if strings.HasSuffix(inputName, "png") {
		png.Encode(file, rgba)
	} else if strings.HasSuffix(inputName, "gif") {
		gif.Encode(file, rgba, nil)
	} else {
		fmt.Println("不支持的图片格式")
	}
}

//图片灰化处理
func hdImageNormal(img string) (imgTarget string, err error) {
	f, e := os.Open(img)
	if e != nil {
		err = e
		return
	}
	defer f.Close()
	dir := path.Join(filepath.Dir(img), "gray")
	os.MkdirAll(dir, 0755)
	m, _, _ := image.Decode(f)
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)

	// 生成灰度图
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, _ := colorRgb.RGBA()
			gUint8 := uint8(((r + g + b) / 3) >> 8)

			newRgba.SetRGBA(i, j, color.RGBA{gUint8, gUint8, gUint8, 255})
		}
	}
	imgTarget = path.Join(dir, path.Base(img))
	fmt.Println("生成灰度图：" + imgTarget)
	encode(imgTarget, newRgba)
	return
}

//图片灰化处理
func hdImageNormal2(img string) (imgTarget string, err error) {
	f, e := os.Open(img)
	if e != nil {
		err = e
		return
	}
	defer f.Close()
	dir := path.Join(filepath.Dir(img), "gray2")
	os.MkdirAll(dir, 0755)
	m, _, _ := image.Decode(f)
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)

	// 生成灰度图
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, _ := colorRgb.RGBA()
			gUint8 := uint8(((r + g + b) / 3) >> 8)

			if gUint8 > 150 {
				gUint8 = 255
			} else {
				// gUint8 = 80
			}
			newRgba.SetRGBA(i, j, color.RGBA{gUint8, gUint8, gUint8, 255})
		}
	}
	imgTarget = path.Join(dir, path.Base(img))
	fmt.Println("生成灰度图2：" + imgTarget)
	encode(imgTarget, newRgba)
	return
}

//图片灰化处理
func hdImageNormal3(img string) (info img, err error) {
	f, e := os.Open(img)
	if e != nil {
		err = e
		return
	}
	defer f.Close()
	dir := path.Join(filepath.Dir(img), "gray3")
	os.MkdirAll(dir, 0755)
	m, _, _ := image.Decode(f)
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)

	var valTotal uint64 = 0
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, _ := colorRgb.RGBA()
			gUint8 := uint8(((r + g + b) / 3) >> 8)

			valTotal += uint64(gUint8)
			newRgba.SetRGBA(i, j, color.RGBA{gUint8, gUint8, gUint8, 255})
		}
	}
	var valPerc = uint8(float32(valTotal/uint64(dx*dy)) / scale)
	// fmt.Println("valPerc =", valPerc)
	// 生成灰度图
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, _ := colorRgb.RGBA()
			gUint8 := uint8(((r + g + b) / 3) >> 8)

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
		c1 := m.At(x, 0)
		r1, g1, b1, _ := c1.RGBA()
		gray1 := uint8(((r1 + g1 + b1) / 3) >> 8)

		c2 := m.At(x, dy-1)
		r2, g2, b2, _ := c2.RGBA()
		gray2 := uint8(((r2 + g2 + b2) / 3) >> 8)
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
		c1 := m.At(0, y)
		r1, g1, b1, _ := c1.RGBA()
		gray1 := uint8(((r1 + g1 + b1) / 3) >> 8)

		c2 := m.At(dx-1, y)
		r2, g2, b2, _ := c2.RGBA()
		gray2 := uint8(((r2 + g2 + b2) / 3) >> 8)
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
		keyCache := fmt.Sprintf("%d_%d", p.x, p.y)
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
				key := fmt.Sprintf("%d_%d", x, y)
				if _, ok := checkedList[key]; ok {
					continue
				}

				colorRgb := m.At(x, y)
				r, g, b, _ := colorRgb.RGBA()
				gUint8 := uint8(((r + g + b) / 3) >> 8)
				if math.Abs(float64(gUint8)-float64(p.gray)) <= distance {
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
			p.grayNew = 255
		}
		checkedList[keyCache] = p
	}
	for _, p := range checkedList {
		if p.grayNew > 0 {
			newRgba.SetRGBA(p.x, p.y, color.RGBA{255, 255, 255, 255})
		}
	}

	imgTarget := path.Join(dir, path.Base(img)+".png")
	fmt.Println("生成灰度图3：" + imgTarget)
	encode(imgTarget, newRgba)
	info.src = imgTarget
	info.width = dx
	info.height = dy
	return
}
func hdImageNormal4(img string) (imgTarget string, err error) {
	f, e := os.Open(img)
	if e != nil {
		err = e
		return
	}
	defer f.Close()
	dir := path.Join(filepath.Dir(img), "gray4")
	os.MkdirAll(dir, 0755)
	m, _, _ := image.Decode(f)
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)

	// 生成灰度图
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			colorRgb := m.At(i, j)
			r, g, b, _ := colorRgb.RGBA()
			gUint8 := uint8(((r + g + b) / 3) >> 8)

			if gUint8 > 160 {
				newRgba.SetRGBA(i, j, color.RGBA{255, 255, 255, 255})
			} else {
				newRgba.SetRGBA(i, j, color.RGBA{0, 0, 255, 255})
			}
		}
	}
	imgTarget = path.Join(dir, path.Base(img))
	fmt.Println("生成灰度图4：" + imgTarget)
	encode(imgTarget, newRgba)
	return
}
func dealImg(img string) (info img, err error) {
	// hdImageNormal(img)
	// hdImageNormal2(img)
	return hdImageNormal3(img)
	// hdImageNormal4(img)
}
func createPdf(list []img, savePath string) {
	var width float64 = 595
	var height float64 = 842
	var padding float64 = 20
	// var padding float64 = 20
	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{
		Unit: gopdf.Unit_PT,
		PageSize: gopdf.Rect{
			W: width,
			H: height,
		},
	})

	width -= padding * 2
	height -= padding * 2

	for _, img := range list {
		pdf.AddPage()
		var w, h float64 = float64(img.width), float64(img.height)
		if w > width || height < h {
			scale := w / width
			if w/width < h/height {
				scale = h / height
			}
			w /= scale
			h /= scale
		}
		pdf.Image(img.src, (width-w)/2+padding, (height-h)/2+padding, &gopdf.Rect{
			W: w,
			H: h,
		})
		// fmt.Println("add", img.src, img.width, img.height, w, h, e, (width-w)/2, (height-h)/2)
	}

	os.MkdirAll(path.Dir(savePath), 0755)
	pdf.WritePdf(savePath)
	fmt.Println("生成pdf:" + savePath)
}
func main() {
	if len(os.Args) == 1 {
		fmt.Println("[ERROR] 请输入要处理的目录!!")
		os.Exit(0)
	}
	dir := os.Args[1]
	var imgList []img
	// dir := "./img"
	list, e := os.ReadDir(dir)
	if e == nil {
		imgList = make([]img, 0)
		for _, f := range list {
			// if f.Name() != "52.jpg" {
			// 	continue
			// }
			if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
				info, e := dealImg(path.Join(dir, f.Name()))
				if e == nil {
					imgList = append(imgList, info)
				}
			}
		}
		createPdf(imgList, path.Join(dir, "pdf/print.pdf"))
	} else {
		fmt.Println("[ERROR] 读取目录[" + dir + "]错误")
	}
}
