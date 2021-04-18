package utils

import (
	"os"
	"path"

	"github.com/signintech/gopdf"
)

// PdfImg img for pdf
type PdfImg struct {
	Src           string
	Width, Height int
}

var width float64 = 595
var height float64 = 842
var padding float64 = 20

// CreatePdf create pdf
func CreatePdf(list []PdfImg, savePath string) error {
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
		var w, h float64 = float64(img.Width), float64(img.Height)
		if w > width || height < h {
			scale := w / width
			if w/width < h/height {
				scale = h / height
			}
			w /= scale
			h /= scale
		}
		e := pdf.Image(img.Src, (width-w)/2+padding, (height-h)/2+padding, &gopdf.Rect{
			W: w,
			H: h,
		})
		if e != nil {
			return e
		}
	}

	os.MkdirAll(path.Dir(savePath), 0755)
	return pdf.WritePdf(savePath)
}
