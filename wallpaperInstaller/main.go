package main

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/thatnerd/betterapng"
)

type LoadingWidget struct {
	widget.BaseWidget
	FrameCounter *widget.Label
	frame        uint64
	Image        *canvas.Raster
	BAPNG 	  *betterapng.BAPNG
	images [][]byte
	currentImage image.Image
}

func (item *LoadingWidget) Update() {
	item.FrameCounter.SetText(fmt.Sprintf("%d", item.frame))
	item.frame++
	if item.frame >= uint64(item.BAPNG.GetNumberOfFrames()) {
		item.frame = 0
	}
	item.currentImage, _ = png.Decode(bytes.NewReader(item.images[item.frame]))
}

func NewMyListItemWidget() *LoadingWidget {

	item := &LoadingWidget{
		FrameCounter: widget.NewLabel("0"),
	}
	f, err := os.Open("installing.bapng")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	bapng := betterapng.NewBAPNG(f)
	if err != nil {
		log.Fatal(err)
	}
	bapng.Open()
	item.BAPNG = bapng
	images , err := item.BAPNG.ReadAllFramesAsPNG()
	fmt.Println(len(images))
	if err != nil {
		log.Fatal(err)
	}
	item.images = images

	item.Image = canvas.NewRasterWithPixels(func(x, y, w, h int) color.Color {
		return item.currentImage.At(x, y)
	})

	image, err := png.Decode(bytes.NewReader(item.images[0]))
	if err != nil {
		log.Fatal(err)
	}
	item.Image.SetMinSize(fyne.NewSize(float32(image.Bounds().Dx()), float32(image.Bounds().Dy())))
	item.frame = 0

	item.ExtendBaseWidget(item)
	go func() {
		ticker := time.NewTicker(time.Second / 60) // 60 FPS
		for range ticker.C {
			item.Update()
			canvas.Refresh(item)
		}
	}()

	return item
}

func (item *LoadingWidget) CreateRenderer() fyne.WidgetRenderer {
	c := item.Image
	return widget.NewSimpleRenderer(c)
}

func main() {
	a := app.New()
	drv := a.Driver()
	if drv, ok := drv.(desktop.Driver); ok {
		w := drv.CreateSplashWindow()

		widget := NewMyListItemWidget()

		w.SetContent(widget)
		w.ShowAndRun()
		// Customize your splash window here
	}
}
