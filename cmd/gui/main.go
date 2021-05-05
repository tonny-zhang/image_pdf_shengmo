package main

import (
	"fmt"
	"image_pdf_shengmo/deal"
	"image_pdf_shengmo/filter"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

var typeCurrent = "floder"
var filterEdge = filter.Edge
var filterGray = filter.Gray
var filterYouhua = filter.Youhua

type writerGui struct {
	writerFn func([]byte)
}

func (writer writerGui) Write(b []byte) (int, error) {
	writer.writerFn(b)
	return 0, nil
}
func main() {
	os.Setenv("FYNE_FONT", "/Library/Fonts/Arial Unicode.ttf")
	a := app.New()
	win := a.NewWindow("put image to pdf 中文")

	btnOpen := widget.NewButton("floder中文", func() {
		if typeCurrent == "floder" {
			dialog.ShowFolderOpen(func(dir fyne.ListableURI, err error) {
				fmt.Println(dir, err)
				if dir != nil {
					// p := dir.Path()
					p := "/Users/tonny/doc_zk/彬彬小学/平时作业/20210430/one"
					deal.Dir(p, true)
				}

			}, win)
		} else {
			dialog.ShowFileOpen(func(a fyne.URIReadCloser, err error) {
				fmt.Println(a, err)
				if a != nil {
					deal.File(a.URI().Path())
				}
			}, win)
		}
	})
	sOpenType := widget.NewSelect([]string{
		"floder",
		"file",
	}, func(t string) {
		btnOpen.SetText(t)
	})
	sFilter := widget.NewSelect([]string{
		filterEdge.Name(),
		filterGray.Name(),
	}, func(t string) {
		if t == filterEdge.Name() {
			deal.SetOptions(deal.Options{
				Filter: filterEdge,
			})
		} else if t == filterGray.Name() {
			deal.SetOptions(deal.Options{
				Filter: filterGray,
			})
		}
	})
	text := widget.NewTextGrid()
	text.SetText("deal process (处理过程):")

	content := container.NewVBox(
		sFilter,
		sOpenType,
		btnOpen,
	)
	content.Add(widget.NewLabel("处理过程:"))

	deal.SetOptions(deal.Options{
		Filter: filterEdge,
		Writer: writerGui{
			writerFn: func(b []byte) {
				line := fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), string(b))
				content.Add(widget.NewLabel(line))
				content.Refresh()
				// cells := make([]widget.TextGridCell, len(line))
				// for j, r := range line {
				// 	cells[j] = widget.TextGridCell{Rune: r}
				// }
				// text.SetRow(len(text.Rows), widget.TextGridRow{
				// 	Cells: cells,
				// })
			},
		},
	})

	contentScroll := container.NewScroll(content)
	win.SetContent(contentScroll)
	win.Resize(fyne.NewSize(480, 360))
	win.ShowAndRun()
}
