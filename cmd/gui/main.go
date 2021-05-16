package main

import (
	"fmt"
	"image_pdf_shengmo/deal"
	"image_pdf_shengmo/filter"
	"image_pdf_shengmo/utils"
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

var logList = make([]fyne.CanvasObject, 0)
var logContainer *fyne.Container

func clearLog() {
	for _, log := range logList {
		logContainer.Remove(log)
	}
	logList = make([]fyne.CanvasObject, 0)
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
					p := dir.Path()
					// p := "/Users/tonny/doc_zk/彬彬小学/平时作业/20210430/one"
					clearLog()
					deal.Dir(p, true)
				}

			}, win)
		} else {
			dialog.ShowFileOpen(func(a fyne.URIReadCloser, err error) {
				if a != nil {
					clearLog()
					img, e := deal.File(a.URI().Path())
					if e == nil {
						utils.Open(img.Src)
					}
				}
			}, win)
		}
	})
	sOpenType := widget.NewSelect([]string{
		"floder",
		"file",
	}, func(t string) {
		btnOpen.SetText(t)
		typeCurrent = t
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

	logContainer = container.NewVBox(
		sFilter,
		sOpenType,
		btnOpen,
	)
	logContainer.Add(widget.NewLabel("处理过程:"))
	logContainerScroller := container.NewScroll(logContainer)
	win.SetContent(logContainerScroller)

	deal.SetOptions(deal.Options{
		Filter: filterEdge,
		Writer: writerGui{
			writerFn: func(b []byte) {
				line := fmt.Sprintf("%s %s", time.Now().Format("2006-01-02 15:04:05"), string(b))
				txt := widget.NewLabel(line)
				logList = append(logList, txt)
				logContainer.Add(txt)
				logContainer.Refresh()

				logContainerScroller.ScrollToBottom()
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

	win.Resize(fyne.NewSize(480, 360))
	win.ShowAndRun()
}
