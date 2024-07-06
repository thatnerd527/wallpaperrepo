package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"github.com/elliotchance/pie/v2"
	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
	"github.com/google/uuid"
	"github.com/thatnerd/betterapng"
)

type LoadingWidget struct {
	widget.BaseWidget
	FrameCounter *widget.Label
	frame        uint64
	Image        *canvas.Raster
	BAPNG        *betterapng.BAPNG
	images       [][]byte
	imageConfigs []betterapng.BAPNGFrame
	currentImage image.Image
}

func GenerateGUID() string {
	id := uuid.New()
	return fmt.Sprint(id.String())
}

func createShortcut(shortcutPath, targetPath, workingDir string) error {
	ole.CoInitialize(0)
	defer ole.CoUninitialize()

	oleShellObject, err := oleutil.CreateObject("WScript.Shell")
	if err != nil {
		return fmt.Errorf("failed to create WScript.Shell object: %v", err)
	}
	defer oleShellObject.Release()

	wshell, err := oleShellObject.QueryInterface(ole.IID_IDispatch)
	if err != nil {
		return fmt.Errorf("failed to query interface: %v", err)
	}
	defer wshell.Release()

	cs, err := oleutil.CallMethod(wshell, "CreateShortcut", shortcutPath)
	if err != nil {
		return fmt.Errorf("failed to create shortcut: %v", err)
	}
	defer cs.ToIDispatch().Release()

	_, err = oleutil.PutProperty(cs.ToIDispatch(), "TargetPath", targetPath)
	if err != nil {
		return fmt.Errorf("failed to set target path: %v", err)
	}

	_, err = oleutil.PutProperty(cs.ToIDispatch(), "WorkingDirectory", workingDir)
	if err != nil {
		return fmt.Errorf("failed to set working directory: %v", err)
	}

	_, err = oleutil.CallMethod(cs.ToIDispatch(), "Save")
	if err != nil {
		return fmt.Errorf("failed to save shortcut: %v", err)
	}

	return nil
}


func (item *LoadingWidget) Update() {
	item.FrameCounter.SetText(fmt.Sprintf("%d", item.frame))
	item.frame++
	if item.frame >= uint64(item.BAPNG.GetNumberOfFrames()) {
		item.frame = 0
	}
	item.currentImage, _ = betterapng.UniversalDecoder(item.images[item.frame], item.imageConfigs[item.frame].Codec)
}

func NewLoadingSplash() *LoadingWidget {

	item := &LoadingWidget{
		FrameCounter: widget.NewLabel("0"),
	}
	f, err := os.Open("installsplash.bapng")
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
	images, configs, err := item.BAPNG.ReadAllFramesAsRAW()
	fmt.Println(len(images))
	if err != nil {
		log.Fatal(err)
	}
	item.images = images
	item.imageConfigs = pie.Map(configs, func(x *betterapng.BAPNGFrame) betterapng.BAPNGFrame {
		return *x
	})

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

func startBackgroundProcess(command string, args []string, cwd string) {
    cmd := exec.Command(command, args...)
	cmd.Dir = cwd
    // For Windows, this detaches the process from the parent
    cmd.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
    err := cmd.Start()
    if err != nil {
        log.Fatalf("Failed to start background process: %v", err)
    }
    // Note: We do not wait for the process to finish.
}


func postSetup(folder string) {
	embedkey := GenerateGUID()
	installerembed := path.Join(folder, "embedkey")
	ufolder, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
	userstartupfolder := path.Join(ufolder, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
	startmenu := path.Join(ufolder, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs")
	err = os.WriteFile(installerembed, []byte(embedkey), 0644)
	if err != nil {
		log.Fatal(err)
	}

	bootstraphtml := path.Join(folder, "bootstrap", "index.html")
	bootstraphtmlcontents, err := os.ReadFile(bootstraphtml)
	if err != nil {
		log.Fatal(err)
	}
	bootstraphtmlcontents = []byte(strings.ReplaceAll(string(bootstraphtmlcontents), "<EMBEDKEY.THIS WILL BE REPLACED DURING INSTALLATION>", embedkey))
	err = os.WriteFile(bootstraphtml, bootstraphtmlcontents, 0644)
	if err != nil {
		log.Fatal(err)
	}
	err = createShortcut(path.Join(userstartupfolder,"Wallpaper System UI Startup.lnk"), path.Join(folder, "wallpaperuiserver.exe"), folder)
	if err != nil {
		log.Fatal(err)
	}
	startBackgroundProcess(path.Join(folder, "wallpaperuiserver.exe"), []string{}, folder)
	resolved, err := filepath.Abs(path.Join(folder, "bootstrap"))
	if err != nil {
		log.Fatal(err)
	}
	startBackgroundProcess("C:/Windows/explorer.exe", []string{resolved}, folder)

	// Create start menu shortcuts

	os.MkdirAll(path.Join(startmenu, "Wallpaper System"), 0755)
	err = createShortcut(path.Join(startmenu, "Wallpaper System", "Wallpaper System.lnk"), path.Join(folder, "wallpaperuiserver.exe"), folder)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	a := app.New()
	drv := a.Driver()
	installdir := flag.String("installdir", "", "The directory to install the program to")
	targetdir := ""
	flag.Parse()
	go func() {
		if *installdir == "" {
			folder, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}

			targetdir = path.Join(folder, "AppData", "Local", "Programs", "Wallpaper System")
			filepath.Abs(targetdir)
			extractZipToFolderWith7z("install.zip", targetdir)
			postSetup(targetdir)
			a.Quit()
		}
	}()
	if drv, ok := drv.(desktop.Driver); ok {
		w := drv.CreateSplashWindow()

		widget := NewLoadingSplash()

		w.SetContent(widget)
		w.ShowAndRun()
		// Customize your splash window here
	}


}
