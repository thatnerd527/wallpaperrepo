package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
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