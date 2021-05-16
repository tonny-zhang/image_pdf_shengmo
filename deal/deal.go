package deal

import (
	"fmt"
	"image_pdf_shengmo/filter"
	"image_pdf_shengmo/utils"
	"io"
	"log"
	"path"
	"strings"
	"time"
)

var currentFilter filter.IFilter = filter.Edge
var defaultWriter io.Writer = log.Default().Writer()
var currentWriter io.Writer = defaultWriter

// Options options for deal
type Options struct {
	Filter filter.IFilter
	Writer io.Writer
}

// SetOptions set options
func SetOptions(opt Options) {
	if opt.Filter != nil {
		currentFilter = opt.Filter
	}

	if opt.Writer != nil {
		currentWriter = opt.Writer
	}
}

// File deal image
func File(imgSource string) (*utils.PdfImg, error) {
	img, _, e := utils.LoadImage(imgSource)
	if e == nil {
		timeStart := time.Now()
		log.Printf("开始处理 %s", imgSource)
		fmt.Fprintln(currentWriter, fmt.Sprintf("start deal %s", imgSource))
		img := currentFilter.Filter(utils.NewImg(img))
		target := path.Join(path.Dir(imgSource), currentFilter.Name(), path.Base(imgSource)+".png")
		e := utils.Save(target, img)
		if e != nil {
			return nil, e
		}
		log.Printf("生成 %s, 用时%v", target, time.Now().Sub(timeStart))
		fmt.Fprintln(currentWriter, fmt.Sprintf("create %s, takes %v", target, time.Now().Sub(timeStart)))

		return &utils.PdfImg{
			Src:    target,
			Width:  img.Bounds().Dx(),
			Height: img.Bounds().Dy(),
		}, nil
	}
	return nil, e
}

// Dir deal dir
func Dir(dir string, createPdf bool) (err error) {
	fileList, e := utils.GetAllFile(dir)

	if e != nil {
		err = fmt.Errorf("[ERROR] 读取目录[" + dir + "]错误")
		return
	}

	imgList := make([]utils.PdfImg, 0)
	for _, f := range fileList {
		if !f.IsDir() && !strings.HasPrefix(f.Name(), ".") {
			img, e := File(path.Join(dir, f.Name()))
			if e == nil {
				imgList = append(imgList, *img)
			}
		}
	}
	if createPdf {
		pdfPath := path.Join(dir, "pdf", currentFilter.Name()+".pdf")
		utils.CreatePdf(imgList, pdfPath)
		// log.Println("生成pdf:", pdfPath)
		fmt.Fprintln(currentWriter, "pdf:", pdfPath)
		utils.Open(pdfPath)
	}

	return
}
