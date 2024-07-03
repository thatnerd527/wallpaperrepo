package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"github.com/sqweek/dialog"
)
var guidtofilepath = make(map[string]string)

func FileNameToMediaType(filename string) string {
	switch filepath.Ext(filename) {
	case ".mp4",".webm",".mov",".avi":
		return "Video"
	case ".jpg", ".jpeg", ".png",".avif", ".webp", ".gif":
		return "Image"
	case ".mp3", ".wav", ".flac", ".ogg":
		return "Audio"
	default:
		return "unknown"
	}
}

func simpleBackgroundHandler(w http.ResponseWriter, r *http.Request) {

	if (r.Method == "GET") {
		if r.URL.Query().Get("delete") == "true" {
			delete(cachedpreferences, "simplebackground")
			preferenceschannel.SendMessage(PreferenceUpdate{cachedpreferences, GenerateGUID(), false})
			w.Header().Add("Access-Control-Allow-Origin", "*")
			return
		}
		if val, ok := cachedpreferences["simplebackground"]; ok {
			w.Header().Set("Content-Type", "application/octet-stream")
			http.ServeFile(w, r, val.(string))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
		w.Header().Add("Access-Control-Allow-Origin", "*")
	} else if (r.Method == "POST") {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		zipPath2, err := dialog.File().Title("Select simple background file").
		Filter("Image files", "jpg", "jpeg", "png", "avif", "webp", "gif").
		Filter("WebM encoded with AV1, VP8, or VP9", "webm").Load()
		if err != nil {
			log.Println(err)
			return
		}
		cachedpreferences["simplebackground"] = zipPath2
		message, err := json.Marshal(cachedpreferences)
		if err != nil {
			log.Println(err)
			return
		}
		os.WriteFile("preferences.json", message, 0644)
		preferenceschannel.SendMessage(PreferenceUpdate{cachedpreferences, GenerateGUID(), false})
	}
}


func LoadAllMediaAsPanels(panelspath string, manifest AddonManifest) ([]TemplateCustomPanel, error) {
	files, err := os.ReadDir(panelspath)
	if err != nil {
		log.Println(err)
		return []TemplateCustomPanel{}, err
	}
	panels := []TemplateCustomPanel{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		guid := GenerateGUID()
		filename := file.Name()
		panel := TemplateCustomPanel{}
		panel.PanelType = FileNameToMediaType(filename)
		panel.LoaderPanelID = filename + "_panel_" + manifest.Name
		url, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v",secureport))
		if err != nil {
			log.Println(err)
			return []TemplateCustomPanel{}, err
		}
		url.Path = "/mediaregistry"
		q := url.Query()
		q.Add("guid", guid)
		url.RawQuery = q.Encode()
		panel.PanelContent = url.String()
		panel.PanelTitle = filename
		guidtofilepath[guid] = filepath.Join(panelspath, filename)
		panels = append(panels, panel)
	}
	return panels, nil
}

func LoadAllMediaAsTemplateBackgrounds(backgroundspath string, manifest AddonManifest) ([]TemplateCustomBackground, error) {
	files, err := os.ReadDir(backgroundspath)
	if err != nil {
		log.Println(err)
		return []TemplateCustomBackground{}, err
	}
	backgrounds := []TemplateCustomBackground{}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		guid := GenerateGUID()
		filename := file.Name()
		background := TemplateCustomBackground{}
		background.BackgroundType = FileNameToMediaType(filename)
		background.LoaderBackgroundID = filename + "_background_" + manifest.Name
		url, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%v",secureport))
		if err != nil {
			log.Println(err)
			return []TemplateCustomBackground{}, err
		}
		url.Path = "/mediaregistry"
		q := url.Query()
		q.Add("guid", guid)
		url.RawQuery = q.Encode()
		background.BackgroundContent = url.String()
		guidtofilepath[guid] = filepath.Join(backgroundspath, filename)
		backgrounds = append(backgrounds, background)
	}
	return backgrounds, nil
}