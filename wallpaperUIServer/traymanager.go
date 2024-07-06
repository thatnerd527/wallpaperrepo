package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/getlantern/systray"
)

func determineHttpOrHttps(url string) string {
	data , err := http.Get(url)
	if err != nil {
		return "https"
	}
	if data.StatusCode == 200 {
		return "http"
	}
	if data.ContentLength == 0 {
		return "https"
	}
	return "https"
}

func createTray() {
	icon, err := os.ReadFile("icon.ico")
	if err != nil {
		panic(err)
	} else {
		systray.SetIcon(icon)
	}

	systray.SetTitle("Wallpaper System")
	systray.SetTooltip("Quick actions and settings for the wallpaper system")
	mSettings := systray.AddMenuItem("Settings", "Open the settings window")
	mRestart := systray.AddMenuItem("Restart", "Restart the wallpaper system")
	mQuit := systray.AddMenuItem("Quit", "Quit the wallpaper system")

	go func ()  {
		for {
			<-mSettings.ClickedCh
			// open settings window
			whattouse := determineHttpOrHttps("http://localhost:"+strconv.Itoa(allowedredirectport))
			res := CreatePopupRequest(whattouse + "://localhost:"+strconv.Itoa(allowedredirectport)+"?mode=settings","system", "Wallpaper System", "settings","Wallpaper Settings")
			res.trackingID = strings.ReplaceAll(GenerateGUID(),"-","")
			popupchannel <- res
		}
	}()

	go func ()  {
		for {
			<-mRestart.ClickedCh
			// restart the system
			restarthub.SendMessage(RestartRequest{false, "restart",GenerateGUID()})
			restart()
		}
	}()

	go func ()  {
		for {
			<-mQuit.ClickedCh
			// quit the system
			shutdown()
		}
	}()

}
