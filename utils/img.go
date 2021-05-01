package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path"
)

// Img img struct
type Img struct {
	i image.Image
}

// Bounds get Bounds
func (img *Img) Bounds() image.Rectangle {
	return img.i.Bounds()
}

// GrayAt get gray at x, y
func (img *Img) GrayAt(x, y int) uint8 {
	c := img.i.At(x, y)
	r, g, b, _ := c.RGBA()
	return uint8(((r + g + b) / 3) >> 8)
}

// At get color
func (img *Img) At(x, y int) color.Color {
	return img.i.At(x, y)
}

// NewImg new image
func NewImg(img image.Image) *Img {
	return &Img{
		i: img,
	}
}

// Save save img
func Save(target string, img image.Image) error {
	os.MkdirAll(path.Dir(target), 0755)

	file, e := os.Create(target)
	if e != nil {
		return e
	}
	defer file.Close()

	ext := path.Ext(target)
	if ext == ".jpg" || target == ".jpeg" {
		return jpeg.Encode(file, img, &jpeg.Options{
			Quality: 100,
		})
	} else if ext == ".png" {
		return png.Encode(file, img)
	} else {
		return fmt.Errorf("不支持的图片格式")
	}
}

// LoadImage load image
func LoadImage(imgPath string) (m image.Image, isRotate bool, e error) {
	f, e := os.Open(imgPath)
	if e != nil {
		return nil, false, e
	}
	defer f.Close()
	m, _, e = image.Decode(f)

	if e != nil {
		return nil, false, e
	}
	b := m.Bounds()
	if b.Dx() > b.Dy() {
		rotate90 := image.NewRGBA(image.Rect(0, 0, m.Bounds().Dy(), m.Bounds().Dx()))
		// 矩阵旋转
		for x := m.Bounds().Min.Y; x < m.Bounds().Max.Y; x++ {
			for y := m.Bounds().Max.X - 1; y >= m.Bounds().Min.X; y-- {
				//  设置像素点
				rotate90.Set(m.Bounds().Max.Y-x, y, m.At(y, x))
			}
		}
		m = rotate90
		isRotate = true
	}
	return
}

// GetGrayColor get gray color
func GetGrayColor(gray uint8) color.RGBA {
	return color.RGBA{gray, gray, gray, 255}
}
