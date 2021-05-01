package utils

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var maxFileNameLen = 30

// GetAllFile get all files
func GetAllFile(dir string) ([]fs.DirEntry, error) {
	info, e := os.Stat(dir)
	if e != nil {
		return nil, e
	}

	if info.IsDir() {
		list, e := os.ReadDir(dir)
		if e == nil {
			fileList := make([]fs.DirEntry, 0)
			for _, f := range list {
				if !f.IsDir() {
					fileList = append(fileList, f)
				}
			}

			sort.Slice(fileList, func(i, j int) bool {
				a := fileList[i].Name()
				b := fileList[j].Name()

				return strings.Compare(strings.Repeat("0", maxFileNameLen-len(a))+a, strings.Repeat("0", maxFileNameLen-len(b))+b) < 0
			})

			return fileList, nil
		}
		return nil, e
	}
	return nil, fmt.Errorf("[%s] not dir", dir)
}

// GetAllImage get all image
func GetAllImage(dir string) (imgList []PdfImg, err error) {
	fileList, err := GetAllFile(dir)
	if err == nil {
		for _, f := range fileList {
			imgPath := filepath.Join(dir, f.Name())
			img, isRotate, e := LoadImage(imgPath)
			if e == nil {
				if isRotate {
					imgPath = filepath.Join(dir, "pdf", f.Name())
					Save(imgPath, img)
				}

				_pdfImg := PdfImg{
					Src:    imgPath,
					Width:  img.Bounds().Dx(),
					Height: img.Bounds().Dy(),
				}
				imgList = append(imgList, _pdfImg)

				msg := imgPath
				if isRotate {
					msg += " 自动旋转"
				}
				log.Printf("读取图片 %s", msg)
			}
		}
	}
	return
}
