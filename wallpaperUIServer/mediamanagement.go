package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"

	"github.com/sqweek/dialog"
)
var guidtofilepath = make(map[string]string)


func FileNameToMediaType(filename string) string {
	switch filepath.Ext(filename) {
	case ".mp4",".webm",".mov",".avi",".mkv",".flv",".wmv",".mpg",".mpeg",".m4v",".3gp",".3g2":
		return "Video"
	case ".jpg", ".jpeg", ".png",".avif", ".webp", ".gif":
		return "Image"
	default:
		return "unknown"
	}
}

func setBackgroundFromCache(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		if !r.URL.Query().Has("delete") {
			filename, err := filepath.Abs(path.Join("result", r.URL.Query().Get("filename")))
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			pref := cachedpreferences
			if _, err := os.Stat(filename); err == nil {
				path2 := filename

				pref["simplebackground"] = path2
				addRecentMediaBackground(fmt.Sprint(hashCode(r.URL.Query().Get("filename"))), r.URL.Query().Get("filename"), "")
				cachedpreferences = pref
				marshalled, _ := json.Marshal(cachedpreferences)
				os.WriteFile("preferences.json", marshalled, 0755)
				fmt.Println("Set simple background to " + path2)

				preferenceschannel.SendMessage(PreferenceUpdate{cachedpreferences, GenerateGUID(), false})
			} else {
				log.Println("File not found")
				w.WriteHeader(http.StatusNotFound)
			}
			return
		} else {
			filename, err := filepath.Abs(path.Join("result", r.URL.Query().Get("filename")))
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			if _, err := os.Stat(filename); err == nil {
				err := os.Remove(filename)
				if err != nil {
					log.Println(err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
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
		filename := r.URL.Query().Get("resultfilename")
		zipPath2, err := dialog.File().Title("Select simple background file").
		Filter("Image files", "jpg", "jpeg", "png", "avif", "webp", "gif").
		Filter("Video Files", "mp4","webm","mov","avi","mkv","flv","wmv","mpg","mpeg","m4v","3gp","3g2").
		Filter("Other files","*").Load()
		if err != nil {
			log.Println(err)
			tmp := make(map[string]string)
			tmp["error"] = err.Error()
			tmp["status"] = "cancelled"
			tmp["guid"] = ""
			tmp["resultfile"] = ""
			data, err := json.Marshal(tmp)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			w.Write(data)
			return
		}

		encodingguid := GenerateGUID()
		encodertochan[encodingguid] = make(chan string)
		resultfile, err := startReencoding(zipPath2, encodertochan[encodingguid], filename)
		if err != nil {
			log.Println(err)
			tmp := make(map[string]string)
			tmp["error"] = err.Error()
			tmp["status"] = "errorprestart"
			tmp["guid"] = ""
			tmp["resultfile"] = ""
			data, err := json.Marshal(tmp)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Write(data)
			return
		}
		tmp := make(map[string]string)
		tmp["error"] = ""
		tmp["status"] = "encoding"
		tmp["guid"] = encodingguid
		tmp["resultfile"] = resultfile
		data, err := json.Marshal(tmp)
		if err != nil {
			log.Println(err)
			return
		}
		w.Write(data)
		return;
	}
}

func getEncodingStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == "GET" {
		if val, ok := encodertochan[r.URL.Query().Get("guid")]; ok {
			if (r.URL.Query().Get("action") == "status") {
				result := <-val
				_ , err := w.Write([]byte(result))
				if err != nil {
					log.Println("GETENCODINGSTATUS " + err.Error())
					val <- result
				}


			} else {
				var tmp = make(map[string]string)
				data2 := <- val
				tmp["status"] = data2
				data, err := json.Marshal(tmp)
				if err != nil {
					log.Println("GETENCODINGSTATUS " + err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_ , err = w.Write(data)
				if err != nil {
					log.Println("GETENCODINGSTATUS " + err.Error())
					val <- data2
				}
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	if r.Method == "POST" {
		if val, ok := encodertochan[r.URL.Query().Get("guid")]; ok {
			if r.URL.Query().Has("action") {
				val <- r.URL.Query().Get("action")

			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
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
		panel.LoaderPanelID = filename + "_panel_" + manifest.Name + "_" + manifest.ClientID
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
		background.LoaderBackgroundID = filename + "_background_" + manifest.Name + "_" + manifest.ClientID
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