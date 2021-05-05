package main

import (
	"fmt"
	"image_pdf_shengmo/deal"
	"image_pdf_shengmo/filter"
	"image_pdf_shengmo/utils"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
)

var version = ""
var buildtime = ""

var filterEdge = filter.Edge
var filterGray = filter.Gray
var filterYouhua = filter.Youhua

var currentFilter filter.IFilter
var filters = map[string]filter.IFilter{}

func init() {
	filters[filterEdge.Name()] = filterEdge
	filters[filterGray.Name()] = filterGray
	filters[filterYouhua.Name()] = filterYouhua
}

// func deal(imgSource string) (*utils.PdfImg, error) {
// 	img, _, e := utils.LoadImage(imgSource)
// 	if e == nil {
// 		timeStart := time.Now()
// 		log.Printf("开始处理 %s", imgSource)
// 		img := currentFilter.Filter(utils.NewImg(img))
// 		target := path.Join(path.Dir(imgSource), currentFilter.Name(), path.Base(imgSource)+".png")
// 		e := utils.Save(target, img)
// 		if e != nil {
// 			return nil, e
// 		}
// 		log.Printf("生成 %s, 用时%v", target, time.Now().Sub(timeStart))

// 		return &utils.PdfImg{
// 			Src:    target,
// 			Width:  img.Bounds().Dx(),
// 			Height: img.Bounds().Dy(),
// 		}, nil
// 	}
// 	return nil, e
// }

func main() {
	app := cli.NewApp()
	app.Name = "imgPrint"
	app.Description = "deal image and create pdf"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "filter",
			Value: filterEdge.Name(),
			Usage: "filter name for deal",
		},
		cli.BoolFlag{
			Name:  "nopdf",
			Usage: "not creat pdf (default to create)",
		},
		cli.Float64Flag{
			Name:  "gray-scale",
			Usage: "scale for average of gray",
		},
		cli.UintFlag{
			Name:  "gray-distance",
			Usage: "distance for check neighbor",
		},
		cli.UintFlag{
			Name:  "edge-graycut",
			Usage: "grayCut for cut light gray",
			Value: 200,
		},
	}

	app.Action = func(c *cli.Context) error {
		filterName := c.String("filter")
		var ok = false
		if currentFilter, ok = filters[filterName]; !ok {
			fmt.Println(fmt.Errorf("not support filter [%s]", filterName))
			os.Exit(0)
		}
		if filterName == filterGray.Name() {
			scale := c.Float64("gray-scale")
			distance := c.Uint("gray-distance")

			filterGray.SetOption(map[string]interface{}{
				"scale":    scale,
				"distance": distance,
			})
		} else if filterName == filterEdge.Name() {
			grayCut := c.Uint("edge-graycut")
			filterEdge.SetOption(grayCut)
		}

		args := c.Args()
		dir := args.First()
		if "" == dir {
			fmt.Println("please enter dir or image path")
			os.Exit(0)
		}

		info, e := os.Stat(dir)
		if e != nil {
			fmt.Println(e)
			os.Exit(0)
		}

		fmt.Println(currentFilter.GetDescription())
		deal.SetOptions(deal.Options{
			Filter: currentFilter,
		})
		if info.IsDir() {
			err := deal.Dir(dir, !c.Bool("nopdf"))
			if err != nil {
				fmt.Println(err)
				os.Exit(0)
			}
		} else {
			_, e := deal.File(dir)
			if e != nil {
				fmt.Println(e)
			}
		}
		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list filter",
			Action: func(c *cli.Context) error {
				fmt.Println("filter name list:")
				index := 1
				for key := range filters {
					fmt.Printf("\t %d. %s\n", index, key)
					index++
				}
				return nil
			},
		},
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "show version",
			Action: func(c *cli.Context) error {
				fmt.Printf("Desc: %s", "处理图片并生成pdf方便打印\n")
				fmt.Printf("  version: %s\n  build time: %s\n", version, buildtime)
				return nil
			},
		},
		{
			Name:  "pdf",
			Usage: "put image to pdf\n\t\t\tUsage: " + filepath.Base(os.Args[0]) + " pdf dir",
			Action: func(c *cli.Context) error {
				dir := c.Args().First()
				nameSave := c.Args().Get(1)
				if nameSave == "" {
					nameSave = time.Now().Format("2006-01-02-150405") + ".pdf"
				}
				imgList, e := utils.GetAllImage(dir)
				if e == nil {
					pdfPath := filepath.Join(dir, "pdf", nameSave)
					utils.CreatePdf(imgList, pdfPath)
					log.Println("生成pdf:", pdfPath)
					utils.Open(pdfPath)
				} else {
					fmt.Println(e)
				}
				return nil
			},
		},
	}

	app.Run(os.Args)
}
