package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

func restart() {
	popupchannel <- PopupRequest{Type: "stop"}
	fmt.Println("Sent stop message")
	time.Sleep(5 * time.Second)
	startPopupApp()
	fmt.Println("Started popup app")

	// Shutdown all addons
	for _, addon := range RuntimeAddons {
		addon.Shutdown()
	}
	fmt.Println("Shutted down all addons")
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

func autoStartHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	ufolder, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	shortcut := path.Join(ufolder, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup", "Wallpaper System UI Startup.lnk")
	cwd := path.Join(ufolder, "AppData", "Local", "Programs","Wallpaper System")
	if r.Method == "GET" {
		_, ok := os.Stat(shortcut)
		if ok != nil {
			w.Write([]byte("false"))
		} else {
			w.Write([]byte("true"))
		}
	}
	if r.Method == "POST" {
		if _, ok := r.URL.Query()["enable"]; ok {
			if _, ok := os.Stat(shortcut); ok != nil {
				path, err:= filepath.Abs(cwd)
				if err != nil {
					log.Fatal(err)
				}

				createShortcut(shortcut, os.Args[0], path)
				w.Write([]byte("true"))
				w.WriteHeader(http.StatusOK)
			} else {
				w.Write([]byte("false"))
				w.WriteHeader(http.StatusConflict)
			}
		} else {
			if _, ok := os.Stat(shortcut); ok == nil {
				os.Remove(shortcut)
				w.Write([]byte("true"))
				w.WriteHeader(http.StatusOK)
			} else {
				w.Write([]byte("false"))
				w.WriteHeader(http.StatusConflict)
			}
		}

	}

}

func restartHandler(w http.ResponseWriter, r *http.Request) {
	// Restart the popup app

	// Send a message to the popup app to restart
	restart()
	w.Header().Set("Access-Control-Allow-Origin", "*")

}
